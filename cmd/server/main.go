package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type User struct {
	Name  string
	Email string
}

type Character struct {
	ID         int
	Name       string
	Team       int
	HP         int
	Stamina    int
	AttackMin  int
	AttackMax  int
	Defense    int
	Initiative int
	Weapon     string
	Shield     string
	Height     int
	Weight     int
	Position   [2]int
	Abilities  []Ability
	Effects    []Effect
}

type Ability struct {
	Name        string
	Type        string
	Description string
	Range       int
}

type Effect struct {
	Name       string
	Duration   int
	StaminaMod int
	AttackMod  int
	DefenseMod int
}

type GameState struct {
	Teams         [2]Team
	CurrentTurn   int
	Phase         string
	Board         [20][10]int
	TeamID        int
	ClientID      string
	GameSessionId string
}

type Team struct {
	Characters []Character
}

type Action struct {
	Type        string `json:"type"`
	CharacterID int    `json:"characterID"`
	Position    [2]int `json:"position"`
	TargetID    int    `json:"targetID"`
	Ability     string `json:"ability"`
	ClientID    string `json:"clientID"`
}

type Client struct {
	Conn      *websocket.Conn
	TeamID    int
	ClientID  string
	Spectator bool
	User      *User
}

type Game struct {
	Connections   map[*websocket.Conn]*Client
	Teams         [2]Team
	CurrentTurn   int
	Phase         string
	Board         [20][10]int
	GameSessionId string
	mutex         sync.Mutex
}

var mutex = sync.Mutex{}
var rooms = make(map[string]*Game)
var users = make(map[string]User)

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

func initGame() *Game {
	game := &Game{
		Connections:   make(map[*websocket.Conn]*Client),
		GameSessionId: uuid.New().String(),
		Teams: [2]Team{
			{Characters: []Character{
				{ID: 1, Name: "Vasya", Team: 0, HP: 1, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 5, Initiative: 8, Weapon: "falchion", Shield: "buckler", Height: 175, Weight: 80, Position: [2]int{2, 2}, Abilities: []Ability{{Name: "Takedown", Type: "wrestle", Description: "Attempts to take down the opponent", Range: 1}}},
				{ID: 2, Name: "Petya", Team: 0, HP: 1, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 180, Weight: 90, Position: [2]int{2, 3}, Abilities: []Ability{{Name: "Throw", Type: "wrestle", Description: "Throws the opponent", Range: 1}}},
				{ID: 3, Name: "Alexei", Team: 0, HP: 1, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "two_handed_sword", Shield: "", Height: 185, Weight: 95, Position: [2]int{7, 7}, Abilities: []Ability{{Name: "Pin", Type: "wrestle", Description: "Pins the opponent down", Range: 1}}},
				{ID: 4, Name: "Misha", Team: 0, HP: 1, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "spear", Shield: "buckler", Height: 170, Weight: 75, Position: [2]int{3, 4}, Abilities: []Ability{{Name: "Grapple", Type: "wrestle", Description: "Grapples the opponent", Range: 1}}},
				{ID: 5, Name: "Sasha", Team: 0, HP: 1, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "dagger", Shield: "shield", Height: 178, Weight: 85, Position: [2]int{4, 5}, Abilities: []Ability{{Name: "Lock", Type: "wrestle", Description: "Locks the opponent", Range: 1}}},
			}},
			{Characters: []Character{
				{ID: 6, Name: "Igor", Team: 1, HP: 1, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "falchion", Shield: "buckler", Height: 172, Weight: 78, Position: [2]int{17, 2}, Abilities: []Ability{{Name: "Takedown", Type: "wrestle", Description: "Attempts to take down the opponent", Range: 1}}},
				{ID: 7, Name: "Dima", Team: 1, HP: 1, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "two_handed_halberd", Shield: "", Height: 182, Weight: 92, Position: [2]int{17, 3}, Abilities: []Ability{{Name: "Throw", Type: "wrestle", Description: "Throws the opponent", Range: 1}}},
				{ID: 8, Name: "Kolya", Team: 1, HP: 1, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 176, Weight: 83, Position: [2]int{16, 4}, Abilities: []Ability{{Name: "Pin", Type: "wrestle", Description: "Pins the opponent down", Range: 1}}},
				{ID: 9, Name: "Roma", Team: 1, HP: 1, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "sword", Shield: "buckler", Height: 188, Weight: 98, Position: [2]int{15, 5}, Abilities: []Ability{{Name: "Grapple", Type: "wrestle", Description: "Grapples the opponent", Range: 1}}},
				{ID: 10, Name: "Zhenya", Team: 1, HP: 1, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 5, Initiative: 6, Weapon: "dagger", Shield: "shield", Height: 174, Weight: 80, Position: [2]int{14, 6}, Abilities: []Ability{{Name: "Lock", Type: "wrestle", Description: "Locks the opponent", Range: 1}}},
			}},
		},
		CurrentTurn: 3,
		Phase:       "move",
		Board:       [20][10]int{},
	}

	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}
	for _, team := range game.Teams {
		for _, char := range team.Characters {
			switch char.Shield {
			case "buckler":
				char.Defense += 3
			case "shield":
				char.Defense += 5
			}
			game.Board[char.Position[0]][char.Position[1]] = char.ID
		}
	}
	return game
}

func findCharacter(game *Game, id int) *Character {
	for i := range game.Teams {
		for j := range game.Teams[i].Characters {
			if game.Teams[i].Characters[j].ID == id {
				return &game.Teams[i].Characters[j]
			}
		}
	}
	return nil
}

func distance(pos1, pos2 [2]int) int {
	dist := abs(pos1[0]-pos2[0]) + abs(pos1[1]-pos2[1])
	log.Printf("Distance from (%d, %d) to (%d, %d) = %d", pos1[0], pos1[1], pos2[0], pos2[1], dist)
	return dist
}

func distanceToAbility(pos1, pos2 [2]int) int {
	dx := abs(pos1[0] - pos2[0])
	dy := abs(pos1[1] - pos2[1])
	dist := max(dx, dy)
	log.Printf("Chebyshev Distance from (%d, %d) to (%d, %d) = %d", pos1[0], pos1[1], pos2[0], pos2[1], dist)
	return dist
}

func distanceToAttack(pos1, pos2 [2]int, weapon string) int {
	isTwoHanded := weapon == "two_handed_halberd" || weapon == "two_handed_sword"
	if isTwoHanded {
		return max(abs(pos1[0]-pos2[0]), abs(pos1[1]-pos2[1]))
	}
	return distanceToAbility(pos1, pos2)
}

func countSurroundingEnemies(game *Game, char *Character) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x, y := char.Position[0]+dx, char.Position[1]+dy
			if x >= 0 && x < 20 && y >= 0 && y < 10 && game.Board[x][y] != -1 {
				target := findCharacter(game, game.Board[x][y])
				if target != nil && target.Team != char.Team && target.HP > 0 {
					count++
				}
			}
		}
	}
	log.Printf("Surrounding enemies for %s at (%d, %d): %d", char.Name, char.Position[0], char.Position[1], count)
	return count
}

func calculateDamage(attacker, target *Character, game *Game) int {
	baseDamage := rand.Intn(attacker.AttackMax-attacker.AttackMin+1) + attacker.AttackMin
	totalDefense := target.Defense
	for _, effect := range target.Effects {
		totalDefense += effect.DefenseMod
	}
	damage := baseDamage - totalDefense
	if damage < 0 {
		damage = 0
	}
	surroundingEnemies := countSurroundingEnemies(game, target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost
	if damage < 0 {
		damage = 0
	}
	log.Printf("Damage calculation: base=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, totalDefense, damageBoost, damage)
	return damage
}

func applyWrestlingMove(game *Game, attacker, target *Character, moveName string) {
	successChance := 20
	partialSuccessChance := 25
	nothingChance := 30
	failureChance := 10
	totalFailureChance := 5

	heightDiff := float64(attacker.Height-target.Height) / 10.0
	weightDiff := float64(attacker.Weight-target.Weight) / 10.0
	mod := int(heightDiff+weightDiff) * 5

	surroundingEnemies := countSurroundingEnemies(game, target)
	successBoost := surroundingEnemies * 5
	successChance += mod + successBoost
	partialSuccessChance += mod + successBoost
	failureChance -= mod / 2
	totalFailureChance -= mod / 2

	if successChance < 5 {
		successChance = 5
	}
	if partialSuccessChance < 5 {
		partialSuccessChance = 5
	}
	if failureChance < 5 {
		failureChance = 5
	}
	if totalFailureChance < 5 {
		totalFailureChance = 5
	}

	total := successChance + partialSuccessChance + nothingChance + failureChance + totalFailureChance
	if total != 100 {
		scale := float64(100) / float64(total)
		successChance = int(float64(successChance) * scale)
		nothingChance = int(float64(nothingChance) * scale)
		partialSuccessChance = int(float64(partialSuccessChance) * scale)
		failureChance = int(float64(failureChance) * scale)
		totalFailureChance = 100 - successChance - partialSuccessChance - failureChance - nothingChance
	}

	r := rand.Intn(100)
	log.Printf("%s attempts %s on %s: success=%d%%, partial=%d%%, failure=%d%%, totalFailure=%d%%, surroundingBoost=%d, roll=%d",
		attacker.Name, moveName, target.Name, successChance, partialSuccessChance, failureChance, totalFailureChance, successBoost, r)

	switch {
	case r < successChance:
		target.HP = 0
		log.Printf("%s successfully used %s on %s, knocking them out!", attacker.Name, moveName, target.Name)
		game.Board[target.Position[0]][target.Position[1]] = -1
	case r < successChance+partialSuccessChance:
		damage := calculateDamage(attacker, target, game)
		target.HP -= damage
		log.Printf("%s partially succeeded with %s on %s, dealing %d damage", attacker.Name, moveName, target.Name, damage)
		if target.HP <= 0 {
			game.Board[target.Position[0]][target.Position[1]] = -1
		}
	case r < successChance+partialSuccessChance+nothingChance:
		log.Printf("%s failed to use %s on %s - the move didn't connect!", attacker.Name, moveName, target.Name)
	case r < successChance+partialSuccessChance+nothingChance+failureChance:
		attacker.HP = 0
		target.HP = 0
		log.Printf("%s failed %s - both %s and %s are knocked out!", attacker.Name, moveName, attacker.Name, target.Name)
		game.Board[attacker.Position[0]][attacker.Position[1]] = -1
		game.Board[target.Position[0]][target.Position[1]] = -1
	default:
		attacker.HP = 0
		log.Printf("%s catastrophically failed %s - %s is knocked out!", attacker.Name, moveName, attacker.Name)
		game.Board[attacker.Position[0]][attacker.Position[1]] = -1
	}
}

func nextTurn(game *Game) {
	liveChars := []Character{}
	for _, team := range game.Teams {
		for _, char := range team.Characters {
			if char.HP > 0 {
				liveChars = append(liveChars, char)
			}
		}
	}
	if len(liveChars) == 0 {
		return
	}
	sortCharactersByInitiative(liveChars)
	currentIndex := -1
	for i, char := range liveChars {
		if char.ID == game.CurrentTurn {
			currentIndex = i
			break
		}
	}
	nextIndex := (currentIndex + 1) % len(liveChars)
	game.CurrentTurn = liveChars[nextIndex].ID
	game.Phase = "move"
	for i := range game.Teams {
		for j := range game.Teams[i].Characters {
			char := &game.Teams[i].Characters[j]
			for k := len(char.Effects) - 1; k >= 0; k-- {
				char.Effects[k].Duration--
				if char.Effects[k].Duration <= 0 {
					char.Effects = append(char.Effects[:k], char.Effects[k+1:]...)
				}
			}
		}
	}
}

func sortCharactersByInitiative(chars []Character) {
	for i := 0; i < len(chars)-1; i++ {
		for j := 0; j < len(chars)-i-1; j++ {
			if chars[j].Initiative < chars[j+1].Initiative {
				chars[j], chars[j+1] = chars[j+1], chars[j]
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	var user User
	user = users[clientID]
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
		log.Printf("New client %s connected to room %s, User: %v", clientID, room, client.User)
	}
	game.Connections[conn] = client
	game.mutex.Unlock()

	broadcastGameState(game)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", clientID, err)
			game.mutex.Lock()
			delete(game.Connections, conn)
			game.mutex.Unlock()
			broadcastGameState(game)
			return
		}

		var action Action
		if err := json.Unmarshal(message, &action); err != nil {
			log.Printf("Invalid action from %s: %v", clientID, err)
			continue
		}

		if client.Spectator && action.Type != "chat" {
			continue
		}

		game.mutex.Lock()
		currentChar := findCharacter(game, game.CurrentTurn)
		log.Printf("Received action: %+v, client.ClientID: %s", action, client.ClientID)
		if currentChar == nil || currentChar.Team != client.TeamID || action.ClientID != client.ClientID {
			log.Printf("Action rejected: currentChar: %v, Team: %d, action.ClientID: %s, client.ClientID: %s", currentChar, client.TeamID, action.ClientID, client.ClientID)
			game.mutex.Unlock()
			continue
		}

		switch action.Type {
		case "move":
			if game.Phase == "move" && game.Board[action.Position[0]][action.Position[1]] == -1 {
				dist := distance(currentChar.Position, action.Position)
				totalStamina := currentChar.Stamina
				for _, effect := range currentChar.Effects {
					totalStamina += effect.StaminaMod
				}
				if totalStamina < 1 {
					totalStamina = 1
				}
				log.Printf("Move attempt: char %s at (%d, %d), target (%d, %d), dist: %d, totalStamina: %d",
					currentChar.Name, currentChar.Position[0], currentChar.Position[1], action.Position[0], action.Position[1], dist, totalStamina)
				if dist <= totalStamina {
					game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
					currentChar.Position = action.Position
					game.Board[action.Position[0]][action.Position[1]] = currentChar.ID
					game.Phase = "action"
					log.Printf("Character %d moved to (%d, %d)", currentChar.ID, action.Position[0], action.Position[1])
				} else {
					log.Printf("Move rejected: distance %d exceeds stamina %d", dist, totalStamina)
				}
			} else {
				log.Printf("Move rejected: phase=%s or target cell occupied", game.Phase)
			}
		case "attack":
			if game.Phase == "action" {
				target := findCharacter(game, action.TargetID)
				if target != nil && target.Team != currentChar.Team {
					dist := distanceToAttack(currentChar.Position, target.Position, currentChar.Weapon)
					isTwoHanded := currentChar.Weapon == "two_handed_halberd" || currentChar.Weapon == "two_handed_sword"
					weaponRange := 1
					if isTwoHanded {
						weaponRange = 2
					}
					if (isTwoHanded && dist == weaponRange) || (!isTwoHanded && dist <= weaponRange) {
						damage := calculateDamage(currentChar, target, game)
						target.HP -= damage
						log.Printf("%s attacked %s for %d damage", currentChar.Name, target.Name, damage)
						if target.HP <= 0 {
							game.Board[target.Position[0]][target.Position[1]] = -1
						}
						nextTurn(game)
					}
				}
			}
		case "ability":
			if game.Phase == "action" {
				target := findCharacter(game, action.TargetID)
				if target != nil && target.Team != currentChar.Team {
					dist := distanceToAbility(currentChar.Position, target.Position)
					for i, ability := range currentChar.Abilities {
						if ability.Name == action.Ability {
							if dist <= ability.Range {
								if ability.Type == "wrestle" {
									applyWrestlingMove(game, currentChar, target, action.Ability)
									log.Printf("%s used ability %s on %s", currentChar.Name, action.Ability, target.Name)
								} else {
									log.Printf("%s used %s on %s (non-wrestling ability not implemented)", currentChar.Name, action.Ability, target.Name)
								}
								currentChar.Abilities = append(currentChar.Abilities[:i], currentChar.Abilities[i+1:]...)
								if target.HP <= 0 {
									game.Board[target.Position[0]][target.Position[1]] = -1
								}
								nextTurn(game)
							} else {
								log.Printf("Ability rejected: target out of range (dist: %d, range: %d)", dist, ability.Range)
							}
							break
						}
					}
				}
			}
		case "end_turn":
			nextTurn(game)
			log.Printf("Turn ended for %s", currentChar.Name)
		case "restart":
			game = initGame()
			for conn, c := range game.Connections {
				c.TeamID = len(game.Connections) % 2
				game.Connections[conn] = c
			}
			log.Printf("Game restarted by %s", client.ClientID)
		case "chat":
			for conn := range game.Connections {
				conn.WriteJSON(map[string]string{"type": "chat", "message": fmt.Sprintf("%s: %s", currentChar.Name, string(message))})
			}
		}
		game.mutex.Unlock()
		broadcastGameState(game)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/ws", enableCORS(http.HandlerFunc(handleWebSocket)))
	http.Handle("/register", enableCORS(http.HandlerFunc(handleRegister)))
	http.Handle("/check-client", enableCORS(http.HandlerFunc(handleCheckClient)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
