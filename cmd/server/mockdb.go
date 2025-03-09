package main

// MockDatabase имитирует базу данных с предопределёнными данными
type MockDatabase struct{}

// GetWeapons возвращает конфигурацию оружия
func (m *MockDatabase) GetWeapons() (map[string]Weapon, error) {
	return map[string]Weapon{
		"falchion":           {Name: "falchion", DisplayName: "Фальшион", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"axe":                {Name: "axe", DisplayName: "Топор", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 8},
		"two_handed_sword":   {Name: "two_handed_sword", DisplayName: "Двуручный меч", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"two_handed_halberd": {Name: "two_handed_halberd", DisplayName: "Алебарда", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 8},
		"sword":              {Name: "sword", DisplayName: "Меч", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetShields() (map[string]Shield, error) {
	return map[string]Shield{
		"buckler": {Name: "buckler", DisplayName: "Баклер", DefenseBonus: 1, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 1},
		"shield":  {Name: "shield", DisplayName: "Тарч", DefenseBonus: 2, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 0},
		"tower":   {Name: "tower", DisplayName: "Ростовой щит", DefenseBonus: 3, ImageURL: "./static/shields/default.png", AttackBonus: 0, GrappleBonus: -1},
		"":        {Name: "none", DefenseBonus: 0, ImageURL: "", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetAbilities() (map[string]Ability, error) {
	abilities := map[string]Ability{
		"yama_arashi": {
			Name:        "yama_arashi",
			DisplayName: "Подхват",
			Description: "Мощный бросок через бедро, использующий силу и импульс противника для стремительного падения.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"trip": {
			Name:        "trip",
			DisplayName: "Зацеп",
			Description: "Точный удар по ноге, нарушающий баланс противника и заставляющий его упасть.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"hip_toss": {
			Name:        "hip_toss",
			DisplayName: "Высед",
			Description: "Бросок через бедро с использованием вращения и силы противника.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"front_sweep": {
			Name:        "front_sweep",
			DisplayName: "Передняя подножка",
			Description: "Быстрый удар по передней ноге, опрокидывающий противника вперёд.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"back_sweep": {
			Name:        "back_sweep",
			DisplayName: "Задняя подножка",
			Description: "Подсечка задней ноги, использующая вес противника для опрокидывания.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"outer_hook": {
			Name:        "outer_hook",
			DisplayName: "Внешний зацеп",
			Description: "Подсечка внешней стороны ноги, нарушающая равновесие противника.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"inner_hook": {
			Name:        "inner_hook",
			DisplayName: "Внутренний зацеп",
			Description: "Подсечка внутренней стороны ноги, выводящая противника из равновесия.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"shoulder_throw": {
			Name:        "shoulder_throw",
			DisplayName: "Бросок через плечо",
			Description: "Рывок противника через плечо с использованием его импульса.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"spinning_throw": {
			Name:        "spinning_throw",
			DisplayName: "Вращательный бросок",
			Description: "Бросок противника через спину с использованием вращения.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"double_leg": {
			Name:        "double_leg",
			DisplayName: "Двойной захват ног",
			Description: "Захват обеих ног противника с последующим опрокидыванием.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"single_leg": {
			Name:        "single_leg",
			DisplayName: "Одинарный захват ноги",
			Description: "Захват одной ноги противника с последующим броском.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"overhook_throw": {
			Name:        "overhook_throw",
			DisplayName: "Бросок через захват",
			Description: "Бросок противника через верхний захват руки.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"underhook_throw": {
			Name:        "underhook_throw",
			DisplayName: "Бросок через подхват",
			Description: "Бросок противника через нижний захват руки.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
	}
	return abilities, nil
}

// GetTeamsConfig возвращает конфигурацию команд
func (m *MockDatabase) GetTeamsConfig() (map[int]TeamConfig, error) {
	return map[int]TeamConfig{
		0: {IconURL: "./static/icons/team0.png"},
		1: {IconURL: "./static/icons/team1.png"},
	}, nil
}

// GetCharacters возвращает список персонажей
func (m *MockDatabase) GetCharacters() ([]Character, error) {
	return []Character{
		{
			ID: 1, Name: "Vasya", Team: 0, HP: 10, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 5, Initiative: 8,
			Weapon: "falchion", Shield: "buckler", Height: 175, Weight: 80, Position: [2]int{1, 0},
			Abilities: []string{"yama_arashi", "outer_hook", "shoulder_throw", "spinning_throw", "double_leg"},
			ImageURL:  "./static/characters/sashya.png",
		},
		{
			ID: 2, Name: "Petya", Team: 0, HP: 10, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 6, Initiative: 7,
			Weapon: "axe", Shield: "shield", Height: 180, Weight: 90, Position: [2]int{1, 2},
			Abilities: []string{"trip", "inner_hook", "shoulder_throw", "single_leg", "overhook_throw"},
			ImageURL:  "./static/characters/suslik.png",
		},
		{
			ID: 3, Name: "Alexei", Team: 0, HP: 10, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9,
			Weapon: "two_handed_sword", Shield: "", Height: 185, Weight: 95, Position: [2]int{4, 4},
			Abilities: []string{"yama_arashi", "trip", "hip_toss", "front_sweep", "back_sweep"},
			ImageURL:  "./static/characters/benya.png",
		},
		{
			ID: 4, Name: "Misha", Team: 0, HP: 10, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6,
			Weapon: "falchion", Shield: "buckler", Height: 170, Weight: 75, Position: [2]int{1, 6},
			Abilities: []string{"hip_toss", "spinning_throw", "double_leg", "underhook_throw", "outer_hook"},
			ImageURL:  "./static/characters/bear_coub.png",
		},
		{
			ID: 5, Name: "Sasha", Team: 0, HP: 10, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8,
			Weapon: "falchion", Shield: "shield", Height: 178, Weight: 85, Position: [2]int{1, 8},
			Abilities: []string{"front_sweep", "overhook_throw", "underhook_throw", "single_leg", "inner_hook"},
			ImageURL:  "./static/characters/kolya.png",
		},
		{
			ID: 6, Name: "Igor", Team: 1, HP: 10, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6,
			Weapon: "falchion", Shield: "buckler", Height: 172, Weight: 78, Position: [2]int{14, 0},
			Abilities: []string{"back_sweep", "outer_hook", "shoulder_throw", "spinning_throw", "double_leg"},
			ImageURL:  "./static/characters/sasha.png",
		},
		{
			ID: 7, Name: "Dima", Team: 1, HP: 10, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8,
			Weapon: "two_handed_halberd", Shield: "", Height: 182, Weight: 92, Position: [2]int{14, 2},
			Abilities: []string{"front_sweep", "inner_hook", "shoulder_throw", "single_leg", "overhook_throw"},
			ImageURL:  "./static/characters/baba.png",
		},
		{
			ID: 8, Name: "Kolya", Team: 1, HP: 10, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 6, Initiative: 7,
			Weapon: "axe", Shield: "shield", Height: 176, Weight: 83, Position: [2]int{12, 4},
			Abilities: []string{"hip_toss", "spinning_throw", "double_leg", "underhook_throw", "outer_hook"},
			ImageURL:  "./static/characters/default.png",
		},
		{
			ID: 9, Name: "Roma", Team: 1, HP: 10, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9,
			Weapon: "sword", Shield: "buckler", Height: 188, Weight: 98, Position: [2]int{14, 6},
			Abilities: []string{"trip", "overhook_throw", "underhook_throw", "single_leg", "inner_hook"},
			ImageURL:  "./static/characters/default.png",
		},
		{
			ID: 10, Name: "Zhenya", Team: 1, HP: 10, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 5, Initiative: 6,
			Weapon: "falchion", Shield: "shield", Height: 174, Weight: 80, Position: [2]int{14, 8},
			Abilities: []string{"yama_arashi", "outer_hook", "shoulder_throw", "spinning_throw", "double_leg"},
			ImageURL:  "./static/characters/default.png",
		},
	}, nil
}

// NewMockDatabase создаёт новый экземпляр MockDatabase
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}
