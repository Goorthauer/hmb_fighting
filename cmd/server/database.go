package main

// Database определяет интерфейс для работы с данными игры
type Database interface {
	GetWeapons() (map[string]Weapon, error)
	GetShields() (map[string]Shield, error)
	GetTeamsConfig() (map[int]TeamConfig, error)
	GetCharacters() ([]Character, error)
	GetAbilities() (map[string]Ability, error) // Новый метод
	GetRoleConfig() (map[string]Role, error)

	SetUser(refreshToken string, user User) error
	GetUserByEmail(email string) (User, error)
	GetUserByRefresh(token string) (User, error)
	GetRoom(roomID string) (*Game, error)
	SetRoom(game *Game) error
}
