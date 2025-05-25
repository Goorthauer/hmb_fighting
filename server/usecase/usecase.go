package usecase

import (
	"encoding/json"
	"fmt"
	"hmb_fighting/server/db"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/jwt"
	"hmb_fighting/server/types"
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

func (u *Usecase) HandleWebSocket(conn *websocket.Conn, roomID, accessToken string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	room, err := u.db.GetRoom(roomID)
	if err != nil || room == nil {
		return fmt.Errorf("Room not found")
	}

	// Добавляем клиента в игру
	room.Mutex.Lock()
	if room.Connections == nil {
		room.Connections = make(map[*websocket.Conn]*entities.Client)
	}
	if room.Players == nil {
		room.Players = make(map[int]string)
	}
	client := &entities.Client{
		Conn:     conn,
		ClientID: claims.ClientID,
		User:     &entities.User{Name: claims.Email, Email: claims.Email},
	}
	if room.Players[0] == claims.ClientID {
		client.TeamID = 0
		client.Spectator = false
	} else if len(room.Players) < 2 && claims.Role == types.UserRoleSpectator {
		client.TeamID = 1
		client.Spectator = false
		room.Players[1] = claims.ClientID
		claims.Role = types.UserRolePlayer
	} else if room.Players[1] == claims.ClientID {
		client.TeamID = 1
		client.Spectator = false
	} else {
		client.TeamID = -1
		client.Spectator = true
	}
	room.Connections[conn] = client
	room.Mutex.Unlock()

	u.broadcastRoomState(room)

	// Читаем сообщения в цикле
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			room.Mutex.Lock()
			delete(room.Connections, conn)
			room.Mutex.Unlock()
			u.broadcastRoomState(room)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Client %s disconnected normally", claims.ClientID)
			} else {
				log.Printf("Error reading message from %s: %v", claims.ClientID, err)
			}
			return nil
		}

		var action entities.Action
		if err := json.Unmarshal(msg, &action); err != nil {
			log.Printf("Invalid action from %s: %v", claims.ClientID, err)
			continue
		}

		if action.ClientID != claims.ClientID {
			log.Printf("ClientID mismatch: expected %s, got %s", claims.ClientID, action.ClientID)
			continue
		}

		// Обрабатываем действие с минимальной блокировкой
		u.processAction(room, client, action, claims)
		u.broadcastRoomState(room)
	}
}

func (u *Usecase) handleSetupPhase(game *entities.Room, client *entities.Client, action entities.Action) {
	if action.Type == types.ActionPlace && client.TeamID >= 0 {
		char := game.FindCharacter(action.CharacterID)
		if char != nil && char.HP > 0 && char.TeamID == client.TeamID && isPositionOnBoard(action.Position) {
			if game.Board[action.Position[0]][action.Position[1]] == -1 && ((char.TeamID == 0 && action.Position[0] < 8) || (char.TeamID == 1 && action.Position[0] >= 8)) {
				if isPositionOnBoard(char.Position) {
					game.Board[char.Position[0]][char.Position[1]] = -1
				}
				char.Position = action.Position
				game.Board[action.Position[0]][action.Position[1]] = char.ID
				log.Printf("%s placed %s at (%d, %d)", client.ClientID, char.Name, action.Position[0], action.Position[1])
			}
		}
	} else if action.Type == types.ActionStart && client.TeamID >= 0 {
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
				}
			}
			if allPlaced {
				game.Phase = types.GamePhaseMove
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
				log.Printf("Room %s started by %s", game.GameSessionId, client.ClientID)
			}
		}
	}
}

func (u *Usecase) handleGamePhase(game *entities.Room, client *entities.Client, action entities.Action, claims *jwt.Claims) {
	currentChar := game.FindCharacter(game.CurrentTurn)
	if currentChar == nil || currentChar.TeamID != client.TeamID {
		log.Printf("Not your turn or invalid character: %s", claims.ClientID)
		return
	}

	switch action.Type {
	case types.ActionMove:
		u.handleMoveAction(game, currentChar, action)
	case types.ActionAttack:
		u.handleAttackAction(game, currentChar, action)
	case types.ActionAbility:
		u.handleAbilityAction(game, currentChar, action)
	case types.ActionEndTurn:
		game.NextTurn()
		log.Printf("%s ended turn", claims.ClientID)
	}
}

func (u *Usecase) handleAbilityAction(game *entities.Room, currentChar *entities.Character, action entities.Action) {
	target := game.FindCharacter(action.TargetID)
	if game.Phase == types.GamePhaseAction && canTarget(currentChar, target) {
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

func (u *Usecase) handleMoveAction(game *entities.Room, currentChar *entities.Character, action entities.Action) {
	haveMove :=
		game.Phase == types.GamePhaseMove &&
			isPositionOnBoard(currentChar.Position) &&
			isPositionOnBoard(action.Position)
	if !haveMove || game.Board[action.Position[0]][action.Position[1]] != -1 {
		return
	}
	path, opportunityAttacks := game.FindPath(currentChar.Position[0], currentChar.Position[1], action.Position[0], action.Position[1], currentChar.Stamina, game.Board, currentChar.ID)
	if len(path) > 0 {
		totalDamage := 0
		for _, oa := range opportunityAttacks {
			attacker := game.FindCharacter(oa.AttackerID)
			if oa.Type == types.OpportunityAttackTrip {
				game.SetBattleLog(fmt.Sprintf("%s проводит подсечку и %s безвольно падает!", attacker.Name, currentChar.Name))
				totalDamage += oa.Damage
			} else if oa.Type == types.OpportunityAttackAttack {
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
			game.Phase = types.GamePhaseAction
			game.SetBattleLog(fmt.Sprintf("%s ходит на (%d, %d)", currentChar.Name, action.Position[0], action.Position[1]))
		} else {
			game.SetBattleLog(fmt.Sprintf("%s был накаутирован во время хода (%d, %d)", currentChar.Name, action.Position[0], action.Position[1]))
			game.NextTurn()
		}
	} else {
		game.SetBattleLog(fmt.Sprintf("%s пытался пройти в (%d, %d), но путь заблокирован", currentChar.Name, action.Position[0], action.Position[1]))
	}
}

func (u *Usecase) handleAttackAction(game *entities.Room, currentChar *entities.Character, action entities.Action) {
	target := game.FindCharacter(action.TargetID)
	if (game.Phase == types.GamePhaseMove || game.Phase == types.GamePhaseAction) && canTarget(currentChar, target) {
		weapon := game.WeaponsConfig[currentChar.Weapon]
		weaponRange := weapon.Range
		if game.DistanceToAttack(currentChar.Position, target.Position, weapon) <= weaponRange {
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

func (u *Usecase) broadcastRoomState(room *entities.Room) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	log.Printf("Broadcasting to %d clients", len(room.Connections))

	teams := [2]entities.Team{}
	for i := 0; i < 2; i++ {
		if team, ok := room.Teams[i]; ok {
			teams[i] = team
		}
	}
	teamsConfig := [2]entities.TeamConfig{}
	for i := 0; i < 2; i++ {
		if config, ok := room.TeamsConfig[i]; ok {
			teamsConfig[i] = config
		}
	}

	for conn, client := range room.Connections {
		state := entities.GameState{
			Teams:           teams,
			Winner:          room.Winner,
			CurrentTurn:     room.CurrentTurn,
			Phase:           room.Phase,
			Board:           room.Board,
			TeamID:          client.TeamID,
			ClientID:        client.ClientID,
			GameSessionId:   room.GameSessionId,
			WeaponsConfig:   room.WeaponsConfig,
			AbilitiesConfig: room.AbilitiesConfig,
			ShieldsConfig:   room.ShieldsConfig,
			TeamsConfig:     teamsConfig,
			Battlelog:       room.Battlelog, // Добавляем Battlelog
		}
		if err := conn.WriteJSON(state); err != nil {
			log.Printf("Error sending room state to %s: %v", client.ClientID, err)
			conn.Close()
			delete(room.Connections, conn)
		}
	}
}

func (u *Usecase) processAction(game *entities.Room, client *entities.Client, action entities.Action, claims *jwt.Claims) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if game.Phase == types.GamePhaseSetup {
		u.handleSetupPhase(game, client, action)
	} else {
		u.handleGamePhase(game, client, action, claims)
	}
}

func generateClientID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func canTarget(attacker, target *entities.Character) bool {
	return attacker != nil &&
		target != nil &&
		attacker.HP > 0 &&
		target.HP > 0 &&
		attacker.TeamID != target.TeamID &&
		isPositionOnBoard(attacker.Position) &&
		isPositionOnBoard(target.Position)
}

func isPositionOnBoard(position [2]int) bool {
	return position[0] >= 0 &&
		position[0] < types.BoardVerticalSize &&
		position[1] >= 0 &&
		position[1] < types.BoardHorizontalSize
}
