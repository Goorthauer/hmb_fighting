package db

import (
	"hmb_fighting/server/entities"
)

type Database interface {
	GetWeapons() (map[string]entities.Weapon, error)
	GetShields() (map[string]entities.Shield, error)
	GetTeams() (map[int]entities.TeamConfig, error)
	GetCharacters() ([]entities.Character, error)
	GetAbilities() (map[string]entities.Ability, error)
	GetRoleConfig() (map[string]entities.Role, error)

	SetUser(refreshToken string, user entities.User) error
	GetUserByEmail(email string) (entities.User, error)
	GetUserByRefresh(token string) (entities.User, error)
	GetRoom(roomID string) (*entities.Room, error)
	SetRoom(game *entities.Room) error
}
