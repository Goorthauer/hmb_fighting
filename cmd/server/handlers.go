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

type Handler struct {
	db Database
}

func NewHandler(db Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var currentUser User
	if err := json.NewDecoder(r.Body).Decode(&currentUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if currentUser.Name == "" || currentUser.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByEmail(currentUser.Email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		user = currentUser
		user.ID = generateClientID()
	}

	tokenPair, err := generateTokenPair(user, "spectator")
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = h.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		http.Error(w, "Failed to save user with refresh token", http.StatusInternalServerError)
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
	log.Printf("Registered user: %s (%s) with ClientID: %s", user.Name, user.Email, response.ClientID)
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByRefresh(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	tokenPair, err := refreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	err = h.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
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

func (h *Handler) handleCheckClient(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.db.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	exists := user.ID == req.ClientID && user.ID == claims.ClientID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"valid":   exists,
	})
	log.Printf("Checked clientID: %s, valid: %v", req.ClientID, exists)
}

func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
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

	game := initGame(h.db)
	game.Players[0] = claims.ClientID

	err = h.db.SetRoom(game)
	if err != nil {
		http.Error(w, "Failed to save room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"roomID":  game.GameSessionId,
	})
	log.Printf("Room %s created by %s", game.GameSessionId, claims.ClientID)
}

func (h *Handler) handleRestart(w http.ResponseWriter, r *http.Request) {
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

	game, err := h.db.GetRoom(req.RoomID)
	if err != nil || game == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if game.Players[0] != claims.ClientID && game.Players[1] != claims.ClientID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	game.mutex.Lock()
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
			// Сбрасываем HP и позицию, но сохраняем эффекты Titan Armour
			char := &game.Teams[teamID].Characters[i]
			char.HP = 100
			char.Position = [2]int{-1, -1}
		}
	}
	game.mutex.Unlock()
	broadcastGameState(game)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Room restarted"))
	log.Printf("Room %s restarted by %s", req.RoomID, claims.ClientID)
}

func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	game, err := h.db.GetRoom(room)
	if err != nil || game == nil {
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
			handleSetupPhase(game, client, action)
		} else {
			handleGamePhase(game, client, action, claims)
		}
		game.mutex.Unlock()
		broadcastGameState(game)
	}
}

func (h *Handler) handleSelectTeam(w http.ResponseWriter, r *http.Request) {

	teams, err := h.db.GetTeamsConfig()
	if err != nil {
		http.Error(w, "Teams not found", http.StatusNotFound)
		return
	}
	characters, err := h.db.GetCharacters()
	if err != nil {
		http.Error(w, "Char not found", http.StatusNotFound)
	}

	outChars := make(map[int][]Character)
	for _, char := range characters {
		if !char.IsActive {
			continue
		}
		outChars[char.TeamID] = append(outChars[char.TeamID], char)
	}

	response := struct {
		AvailableTeams map[int]TeamConfig  `json:"availableTeams"`
		Characters     map[int][]Character `json:"characters"`
	}{
		AvailableTeams: teams,
		Characters:     outChars,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleSetTeam(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID      string `json:"roomID"`
		RealTeamID  int    `json:"realTeamID"`
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

	game, err := h.db.GetRoom(req.RoomID)
	if err != nil || game == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	teamID := -1
	if game.Players[0] == claims.ClientID {
		teamID = 0
		claims.Role = "player"
	} else if len(game.Players) < 2 && claims.Role == "spectator" {
		teamID = 1
		game.Players[1] = claims.ClientID
		claims.Role = "player"
	} else if game.Players[1] == claims.ClientID {
		teamID = 1
	}
	for i := range game.Connections {
		if game.Connections[i].ClientID == claims.ClientID {
			game.Connections[i].Spectator = false
		}
	}
	if teamID == -1 {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}
	teams, err := h.db.GetTeamsConfig()
	if err != nil {
		http.Error(w, "Teams not found", http.StatusNotFound)
		return
	}
	characters, err := h.db.GetCharacters()
	if err != nil {
		http.Error(w, "Char not found", http.StatusNotFound)
	}

	characterTeam := make([]Character, 0)
	for _, char := range characters {
		if !char.IsActive {
			continue
		}
		char.SetAbilities(game.AbilitiesConfig)
		char.Position = [2]int{-1, -1}
		// Применяем эффекты Titan Armour при инициализации
		if char.IsTitanArmour {
			char.Wrestling += 1
			char.Stamina += 1
			char.Initiative += 1
			char.Defense -= 2
			char.HP -= 5
			if char.HP < 1 {
				char.HP = 1 // Минимальное значение HP
			}
			if char.Defense < 0 {
				char.Defense = 0 // Минимальное значение защиты
			}
		}
		if char.TeamID == req.RealTeamID {
			char.TeamID = teamID
			characterTeam = append(characterTeam, char)
		}
	}

	game.mutex.Lock()
	teamConfigList := make(map[int]TeamConfig)
	teamList := make(map[int]Team)
	if game.TeamsConfig == nil {
		game.TeamsConfig = teamConfigList
	}
	if game.Teams == nil {
		game.Teams = teamList
	}
	game.TeamsConfig[teamID] = teams[req.RealTeamID]
	game.Teams[teamID] = Team{characterTeam}

	if len(game.Players) == 2 {
		for _, team := range game.Teams {
			for i := range team.Characters {
				char := &team.Characters[i]
				if shield, ok := game.ShieldsConfig[char.Shield]; ok {
					char.Defense += shield.DefenseBonus
					char.AttackMin += shield.AttackBonus
					char.AttackMax += shield.AttackBonus
				}
				if weapon, ok := game.WeaponsConfig[char.Weapon]; ok {
					char.AttackMin += weapon.AttackBonus
					char.AttackMax += weapon.AttackBonus
				}
			}
		}
		game.Phase = "setup"
	}
	game.mutex.Unlock()

	h.db.SetRoom(game)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Team assigned",
	})
}

func (h *Handler) handleCheckTeams(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID      string `json:"roomID"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := validateToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	game, err := h.db.GetRoom(req.RoomID)
	if err != nil || game == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	// Логика для проверки, выбрали ли обе команды свои команды
	allTeamsSelected := game.Phase == "setup"

	response := struct {
		AllTeamsSelected bool `json:"allTeamsSelected"`
	}{
		AllTeamsSelected: allTeamsSelected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSetupPhase(game *Game, client *Client, action Action) {
	if action.Type == "place" && client.TeamID >= 0 {
		char := findCharacter(game, action.CharacterID)
		if char != nil && char.TeamID == client.TeamID && action.Position[0] >= 0 && action.Position[0] < 16 && action.Position[1] >= 0 && action.Position[1] < 9 {
			if game.Board[action.Position[0]][action.Position[1]] == -1 && ((char.TeamID == 0 && action.Position[0] < 8) || (char.TeamID == 1 && action.Position[0] >= 8)) {
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
				for i := range team.Characters {
					char := &team.Characters[i]
					if char.Position[0] == -1 && char.Position[1] == -1 {
						char.HP = 0 // Убиваем непоставленных
						log.Printf("%s was killed due to not being placed", char.Name)
					} else if char.HP > 0 {
						placed++
					}
				}
				if placed < 5 {
					allPlaced = false
					break
				}
			}
			if allPlaced {
				game.Phase = "move"
				var maxInitiativeCharacterID int
				maxInitiative := -1
				for _, team := range game.Teams {
					for _, char := range team.Characters {
						if char.Initiative > maxInitiative && char.HP > 0 {
							maxInitiative = char.Initiative
							maxInitiativeCharacterID = char.ID
						}
					}
				}
				game.CurrentTurn = maxInitiativeCharacterID
				log.Printf("Game %s started by %s", game.GameSessionId, client.ClientID)
			}
		}
	}
}

func handleGamePhase(game *Game, client *Client, action Action, claims *Claims) {
	currentChar := findCharacter(game, game.CurrentTurn)
	if currentChar == nil || currentChar.TeamID != client.TeamID {
		log.Printf("Not your turn or invalid character: %s", claims.ClientID)
		return
	}

	switch action.Type {
	case "move":
		handleMoveAction(game, currentChar, action)
	case "attack":
		handleAttackAction(game, currentChar, action)
	case "ability":
		handleAbilityAction(game, currentChar, action)
	case "end_turn":
		nextTurn(game)
		log.Printf("%s ended turn", claims.ClientID)
	}
}

func handleMoveAction(game *Game, currentChar *Character, action Action) {
	if game.Phase == "move" && action.Position[0] >= 0 && action.Position[0] < 16 && action.Position[1] >= 0 && action.Position[1] < 9 {
		if game.Board[action.Position[0]][action.Position[1]] == -1 {
			path, opportunityAttacks := findPath(currentChar.Position[0], currentChar.Position[1], action.Position[0], action.Position[1], currentChar.Stamina, game.Board, currentChar.ID, game)
			if len(path) > 0 {
				// Применяем атаки в догонку
				totalDamage := 0
				for _, oa := range opportunityAttacks {
					attacker := findCharacter(game, oa.AttackerID)
					if oa.Type == "trip" {
						log.Printf("%s tripped %s, knocking them out!", attacker.Name, currentChar.Name)
						totalDamage += oa.Damage
					} else if oa.Type == "attack" {
						log.Printf("%s hit %s for %d damage in pursuit!", attacker.Name, currentChar.Name, oa.Damage)
						totalDamage += oa.Damage
					}
					currentChar.HP -= oa.Damage
					if currentChar.HP <= 0 {
						game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
						break
					}
				}

				// Если персонаж выжил, выполняем перемещение
				if currentChar.HP > 0 {
					game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
					currentChar.Position = action.Position
					game.Board[action.Position[0]][action.Position[1]] = currentChar.ID
					game.Phase = "action"
					log.Printf("%s moved to (%d, %d)", currentChar.Name, action.Position[0], action.Position[1])
				} else {
					log.Printf("%s was knocked out during move to (%d, %d)", currentChar.Name, action.Position[0], action.Position[1])
					nextTurn(game) // Завершаем ход, если персонаж погиб
				}
			} else {
				log.Printf("%s tried to move %s to (%d, %d), but path blocked or out of stamina", currentChar.Name, currentChar.Name, action.Position[0], action.Position[1])
			}
		}
	}
}

func handleAttackAction(game *Game, currentChar *Character, action Action) {
	target := findCharacter(game, action.TargetID)
	if (game.Phase == "move" || game.Phase == "action") && target != nil && target.TeamID != currentChar.TeamID {
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
}

func handleAbilityAction(game *Game, currentChar *Character, action Action) {
	target := findCharacter(game, action.TargetID)
	if game.Phase == "action" && target != nil && target.TeamID != currentChar.TeamID {
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
		//log.Printf("client %s is a %s spectator", client.TeamID, client.Spectator)
		//if client.Spectator && client.TeamID >= 0 {
		//	client.Spectator = false
		//}
		state := GameState{
			Teams:           teams,
			Winner:          game.Winner,
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
