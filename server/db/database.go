package db

import "hmb_fighting/server/types"

type Database interface {
	GetWeapons() (map[string]types.Weapon, error)
	GetShields() (map[string]types.Shield, error)
	GetTeams() (map[int]types.TeamConfig, error)
	GetCharacters() ([]types.Character, error)
	GetAbilities() (map[string]types.Ability, error)
	GetRoleConfig() (map[string]types.Role, error)

	SetUser(refreshToken string, user types.User) error
	GetUserByEmail(email string) (types.User, error)
	GetUserByRefresh(token string) (types.User, error)
	GetRoom(roomID string) (*types.Game, error)
	SetRoom(game *types.Game) error
}
