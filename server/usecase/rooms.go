package usecase

import (
	"fmt"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/jwt"
	"hmb_fighting/server/types"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (u *Usecase) CreateRoom(accessToken string) (string, error) {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return "", fmt.Errorf("Invalid token: %v", err)
	}

	room := u.initRoom()
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.Players[0] = claims.ClientID

	err = u.db.SetRoom(room)
	if err != nil {
		return "", fmt.Errorf("Failed to save room: %v", err)
	}

	return room.GameSessionId, nil
}

func (u *Usecase) RestartRoom(accessToken, roomID string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	room, err := u.db.GetRoom(roomID)
	if err != nil || room == nil {
		return fmt.Errorf("Room not found")
	}

	if room.Players[0] != claims.ClientID && room.Players[1] != claims.ClientID {
		return fmt.Errorf("Unauthorized")
	}

	room.Mutex.Lock()

	room.Phase = types.GamePhaseSetup
	room.Winner = -1
	room.CurrentTurn = -1
	room.InitialOrder = nil
	room.Battlelog = nil
	for i := range room.Board {
		for j := range room.Board[i] {
			room.Board[i][j] = -1
		}
	}
	for teamID := range room.Teams {
		for i := range room.Teams[teamID].Characters {
			char := &room.Teams[teamID].Characters[i]
			char.HP = 100
			char.Position = [2]int{-1, -1}
			char.Effects = nil
			char.SetAbilities(room.AbilitiesConfig)
		}
	}

	err = u.db.SetRoom(room)
	room.Mutex.Unlock()
	if err != nil {
		return fmt.Errorf("Failed to save room: %v", err)
	}

	u.broadcastRoomState(room)
	log.Printf("Room %s restarted by %s", roomID, claims.ClientID)
	return nil
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

	playerIndex := -1
	for i, playerID := range game.Players {
		if playerID == claims.ClientID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		game.Mutex.Unlock()
		return fmt.Errorf("You are not a player in this room")
	}

	delete(game.Players, playerIndex)

	for conn, client := range game.Connections {
		if client.ClientID == claims.ClientID {
			delete(game.Connections, conn)
		}
	}

	if len(game.Players) == 0 {
		game.Phase = types.GamePhasePickTeam
		log.Printf("Room %s is now empty after %s left", roomID, claims.ClientID)
	} else {
		log.Printf("Player %s left room %s, %d players remaining", claims.ClientID, roomID, len(game.Players))
	}

	err = u.db.SetRoom(game)
	game.Mutex.Unlock()
	if err != nil {
		return fmt.Errorf("Failed to update room: %v", err)
	}

	if len(game.Players) > 0 {
		u.broadcastRoomState(game)
	}

	log.Printf("Player %s left room %s", claims.ClientID, roomID)
	return nil
}

func (u *Usecase) initRoom() *entities.Room {
	weaponsConfig, err := u.db.GetWeapons()
	if err != nil {
		log.Fatalf("Failed to get weapons config: %v", err)
	}

	shieldsConfig, err := u.db.GetShields()
	if err != nil {
		log.Fatalf("Failed to get shields config: %v", err)
	}

	roleConfig, err := u.db.GetRoleConfig()
	if err != nil {
		log.Fatalf("Failed to get role config: %v", err)
	}

	abilitiesConfig, err := u.db.GetAbilities()
	if err != nil {
		log.Fatalf("Failed to get abilities config: %v", err)
	}

	game := &entities.Room{
		Connections:     make(map[*websocket.Conn]*entities.Client),
		GameSessionId:   uuid.New().String(),
		WeaponsConfig:   weaponsConfig,
		ShieldsConfig:   shieldsConfig,
		AbilitiesConfig: abilitiesConfig,
		RoleConfig:      roleConfig,
		CurrentTurn:     -1,
		Winner:          -1,
		Phase:           types.GamePhasePickTeam,
		Players:         make(map[int]string),
		Board:           types.BoardTypes{},
		Battlelog:       []entities.Battlelog{},
	}

	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}

	return game
}
