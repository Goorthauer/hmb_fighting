package main

import (
	"log"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var mutex = sync.Mutex{}
var rooms = make(map[string]*Game)
var users = make(map[string]User)

func initGame() *Game {
	weaponsConfig := map[string]Weapon{
		"falchion":           {Name: "falchion", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png"},
		"axe":                {Name: "axe", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png"},
		"two_handed_sword":   {Name: "two_handed_sword", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png"},
		"spear":              {Name: "spear", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png"},
		"dagger":             {Name: "dagger", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png"},
		"two_handed_halberd": {Name: "two_handed_halberd", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png"},
		"sword":              {Name: "sword", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png"},
	}

	shieldsConfig := map[string]Shield{
		"buckler": {Name: "buckler", DefenseBonus: 3, ImageURL: "./static/shields/default.png"},
		"shield":  {Name: "shield", DefenseBonus: 5, ImageURL: "./static/shields/default.png"},
		"":        {Name: "none", DefenseBonus: 0, ImageURL: ""},
	}

	teamsConfig := [2]TeamConfig{
		{IconURL: "./static/icons/default.png"},
		{IconURL: "./static/icons/default.png"},
	}

	game := &Game{
		Connections:   make(map[*websocket.Conn]*Client),
		GameSessionId: uuid.New().String(),
		WeaponsConfig: weaponsConfig,
		ShieldsConfig: shieldsConfig,
		TeamsConfig:   teamsConfig,
		Teams: [2]Team{
			{Characters: []Character{
				{ID: 1, Name: "Vasya", Team: 0, HP: 100, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 5, Initiative: 8, Weapon: "falchion", Shield: "buckler", Height: 175, Weight: 80, Position: [2]int{2, 2}, Abilities: []Ability{{Name: "Takedown", Type: "wrestle", Description: "Attempts to take down the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 2, Name: "Petya", Team: 0, HP: 100, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 180, Weight: 90, Position: [2]int{2, 3}, Abilities: []Ability{{Name: "Throw", Type: "wrestle", Description: "Throws the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 3, Name: "Alexei", Team: 0, HP: 100, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "two_handed_sword", Shield: "", Height: 185, Weight: 95, Position: [2]int{7, 7}, Abilities: []Ability{{Name: "Pin", Type: "wrestle", Description: "Pins the opponent down", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 4, Name: "Misha", Team: 0, HP: 100, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "spear", Shield: "buckler", Height: 170, Weight: 75, Position: [2]int{3, 4}, Abilities: []Ability{{Name: "Grapple", Type: "wrestle", Description: "Grapples the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 5, Name: "Sasha", Team: 0, HP: 100, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "dagger", Shield: "shield", Height: 178, Weight: 85, Position: [2]int{4, 5}, Abilities: []Ability{{Name: "Lock", Type: "wrestle", Description: "Locks the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
			}},
			{Characters: []Character{
				{ID: 6, Name: "Igor", Team: 1, HP: 100, Stamina: 5, AttackMin: 9, AttackMax: 19, Defense: 5, Initiative: 6, Weapon: "falchion", Shield: "buckler", Height: 172, Weight: 78, Position: [2]int{17, 2}, Abilities: []Ability{{Name: "Takedown", Type: "wrestle", Description: "Attempts to take down the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 7, Name: "Dima", Team: 1, HP: 100, Stamina: 6, AttackMin: 11, AttackMax: 21, Defense: 7, Initiative: 8, Weapon: "two_handed_halberd", Shield: "", Height: 182, Weight: 92, Position: [2]int{17, 3}, Abilities: []Ability{{Name: "Throw", Type: "wrestle", Description: "Throws the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 8, Name: "Kolya", Team: 1, HP: 100, Stamina: 5, AttackMin: 10, AttackMax: 20, Defense: 6, Initiative: 7, Weapon: "axe", Shield: "shield", Height: 176, Weight: 83, Position: [2]int{16, 4}, Abilities: []Ability{{Name: "Pin", Type: "wrestle", Description: "Pins the opponent down", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 9, Name: "Roma", Team: 1, HP: 100, Stamina: 7, AttackMin: 12, AttackMax: 22, Defense: 4, Initiative: 9, Weapon: "sword", Shield: "buckler", Height: 188, Weight: 98, Position: [2]int{15, 5}, Abilities: []Ability{{Name: "Grapple", Type: "wrestle", Description: "Grapples the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
				{ID: 10, Name: "Zhenya", Team: 1, HP: 100, Stamina: 6, AttackMin: 8, AttackMax: 18, Defense: 5, Initiative: 6, Weapon: "dagger", Shield: "shield", Height: 174, Weight: 80, Position: [2]int{14, 6}, Abilities: []Ability{{Name: "Lock", Type: "wrestle", Description: "Locks the opponent", Range: 1, ImageURL: "./static/abilities/default.jpg"}}, ImageURL: "./static/characters/default.png"},
			}},
		},
		CurrentTurn: 3,
		Phase:       "move",
		Board:       [20][10]int{},
	}

	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}
	for _, team := range game.Teams {
		for _, char := range team.Characters {
			if shield, ok := game.ShieldsConfig[char.Shield]; ok {
				char.Defense += shield.DefenseBonus
			}
			game.Board[char.Position[0]][char.Position[1]] = char.ID
		}
	}
	return game
}

func findCharacter(game *Game, id int) *Character {
	for i := range game.Teams {
		for j := range game.Teams[i].Characters {
			if game.Teams[i].Characters[j].ID == id {
				return &game.Teams[i].Characters[j]
			}
		}
	}
	return nil
}

func countSurroundingEnemies(game *Game, char *Character) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x, y := char.Position[0]+dx, char.Position[1]+dy
			if x >= 0 && x < 20 && y >= 0 && y < 10 && game.Board[x][y] != -1 {
				target := findCharacter(game, game.Board[x][y])
				if target != nil && target.Team != char.Team && target.HP > 0 {
					count++
				}
			}
		}
	}
	log.Printf("Surrounding enemies for %s at (%d, %d): %d", char.Name, char.Position[0], char.Position[1], count)
	return count
}

func calculateDamage(attacker, target *Character, game *Game) int {
	baseDamage := rand.Intn(attacker.AttackMax-attacker.AttackMin+1) + attacker.AttackMin
	totalDefense := target.Defense
	for _, effect := range target.Effects {
		totalDefense += effect.DefenseMod
	}
	damage := baseDamage - totalDefense
	if damage < 0 {
		damage = 0
	}
	surroundingEnemies := countSurroundingEnemies(game, target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost
	if damage < 0 {
		damage = 0
	}
	log.Printf("Damage calculation: base=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, totalDefense, damageBoost, damage)
	return damage
}

func applyWrestlingMove(game *Game, attacker, target *Character, moveName string) {
	successChance := 20
	partialSuccessChance := 25
	nothingChance := 30
	failureChance := 10
	totalFailureChance := 5

	heightDiff := float64(attacker.Height-target.Height) / 10.0
	weightDiff := float64(attacker.Weight-target.Weight) / 10.0
	mod := int(heightDiff+weightDiff) * 5

	surroundingEnemies := countSurroundingEnemies(game, target)
	successBoost := surroundingEnemies * 5
	successChance += mod + successBoost
	partialSuccessChance += mod + successBoost
	failureChance -= mod / 2
	totalFailureChance -= mod / 2

	if successChance < 5 {
		successChance = 5
	}
	if partialSuccessChance < 5 {
		partialSuccessChance = 5
	}
	if failureChance < 5 {
		failureChance = 5
	}
	if totalFailureChance < 5 {
		totalFailureChance = 5
	}

	total := successChance + partialSuccessChance + nothingChance + failureChance + totalFailureChance
	if total != 100 {
		scale := float64(100) / float64(total)
		successChance = int(float64(successChance) * scale)
		nothingChance = int(float64(nothingChance) * scale)
		partialSuccessChance = int(float64(partialSuccessChance) * scale)
		failureChance = int(float64(failureChance) * scale)
		totalFailureChance = 100 - successChance - partialSuccessChance - failureChance - nothingChance
	}

	r := rand.Intn(100)
	log.Printf("%s attempts %s on %s: success=%d%%, partial=%d%%, failure=%d%%, totalFailure=%d%%, surroundingBoost=%d, roll=%d",
		attacker.Name, moveName, target.Name, successChance, partialSuccessChance, failureChance, totalFailureChance, successBoost, r)

	switch {
	case r < successChance:
		target.HP = 0
		log.Printf("%s successfully used %s on %s, knocking them out!", attacker.Name, moveName, target.Name)
		game.Board[target.Position[0]][target.Position[1]] = -1
	case r < successChance+partialSuccessChance:
		damage := calculateDamage(attacker, target, game)
		target.HP -= damage
		log.Printf("%s partially succeeded with %s on %s, dealing %d damage", attacker.Name, moveName, target.Name, damage)
		if target.HP <= 0 {
			game.Board[target.Position[0]][target.Position[1]] = -1
		}
	case r < successChance+partialSuccessChance+nothingChance:
		log.Printf("%s failed to use %s on %s - the move didn't connect!", attacker.Name, moveName, target.Name)
	case r < successChance+partialSuccessChance+nothingChance+failureChance:
		attacker.HP = 0
		target.HP = 0
		log.Printf("%s failed %s - both %s and %s are knocked out!", attacker.Name, moveName, attacker.Name, target.Name)
		game.Board[attacker.Position[0]][attacker.Position[1]] = -1
		game.Board[target.Position[0]][target.Position[1]] = -1
	default:
		attacker.HP = 0
		log.Printf("%s catastrophically failed %s - %s is knocked out!", attacker.Name, moveName, attacker.Name)
		game.Board[attacker.Position[0]][attacker.Position[1]] = -1
	}
}

func nextTurn(game *Game) {
	liveChars := []Character{}
	for _, team := range game.Teams {
		for _, char := range team.Characters {
			if char.HP > 0 {
				liveChars = append(liveChars, char)
			}
		}
	}
	if len(liveChars) == 0 {
		return
	}
	sortCharactersByInitiative(liveChars)
	currentIndex := -1
	for i, char := range liveChars {
		if char.ID == game.CurrentTurn {
			currentIndex = i
			break
		}
	}
	nextIndex := (currentIndex + 1) % len(liveChars)
	game.CurrentTurn = liveChars[nextIndex].ID
	game.Phase = "move"
	for i := range game.Teams {
		for j := range game.Teams[i].Characters {
			char := &game.Teams[i].Characters[j]
			for k := len(char.Effects) - 1; k >= 0; k-- {
				char.Effects[k].Duration--
				if char.Effects[k].Duration <= 0 {
					char.Effects = append(char.Effects[:k], char.Effects[k+1:]...)
				}
			}
		}
	}
}

func sortCharactersByInitiative(chars []Character) {
	for i := 0; i < len(chars)-1; i++ {
		for j := 0; j < len(chars)-i-1; j++ {
			if chars[j].Initiative < chars[j+1].Initiative {
				chars[j], chars[j+1] = chars[j+1], chars[j]
			}
		}
	}
}
