package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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
	IconURL string `json:"iconURL"`
}

type Character struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Team       int      `json:"team"`
	HP         int      `json:"hp"`
	Stamina    int      `json:"stamina"`
	AttackMin  int      `json:"attackMin"`
	AttackMax  int      `json:"attackMax"`
	Defense    int      `json:"defense"`
	Initiative int      `json:"initiative"`
	Weapon     string   `json:"weapon"`
	Shield     string   `json:"shield"`
	Height     int      `json:"height"`
	Weight     int      `json:"weight"`
	Position   [2]int   `json:"position"`
	Abilities  []string `json:"abilities"` // Теперь это список идентификаторов способностей
	Effects    []Effect `json:"effects"`
	ImageURL   string   `json:"imageURL"`
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
	TeamsConfig     map[int]TeamConfig // Теперь map
	mutex           sync.Mutex
	SetupPhase      bool // Фаза расстановки
	Winner          int  // ID команды-победителя, -1 если нет
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
