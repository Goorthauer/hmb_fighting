package usecase

import (
	"encoding/json"
	"fmt"
	"hmb_fighting/cmd/server/db"
	"hmb_fighting/cmd/server/dtos"
	"hmb_fighting/cmd/server/game"
	"hmb_fighting/cmd/server/jwt"
	"hmb_fighting/cmd/server/types"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Usecase struct {
	db db.Database
}

func NewUsecase(db db.Database) *Usecase {
	return &Usecase{db: db}
}

func generateClientID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func (u *Usecase) RegisterUser(currentUser types.User) (*dtos.RegisterUserResp, error) {
	user, err := u.db.GetUserByEmail(currentUser.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %v", err)
	}

	if user.Email == "" {
		user = currentUser
		user.ID = generateClientID()
	}

	tokenPair, err := jwt.GenerateTokenPair(user, "spectator")
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %v", err)
	}

	err = u.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to save user with refresh token: %v", err)
	}

	return &dtos.RegisterUserResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}, nil
}

func (u *Usecase) RefreshToken(refreshToken string) (*dtos.RegisterUserResp, error) {
	user, err := u.db.GetUserByRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %v", err)
	}

	tokenPair, err := jwt.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %v", err)
	}

	err = u.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to update refresh token: %v", err)
	}

	return &dtos.RegisterUserResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}, nil
}

func (u *Usecase) CheckClient(clientID, accessToken string) (bool, error) {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return false, fmt.Errorf("Invalid token: %v", err)
	}

	user, err := u.db.GetUserByEmail(claims.Email)
	if err != nil {
		return false, fmt.Errorf("Failed to get user: %v", err)
	}

	return user.ID == clientID && user.ID == claims.ClientID, nil
}

func (u *Usecase) CreateRoom(accessToken string) (string, error) {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return "", fmt.Errorf("Invalid token: %v", err)
	}

	game := game.InitGame(u.db)
	game.Mutex.Lock()
	defer game.Mutex.Unlock()
	game.Players[0] = claims.ClientID

	err = u.db.SetRoom(game)
	if err != nil {
		return "", fmt.Errorf("Failed to save room: %v", err)
	}

	return game.GameSessionId, nil
}

func (u *Usecase) RestartRoom(accessToken, roomID string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return fmt.Errorf("Room not found")
	}

	if game.Players[0] != claims.ClientID && game.Players[1] != claims.ClientID {
		return fmt.Errorf("Unauthorized")
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

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
			char := &game.Teams[teamID].Characters[i]
			char.HP = 100
			char.Position = [2]int{-1, -1}
		}
	}

	err = u.db.SetRoom(game)
	if err != nil {
		return fmt.Errorf("Failed to save room: %v", err)
	}

	u.broadcastGameState(game)
	log.Printf("Room %s restarted by %s", roomID, claims.ClientID)
	return nil
}

func (u *Usecase) HandleWebSocket(conn *websocket.Conn, room, accessToken string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(room)
	if err != nil || game == nil {
		return fmt.Errorf("Room not found")
	}

	// Добавляем клиента в игру
	game.Mutex.Lock()
	client := &types.Client{
		Conn:     conn,
		ClientID: claims.ClientID,
		User:     &types.User{Name: claims.Email, Email: claims.Email},
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
	game.Mutex.Unlock()

	u.broadcastGameState(game)

	// Читаем сообщения в цикле
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			game.Mutex.Lock()
			delete(game.Connections, conn)
			game.Mutex.Unlock()
			u.broadcastGameState(game)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Client %s disconnected normally", claims.ClientID)
			} else {
				log.Printf("Error reading message from %s: %v", claims.ClientID, err)
			}
			return nil
		}

		var action types.Action
		if err := json.Unmarshal(msg, &action); err != nil {
			log.Printf("Invalid action from %s: %v", claims.ClientID, err)
			continue
		}

		if action.ClientID != claims.ClientID {
			log.Printf("ClientID mismatch: expected %s, got %s", claims.ClientID, action.ClientID)
			continue
		}

		// Обрабатываем действие с минимальной блокировкой
		u.processAction(game, client, action, claims)
		u.broadcastGameState(game)
	}
}

func (u *Usecase) processAction(game *types.Game, client *types.Client, action types.Action, claims *jwt.Claims) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if game.Phase == "setup" {
		u.handleSetupPhase(game, client, action)
	} else {
		u.handleGamePhase(game, client, action, claims)
	}
}

func (u *Usecase) SelectTeam() (*dtos.SelectTeamResp, error) {
	teams, err := u.db.GetTeamsConfig()
	if err != nil {
		return nil, fmt.Errorf("Teams not found: %v", err)
	}
	characters, err := u.db.GetCharacters()
	if err != nil {
		return nil, fmt.Errorf("Characters not found: %v", err)
	}

	outChars := make(map[int][]types.Character)
	for _, char := range characters {
		if !char.IsActive {
			continue
		}
		outChars[char.TeamID] = append(outChars[char.TeamID], char)
	}

	return &dtos.SelectTeamResp{
		AvailableTeams: teams,
		Characters:     outChars,
	}, nil
}

func (u *Usecase) SetTeam(roomID string, realTeamID int, accessToken string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return fmt.Errorf("Room not found")
	}
	for _, v := range game.TeamsConfig {
		if realTeamID == v.ID {
			return fmt.Errorf("Team exists")
		}
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
	if teamID == -1 {
		return fmt.Errorf("Invalid team ID")
	}

	teams, err := u.db.GetTeamsConfig()
	if err != nil {
		return fmt.Errorf("Teams not found: %v", err)
	}
	characters, err := u.db.GetCharacters()
	if err != nil {
		return fmt.Errorf("Characters not found: %v", err)
	}

	characterTeam := make([]types.Character, 0)
	for _, char := range characters {
		if !char.IsActive {
			continue
		}
		char.SetAbilities(game.AbilitiesConfig)
		char.Position = [2]int{-1, -1}
		if char.IsTitanArmour {
			char.Wrestling += 1
			char.Stamina += 1
			char.Initiative += 1
			char.Defense -= 2
			char.HP -= 5
			if char.HP < 1 {
				char.HP = 1
			}
			if char.Defense < 0 {
				char.Defense = 0
			}
		}
		if char.TeamID == realTeamID {
			char.TeamID = teamID
			characterTeam = append(characterTeam, char)
		}
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if game.TeamsConfig == nil {
		game.TeamsConfig = make(map[int]types.TeamConfig)
	}
	if game.Teams == nil {
		game.Teams = make(map[int]types.Team)
	}
	game.TeamsConfig[teamID] = teams[realTeamID]
	game.Teams[teamID] = types.Team{Characters: characterTeam}

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
		if len(game.InitialOrder) == 0 {
			game.InitTurnOrder()
		}
		game.Phase = "setup"
	}

	return u.db.SetRoom(game)
}

func (u *Usecase) CheckTeams(roomID, accessToken string) (bool, error) {
	_, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return false, fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return false, fmt.Errorf("Room not found")
	}

	return game.Phase == "setup", nil
}

func (u *Usecase) LeaveRoom(accessToken, roomID string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return fmt.Errorf("Room not found")
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	playerIndex := -1
	for i, playerID := range game.Players {
		if playerID == claims.ClientID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return fmt.Errorf("You are not a player in this room")
	}

	delete(game.Players, playerIndex)

	for conn, client := range game.Connections {
		if client.ClientID == claims.ClientID {
			delete(game.Connections, conn)
		}
	}

	if len(game.Players) == 0 {
		game.Phase = "pick_team"
		log.Printf("Room %s is now empty after %s left", roomID, claims.ClientID)
	} else {
		log.Printf("Player %s left room %s, %d players remaining", claims.ClientID, roomID, len(game.Players))
	}

	err = u.db.SetRoom(game)
	if err != nil {
		return fmt.Errorf("Failed to update room: %v", err)
	}

	if len(game.Players) > 0 {
		u.broadcastGameState(game)
	}

	log.Printf("Player %s left room %s", claims.ClientID, roomID)
	return nil
}

func (u *Usecase) handleSetupPhase(game *types.Game, client *types.Client, action types.Action) {
	if action.Type == "place" && client.TeamID >= 0 {
		char := game.FindCharacter(action.CharacterID)
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
						char.HP = 0
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

func (u *Usecase) handleGamePhase(game *types.Game, client *types.Client, action types.Action, claims *jwt.Claims) {
	currentChar := game.FindCharacter(game.CurrentTurn)
	if currentChar == nil || currentChar.TeamID != client.TeamID {
		log.Printf("Not your turn or invalid character: %s", claims.ClientID)
		return
	}

	switch action.Type {
	case "move":
		u.handleMoveAction(game, currentChar, action)
	case "attack":
		u.handleAttackAction(game, currentChar, action)
	case "ability":
		u.handleAbilityAction(game, currentChar, action)
	case "end_turn":
		game.NextTurn()
		log.Printf("%s ended turn", claims.ClientID)
	}
}

func (u *Usecase) handleAbilityAction(game *types.Game, currentChar *types.Character, action types.Action) {
	target := game.FindCharacter(action.TargetID)
	if game.Phase == "action" && target != nil && target.TeamID != currentChar.TeamID {
		ability, exists := game.AbilitiesConfig[strings.ToLower(action.Ability)]
		if exists && game.DistanceToAbility(currentChar.Position, target.Position) <= ability.Range {
			for i, abilityID := range currentChar.Abilities {
				if abilityID == action.Ability {
					game.ApplyWrestlingMove(currentChar, target, strings.ToLower(ability.Name))
					currentChar.Abilities = append(currentChar.Abilities[:i], currentChar.Abilities[i+1:]...)
					game.NextTurn()
					break
				}
			}
		}
	}
}

func (u *Usecase) handleMoveAction(game *types.Game, currentChar *types.Character, action types.Action) {
	if game.Phase == "move" && action.Position[0] >= 0 && action.Position[0] < 16 && action.Position[1] >= 0 && action.Position[1] < 9 {
		if game.Board[action.Position[0]][action.Position[1]] == -1 {
			path, opportunityAttacks := game.FindPath(currentChar.Position[0], currentChar.Position[1], action.Position[0], action.Position[1], currentChar.Stamina, game.Board, currentChar.ID)
			if len(path) > 0 {
				totalDamage := 0
				for _, oa := range opportunityAttacks {
					attacker := game.FindCharacter(oa.AttackerID)
					if oa.Type == "trip" {
						game.SetBattleLog(fmt.Sprintf("%s проводит подсечку и %s безвольно падает!", attacker.Name, currentChar.Name))
						totalDamage += oa.Damage
					} else if oa.Type == "attack" {
						game.SetBattleLog(fmt.Sprintf("%s атакует вслед %s на  %d урона!", attacker.Name, currentChar.Name, oa.Damage))
						totalDamage += oa.Damage
					}
					currentChar.HP -= oa.Damage
					if currentChar.HP <= 0 {
						game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
						break
					}
				}

				if currentChar.HP > 0 {
					game.Board[currentChar.Position[0]][currentChar.Position[1]] = -1
					currentChar.Position = action.Position
					game.Board[action.Position[0]][action.Position[1]] = currentChar.ID
					game.Phase = "action"
					game.SetBattleLog(fmt.Sprintf("%s ходит на (%d, %d)", currentChar.Name, action.Position[0], action.Position[1]))
				} else {
					game.SetBattleLog(fmt.Sprintf("%s был накаутирован во время хода (%d, %d)", currentChar.Name, action.Position[0], action.Position[1]))
					game.NextTurn()
				}
			} else {
				game.SetBattleLog(fmt.Sprintf("%s пытался пройти в (%d, %d), но путь заблокирован", currentChar.Name, action.Position[0], action.Position[1]))
			}
		}
	}
}

func (u *Usecase) handleAttackAction(game *types.Game, currentChar *types.Character, action types.Action) {
	target := game.FindCharacter(action.TargetID)
	if (game.Phase == "move" || game.Phase == "action") && target != nil && target.TeamID != currentChar.TeamID {
		weaponRange := game.WeaponsConfig[currentChar.Weapon].Range
		if game.DistanceToAttack(currentChar.Position, target.Position, game.WeaponsConfig[currentChar.Weapon]) <= weaponRange {
			damage := game.CalculateDamage(currentChar, target)
			target.HP -= damage
			if target.HP <= 0 {
				game.Board[target.Position[0]][target.Position[1]] = -1
				game.SetBattleLog(
					fmt.Sprintf("%s атаковал %s на %d урона и поверг его!",
						currentChar.Name,
						target.Name, damage))
			} else {
				game.SetBattleLog(
					fmt.Sprintf("%s атаковал %s на  %d урона (Осталось здоровья: %d)", currentChar.Name, target.Name, damage, target.HP))
			}
			game.NextTurn()
		}
	}
}

func (u *Usecase) broadcastGameState(game *types.Game) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	log.Printf("Broadcasting to %d clients", len(game.Connections))

	teams := [2]types.Team{}
	for i := 0; i < 2; i++ {
		if team, ok := game.Teams[i]; ok {
			teams[i] = team
		}
	}
	teamsConfig := [2]types.TeamConfig{}
	for i := 0; i < 2; i++ {
		if config, ok := game.TeamsConfig[i]; ok {
			teamsConfig[i] = config
		}
	}

	for conn, client := range game.Connections {
		state := types.GameState{
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
			Battlelog:       game.Battlelog, // Добавляем Battlelog
		}
		if err := conn.WriteJSON(state); err != nil {
			log.Printf("Error sending game state to %s: %v", client.ClientID, err)
			conn.Close()
			delete(game.Connections, conn)
		}
	}
}
