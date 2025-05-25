package entities

import (
	"hmb_fighting/server/types"
)

type Weapon struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Range        int    `json:"range"`
	IsTwoHanded  bool   `json:"isTwoHanded"`
	ImageURL     string `json:"imageURL"`
	AttackBonus  int    `json:"attackBonus"`  // Бонус к атаке
	GrappleBonus int    `json:"grappleBonus"` // Бонус к успешным состояниям борьбы
}

type Shield struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	DefenseBonus int    `json:"defenseBonus"`
	ImageURL     string `json:"imageURL"`
	AttackBonus  int    `json:"attackBonus"`  // Бонус к атаке
	GrappleBonus int    `json:"grappleBonus"` // Бонус к успешным состояниям борьбы
}

type Ability struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Range       int    `json:"range"`
	ImageURL    string `json:"imageURL"`
}

type Effect struct {
	Name       string `json:"name"`
	Duration   int    `json:"duration"`
	StaminaMod int    `json:"staminaMod"`
	AttackMod  int    `json:"attackMod"`
	DefenseMod int    `json:"defenseMod"`
}

type Action struct {
	Type        types.ActionTypes `json:"type"`
	CharacterID int               `json:"characterID"`
	Position    [2]int            `json:"position"`
	TargetID    int               `json:"targetID"`
	Ability     string            `json:"ability"`
	ClientID    string            `json:"clientID"`
}

// OpportunityAttack Структура для результата атаки вдогонку.
type OpportunityAttack struct {
	AttackerID int
	Type       types.OpportunityAttackTypes
	Damage     int
}

type Battlelog struct {
	Time   string
	Action string
}
