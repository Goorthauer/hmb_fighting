package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
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
	var currentUser User
	if err := json.NewDecoder(r.Body).Decode(&currentUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if currentUser.Name == "" || currentUser.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}
	user, ok := users[currentUser.Email]
	if !ok {
		user = currentUser
		user.ID = generateClientID()
		users[user.Email] = user
	}
	tokenPair, err := generateTokenPair(user, "spectator")
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	usersWithRefresh[tokenPair.RefreshToken] = user
	mutex.Unlock()

	response := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ClientID     string `json:"clientID"`
	}{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Registered user: %s (%s) with ClientID: %s", user.Name, user.Email, response.ClientID)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, ok := findUserByRefreshToken(req.RefreshToken)
	if !ok {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	tokenPair, err := refreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	response := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ClientID     string `json:"clientID"`
	}{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func findUserByRefreshToken(refreshToken string) (User, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	user, ok := usersWithRefresh[refreshToken]
	return user, ok
}

func handleCheckClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID    string `json:"clientID"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := validateToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	mutex.Lock()
	_, exists := users[req.ClientID]
	mutex.Unlock()

	if claims.ClientID != req.ClientID {
		exists = false
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": exists})
	log.Printf("Checked clientID: %s, valid: %v", req.ClientID, exists)
}

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := validateToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	db := NewMockDatabase()
	game := initGame(db)
	game.Players[0] = claims.ClientID

	mutex.Lock()
	rooms[game.GameSessionId] = game
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"roomID": game.GameSessionId})
	log.Printf("Room %s created by %s", game.GameSessionId, claims.ClientID)
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"accessToken"`
		RoomID      string `json:"roomID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := validateToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	mutex.Lock()
	game, exists := rooms[req.RoomID]
	mutex.Unlock()
	if !exists || (game.Players[0] != claims.ClientID && game.Players[1] != claims.ClientID) {
		http.Error(w, "Room not found or unauthorized", http.StatusForbidden)
		return
	}

	game.mutex.Lock()
	defer game.mutex.Unlock()
	game.SetupPhase = true
	game.Phase = "setup"
	game.Winner = -1
	game.CurrentTurn = -1
	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}
	for teamID := range game.Teams {
		for i := range game.Teams[teamID].Characters {
			game.Teams[teamID].Characters[i].HP = 10
			game.Teams[teamID].Characters[i].Position = [2]int{-1, -1}
		}
	}

	broadcastGameState(game)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Room restarted"))
	log.Printf("Room %s restarted by %s", req.RoomID, claims.ClientID)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	accessToken := r.URL.Query().Get("accessToken")

	log.Printf("WebSocket request received: room=%s, accessToken=%s", room, accessToken)

	if room == "" || room == "null" {
		http.Error(w, "Room parameter is required", http.StatusBadRequest)
		return
	}
	if accessToken == "" {
		http.Error(w, "AccessToken is required", http.StatusBadRequest)
		return
	}

	claims, err := validateToken(accessToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	game, exists := rooms[room]
	mutex.Unlock()
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Room not found"})
		return
	}

	game.mutex.Lock()
	client := &Client{
		Conn:     conn,
		ClientID: claims.ClientID,
		User:     &User{Name: claims.Email, Email: claims.Email},
	}
	if game.Players[0] == claims.ClientID {
		client.TeamID = 0
		client.Spectator = false
	} else if len(game.Players) < 2 && claims.Role == "spectator" {
		client.TeamID = 1
		client.Spectator = false
		game.Players[1] = claims.ClientID
		claims.Role = "player"
	} else if game.Players[1] == claims.ClientID {
		client.TeamID = 1
		client.Spectator = false
	} else {
		client.TeamID = -1
		client.Spectator = true
	}
	game.Connections[conn] = client
	game.mutex.Unlock()

	broadcastGameState(game)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			game.mutex.Lock()
			delete(game.Connections, conn)
			game.mutex.Unlock()
			broadcastGameState(game)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Client %s disconnected normally", claims.ClientID)
			} else {
				log.Printf("Error reading message from %s: %v", claims.ClientID, err)
			}
			return
		}

		var action Action
		if err := json.Unmarshal(msg, &action); err != nil {
			log.Printf("Invalid action from %s: %v", claims.ClientID, err)
			continue
		}

		if action.ClientID != claims.ClientID {
			log.Printf("ClientID mismatch: expected %s, got %s", claims.ClientID, action.ClientID)
			continue
		}

		game.mutex.Lock()
		if game.Phase == "setup" {
			if action.Type == "place" && client.TeamID >= 0 {
				char := findCharacter(game, action.CharacterID)
				if char != nil && char.Team == client.TeamID && action.Position[0] >= 0 && action.Position[0] < 16 && action.Position[1] >= 0 && action.Position[1] < 9 {
					if game.Board[action.Position[0]][action.Position[1]] == -1 && ((char.Team == 0 && action.Position[0] < 8) || (char.Team == 1 && action.Position[0] >= 8)) {
						char.Position = action.Position
						game.Board[action.Position[0]][action.Position[1]] = char.ID
						log.Printf("%s placed %s at (%d, %d)", client.ClientID, char.Name, action.Position[0], action.Position[1])
					}
				}
			} else if action.Type == "start" && client.TeamID >= 0 {
				if len(game.Players) == 2 {
					allPlaced := true
					for _, team := range game.Teams {
						placed := 0
						for _, char := range team.Characters {
							if char.Position[0] != -1 {
								placed++
							}
						}
						if placed < 5 {
							allPlaced = false
							break
						}
					}
					if allPlaced {
						game.SetupPhase = false
						game.Phase = "move"
						game.CurrentTurn = game.Teams[0].Characters[0].ID
						log.Printf("Game %s started by %s", game.GameSessionId, claims.ClientID)
					}
				}
			}
		} else {
			currentChar := findCharacter(game, game.CurrentTurn)
			if currentChar == nil || currentChar.Team != client.TeamID {
				log.Printf("Not your turn or invalid character: %s", claims.ClientID)
				game.mutex.Unlock()
				continue
			}

			switch action.Type {
			case "move":
				if game.Phase == "move" && action.Position[0] >= 0 && action.Position[0] < 16 && action.Position[1] >= 0 && action.Position[1] < 9 {
					if distance(currentChar.Position, action.Position) <= currentChar.Stamina && game.Board[action.Position[0]][action.Position[1]] == -1 {
						game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
						currentChar.Position = action.Position
						game.Board[action.Position[0]][action.Position[1]] = currentChar.ID
						game.Phase = "action"
						log.Printf("%s moved %s to (%d, %d)", claims.ClientID, currentChar.Name, action.Position[0], action.Position[1])
					}
				}
			case "attack":
				target := findCharacter(game, action.TargetID)
				if game.Phase == "action" && target != nil && target.Team != currentChar.Team {
					weaponRange := game.WeaponsConfig[currentChar.Weapon].Range
					if distanceToAttack(currentChar.Position, target.Position, game.WeaponsConfig[currentChar.Weapon]) <= weaponRange {
						damage := calculateDamage(currentChar, target, game)
						target.HP -= damage
						if target.HP <= 0 {
							game.Board[target.Position[0]][target.Position[1]] = -1
						}
						log.Printf("%s attacked %s for %d damage (HP left: %d)", currentChar.Name, target.Name, damage, target.HP)
						nextTurn(game)
					}
				}
			case "ability":
				target := findCharacter(game, action.TargetID)
				if game.Phase == "action" && target != nil && target.Team != currentChar.Team {
					ability, exists := game.AbilitiesConfig[strings.ToLower(action.Ability)]
					if exists && distanceToAbility(currentChar.Position, target.Position) <= ability.Range {
						for i, abilityID := range currentChar.Abilities {
							if abilityID == action.Ability {
								applyWrestlingMove(game, currentChar, target, strings.ToLower(ability.Name))
								currentChar.Abilities = append(currentChar.Abilities[:i], currentChar.Abilities[i+1:]...)
								nextTurn(game)
								break
							}
						}
					}
				}
			case "end_turn":
				if game.Phase == "move" || game.Phase == "action" {
					nextTurn(game)
					log.Printf("%s ended turn", claims.ClientID)
				}
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

	teams := [2]Team{}
	for i := 0; i < 2; i++ {
		if team, ok := game.Teams[i]; ok {
			teams[i] = team
		}
	}
	teamsConfig := [2]TeamConfig{}
	for i := 0; i < 2; i++ {
		if config, ok := game.TeamsConfig[i]; ok {
			teamsConfig[i] = config
		}
	}

	for conn, client := range game.Connections {
		state := GameState{
			Teams:           teams,
			CurrentTurn:     game.CurrentTurn,
			Phase:           game.Phase,
			Board:           game.Board,
			TeamID:          client.TeamID,
			ClientID:        client.ClientID,
			GameSessionId:   game.GameSessionId,
			WeaponsConfig:   game.WeaponsConfig,
			AbilitiesConfig: game.AbilitiesConfig,
			ShieldsConfig:   game.ShieldsConfig,
			TeamsConfig:     teamsConfig,
		}
		if err := conn.WriteJSON(state); err != nil {
			log.Printf("Error sending game state to %s: %v", client.ClientID, err)
			conn.Close()
			delete(game.Connections, conn)
		}
	}
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("invalid or expired token")
	}
	return claims, nil
}

func countPlacedCharacters(team Team) int {
	count := 0
	for _, char := range team.Characters {
		if char.Position[0] != -1 && char.Position[1] != -1 {
			count++
		}
	}
	return count
}
