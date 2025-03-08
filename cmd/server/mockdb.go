package main

// MockDatabase имитирует базу данных с предопределёнными данными
type MockDatabase struct{}

// GetWeapons возвращает конфигурацию оружия
func (m *MockDatabase) GetWeapons() (map[string]Weapon, error) {
	return map[string]Weapon{
		"falchion":           {Name: "falchion", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"axe":                {Name: "axe", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 8},
		"two_handed_sword":   {Name: "two_handed_sword", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"spear":              {Name: "spear", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 0},
		"dagger":             {Name: "dagger", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 0},
		"two_handed_halberd": {Name: "two_handed_halberd", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 8},
		"sword":              {Name: "sword", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetShields() (map[string]Shield, error) {
	return map[string]Shield{
		"buckler": {Name: "buckler", DefenseBonus: 1, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 1},
		"shield":  {Name: "shield", DefenseBonus: 2, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 0},
		"tower":   {Name: "tower", DefenseBonus: 3, ImageURL: "./static/shields/default.png", AttackBonus: 0, GrappleBonus: -1},
		"":        {Name: "none", DefenseBonus: 0, ImageURL: "", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetAbilities() (map[string]Ability, error) {
	return map[string]Ability{
		"takedown": {
			Name:        "Takedown",
			Type:        "wrestle",
			Description: "Attempts to take down the opponent",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"throw": {
			Name:        "Throw",
			Type:        "wrestle",
			Description: "Throws the opponent",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"pin": {
			Name:        "Pin",
			Type:        "wrestle",
			Description: "Pins the opponent down",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"grapple": {
			Name:        "Grapple",
			Type:        "wrestle",
			Description: "Grapples the opponent",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"lock": {
			Name:        "Lock",
			Type:        "wrestle",
			Description: "Locks the opponent",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
	}, nil
}

// GetTeamsConfig возвращает конфигурацию команд
func (m *MockDatabase) GetTeamsConfig() ([2]TeamConfig, error) {
	return [2]TeamConfig{
		{IconURL: "./static/icons/team0.png"},
		{IconURL: "./static/icons/team1.png"},
	}, nil
}

// GetCharacters возвращает список персонажей
func (m *MockDatabase) GetCharacters() ([]Character, error) {
	return []Character{
		{ID: 1, Name: "Vasya", Team: 0, HP: 10, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 5, Initiative: 8, Weapon: "falchion", Shield: "buckler", Height: 175, Weight: 80, Position: [2]int{1, 0}, Abilities: []string{"takedown"}, ImageURL: "./static/characters/sashya.png"},
		{ID: 2, Name: "Petya", Team: 0, HP: 10, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 180, Weight: 90, Position: [2]int{1, 2}, Abilities: []string{"throw"}, ImageURL: "./static/characters/suslik.png"},
		{ID: 3, Name: "Alexei", Team: 0, HP: 10, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "two_handed_sword", Shield: "", Height: 185, Weight: 95, Position: [2]int{4, 4}, Abilities: []string{"pin", "grapple", "takedown", "throw", "lock"}, ImageURL: "./static/characters/benya.png"},
		{ID: 4, Name: "Misha", Team: 0, HP: 10, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "spear", Shield: "buckler", Height: 170, Weight: 75, Position: [2]int{1, 6}, Abilities: []string{"grapple"}, ImageURL: "./static/characters/bear_coub.png"},
		{ID: 5, Name: "Sasha", Team: 0, HP: 10, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "dagger", Shield: "shield", Height: 178, Weight: 85, Position: [2]int{1, 8}, Abilities: []string{"lock"}, ImageURL: "./static/characters/kolya.png"},
		{ID: 6, Name: "Igor", Team: 1, HP: 10, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "falchion", Shield: "buckler", Height: 172, Weight: 78, Position: [2]int{14, 0}, Abilities: []string{"takedown"}, ImageURL: "./static/characters/sasha.png"},
		{ID: 7, Name: "Dima", Team: 1, HP: 10, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "two_handed_halberd", Shield: "", Height: 182, Weight: 92, Position: [2]int{14, 2}, Abilities: []string{"throw"}, ImageURL: "./static/characters/baba.png"},
		{ID: 8, Name: "Kolya", Team: 1, HP: 10, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 176, Weight: 83, Position: [2]int{12, 4}, Abilities: []string{"pin"}, ImageURL: "./static/characters/default.png"},
		{ID: 9, Name: "Roma", Team: 1, HP: 10, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "sword", Shield: "buckler", Height: 188, Weight: 98, Position: [2]int{14, 6}, Abilities: []string{"grapple"}, ImageURL: "./static/characters/default.png"},
		{ID: 10, Name: "Zhenya", Team: 1, HP: 10, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 5, Initiative: 6, Weapon: "dagger", Shield: "shield", Height: 174, Weight: 80, Position: [2]int{14, 8}, Abilities: []string{"lock"}, ImageURL: "./static/characters/default.png"},
	}, nil
}

// NewMockDatabase создаёт новый экземпляр MockDatabase
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}
