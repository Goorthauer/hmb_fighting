package main

// Database определяет интерфейс для работы с данными игры
type Database interface {
	GetWeapons() (map[string]Weapon, error)
	GetShields() (map[string]Shield, error)
	GetTeamsConfig() (map[int]TeamConfig, error)
	GetCharacters() ([]Character, error)
	GetAbilities() (map[string]Ability, error) // Новый метод
}
