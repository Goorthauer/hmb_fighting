package entities

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID       string `json:"ID"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

type Team struct {
	Characters []Character `json:"characters"`
}

type Client struct {
	Conn      *websocket.Conn
	TeamID    int
	ClientID  string
	Spectator bool
	User      *User
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
