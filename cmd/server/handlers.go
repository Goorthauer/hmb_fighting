package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func generateClientID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	clientID := generateClientID()
	mutex.Lock()
	users[clientID] = user
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"clientID": clientID})
	log.Printf("Registered user: %s (%s) with ClientID: %s", user.Name, user.Email, clientID)
}

func handleCheckClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID string `json:"clientID"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	_, exists := users[req.ClientID]
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": exists})
	log.Printf("Checked clientID: %s, valid: %v", req.ClientID, exists)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	spectator := r.URL.Query().Get("spectator") == "true"
	clientID := r.URL.Query().Get("clientID")

	log.Printf("WebSocket request received: room=%s, spectator=%t, clientID=%s", room, spectator, clientID)

	if room == "" || room == "null" {
		log.Println("Rejecting connection: room parameter is missing or invalid")
		http.Error(w, "Room parameter is required", http.StatusBadRequest)
		return
	}

	if clientID == "" || clientID == "undefined" {
		log.Println("Rejecting connection: clientID is required")
		http.Error(w, "ClientID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	user := users[clientID]
	game, exists := rooms[room]
	if !exists {
		game = initGame()
		rooms[room] = game
	}
	mutex.Unlock()

	game.mutex.Lock()
	var client *Client
	for _, existingClient := range game.Connections {
		if existingClient.ClientID == clientID {
			client = existingClient
			client.Conn = conn
			log.Printf("Reconnected client %s to room %s", clientID, room)
			break
		}
	}
	if client == nil {
		teamID := 0
		if len(game.Connections)%2 == 1 && !spectator {
			teamID = 1
		}
		client = &Client{
			Conn:      conn,
			TeamID:    teamID,
			ClientID:  clientID,
			Spectator: spectator,
			User:      &user,
		}
		game.Connections[conn] = client
		log.Printf("New client %s joined room %s as team %d (spectator: %t)", clientID, room, teamID, spectator)
	}
	game.mutex.Unlock()

	broadcastGameState(game)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", clientID, err)
			game.mutex.Lock()
			delete(game.Connections, conn)
			game.mutex.Unlock()
			broadcastGameState(game)
			return
		}

		var action Action
		if err := json.Unmarshal(msg, &action); err != nil {
			log.Printf("Invalid action from %s: %v", clientID, err)
			continue
		}

		if action.ClientID != clientID {
			log.Printf("ClientID mismatch: expected %s, got %s", clientID, action.ClientID)
			continue
		}

		game.mutex.Lock()
		currentChar := findCharacter(game, game.CurrentTurn)
		if currentChar == nil || currentChar.Team != client.TeamID {
			log.Printf("Not your turn or invalid character: %s", clientID)
			game.mutex.Unlock()
			continue
		}

		switch action.Type {
		case "move":
			if game.Phase != "move" || distance(currentChar.Position, action.Position) > currentChar.Stamina || game.Board[action.Position[0]][action.Position[1]] != -1 {
				log.Printf("Invalid move by %s: phase=%s, distance=%d, stamina=%d, target=%d", clientID, game.Phase, distance(currentChar.Position, action.Position), currentChar.Stamina, game.Board[action.Position[0]][action.Position[1]])
			} else {
				game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
				currentChar.Position = action.Position
				game.Board[action.Position[0]][action.Position[1]] = currentChar.ID
				game.Phase = "action"
				log.Printf("%s moved %s to (%d, %d)", clientID, currentChar.Name, action.Position[0], action.Position[1])
			}
		case "attack":
			target := findCharacter(game, action.TargetID)
			if game.Phase != "action" || target == nil || target.Team == currentChar.Team || distanceToAttack(currentChar.Position, target.Position, game.WeaponsConfig[currentChar.Weapon]) > game.WeaponsConfig[currentChar.Weapon].Range {
				log.Printf("Invalid attack by %s on %d", clientID, action.TargetID)
			} else {
				damage := calculateDamage(currentChar, target, game)
				target.HP -= damage
				if target.HP <= 0 {
					game.Board[target.Position[0]][target.Position[1]] = -1
				}
				log.Printf("%s attacked %s for %d damage (HP left: %d)", currentChar.Name, target.Name, damage, target.HP)
				nextTurn(game)
			}
		case "ability":
			target := findCharacter(game, action.TargetID)
			if game.Phase != "action" || target == nil || target.Team == currentChar.Team {
				log.Printf("Invalid ability use by %s on %d", clientID, action.TargetID)
			} else {
				for i, ability := range currentChar.Abilities {
					if ability.Name == action.Ability && distanceToAbility(currentChar.Position, target.Position) <= ability.Range {
						applyWrestlingMove(game, currentChar, target, ability.Name)
						currentChar.Abilities = append(currentChar.Abilities[:i], currentChar.Abilities[i+1:]...)
						break
					}
				}
				nextTurn(game)
			}
		case "end_turn":
			if game.Phase == "move" || game.Phase == "action" {
				nextTurn(game)
				log.Printf("%s ended turn", clientID)
			}
		}
		game.mutex.Unlock()
		broadcastGameState(game)
	}
}

func broadcastGameState(game *Game) {
	game.mutex.Lock()
	defer game.mutex.Unlock()
	log.Printf("Broadcasting to %d clients", len(game.Connections))
	for conn, client := range game.Connections {
		state := GameState{
			Teams:         game.Teams,
			CurrentTurn:   game.CurrentTurn,
			Phase:         game.Phase,
			Board:         game.Board,
			TeamID:        client.TeamID,
			ClientID:      client.ClientID,
			GameSessionId: game.GameSessionId,
			WeaponsConfig: game.WeaponsConfig,
			ShieldsConfig: game.ShieldsConfig,
			TeamsConfig:   game.TeamsConfig,
		}
		stateJSON, err := json.Marshal(state)
		if err != nil {
			log.Printf("Error marshaling game state for client %s: %v", client.ClientID, err)
			continue
		}
		log.Printf("Sending state to client %s: %s", client.ClientID, string(stateJSON))
		err = conn.WriteJSON(state)
		if err != nil {
			log.Printf("Error sending game state to %s: %v", client.ClientID, err)
			conn.Close()
			delete(game.Connections, conn)
		}
	}
}
