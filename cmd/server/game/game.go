package game

import (
	"hmb_fighting/cmd/server/db"
	"hmb_fighting/cmd/server/types"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func InitGame(db db.Database) *types.Game {
	weaponsConfig, err := db.GetWeapons()
	if err != nil {
		log.Fatalf("Failed to get weapons config: %v", err)
	}

	shieldsConfig, err := db.GetShields()
	if err != nil {
		log.Fatalf("Failed to get shields config: %v", err)
	}

	roleConfig, err := db.GetRoleConfig()
	if err != nil {
		log.Fatalf("Failed to get role config: %v", err)
	}

	abilitiesConfig, err := db.GetAbilities()
	if err != nil {
		log.Fatalf("Failed to get abilities config: %v", err)
	}

	game := &types.Game{
		Connections:     make(map[*websocket.Conn]*types.Client),
		GameSessionId:   uuid.New().String(),
		WeaponsConfig:   weaponsConfig,
		ShieldsConfig:   shieldsConfig,
		AbilitiesConfig: abilitiesConfig,
		RoleConfig:      roleConfig,
		CurrentTurn:     -1,
		Phase:           "pick_team",
		Players:         make(map[int]string),
		Board:           [16][9]int{},
		Battlelog:       []string{},
	}

	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}

	return game
}
