package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

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

type TeamConfig struct {
	IconURL     string `json:"iconURL"`
	Name        string `json:"name"`
	ID          int    `json:"ID"`
	Description string `json:"description"`
}

type Role struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}

type Character struct {
	//base
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TeamID         int    `json:"team"`
	RoleID         int    `json:"role"`
	CountOfAbility int    `json:"-"`
	ImageURL       string `json:"imageURL"`
	IsActive       bool   `json:"isActive"`
	//заполняются в бою или перед инициализацией
	Abilities []string `json:"abilities"`
	Effects   []Effect `json:"effects"`
	Position  [2]int   `json:"position"`
	// Снаряжение
	Weapon        string `json:"weapon"`
	Shield        string `json:"shield"`
	IsTitanArmour bool   `json:"IsTitanArmour"` //новая, показывает из титана комплект доспехов у человека или нет.
	// Основные характеристики
	Height     int `json:"height"`
	Weight     int `json:"weight"`
	HP         int `json:"hp"`
	Stamina    int `json:"stamina"`
	Initiative int `json:"initiative"`
	Wrestling  int `json:"wrestling"` //новая, показатель борьбы. Если показатель у персонажа показатель высокий, а у противника низкий, то шанс успеха приема большой, если наоборот - шанс неуспеха приема большой.
	Attack     int `json:"attack"`    //новая, показатель атаки. Нужен для невилирования защиты противника.
	Defense    int `json:"defense"`
	// остальное
	AttackMin int `json:"attackMin"`
	AttackMax int `json:"attackMax"`
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

type GameState struct {
	Teams           [2]Team            `json:"teams"`
	Winner          int                `json:"winner"`
	CurrentTurn     int                `json:"currentTurn"`
	Phase           string             `json:"phase"`
	Board           [16][9]int         `json:"board"` // Обновляем с [20][10] на [16][9]
	TeamID          int                `json:"teamID"`
	ClientID        string             `json:"clientID"`
	GameSessionId   string             `json:"gameSessionId"`
	WeaponsConfig   map[string]Weapon  `json:"weaponsConfig"`
	AbilitiesConfig map[string]Ability `json:"abilitiesConfig"`
	ShieldsConfig   map[string]Shield  `json:"shieldsConfig"`
	TeamsConfig     [2]TeamConfig      `json:"teamsConfig"`
	Battlelog       []string           `json:"battlelog"`
}

type Team struct {
	Characters []Character `json:"characters"`
}

type Action struct {
	Type        string `json:"type"`
	CharacterID int    `json:"characterID"`
	Position    [2]int `json:"position"`
	TargetID    int    `json:"targetID"`
	Ability     string `json:"ability"`
	ClientID    string `json:"clientID"`
}

type Client struct {
	Conn      *websocket.Conn
	TeamID    int
	ClientID  string
	Spectator bool
	User      *User
}

type Game struct {
	Connections     map[*websocket.Conn]*Client
	Teams           map[int]Team   // Теперь map для гибкости
	Players         map[int]string // teamID -> clientID (только 2 игрока)
	CurrentTurn     int
	Phase           string
	Board           [16][9]int
	GameSessionId   string
	WeaponsConfig   map[string]Weapon
	ShieldsConfig   map[string]Shield
	AbilitiesConfig map[string]Ability
	RoleConfig      map[string]Role
	TeamsConfig     map[int]TeamConfig // Теперь map
	Mutex           sync.Mutex
	InitialOrder    []int
	Battlelog       []string
	Winner          int // ID команды-победителя, -1 если нет
}

// Node для алгоритма A*
type Node struct {
	X      int
	Y      int
	G      int
	H      int
	F      int
	Parent *Node
}

// Структура для результата атаки в догонку
type OpportunityAttack struct {
	AttackerID int
	Type       string // "trip" или "attack"
	Damage     int
}
