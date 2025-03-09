package main

import (
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func initGame(db Database) *Game {
	weaponsConfig, err := db.GetWeapons()
	if err != nil {
		log.Fatalf("Failed to get weapons config: %v", err)
	}

	shieldsConfig, err := db.GetShields()
	if err != nil {
		log.Fatalf("Failed to get shields config: %v", err)
	}

	teamsConfig, err := db.GetTeamsConfig()
	if err != nil {
		log.Fatalf("Failed to get teams config: %v", err)
	}

	characters, err := db.GetCharacters()
	if err != nil {
		log.Fatalf("Failed to get characters: %v", err)
	}

	abilitiesConfig, err := db.GetAbilities()
	if err != nil {
		log.Fatalf("Failed to get abilities config: %v", err)
	}

	firstCharactersTeam := make([]Character, 0)
	secondCharactersTeam := make([]Character, 0)
	for _, char := range characters {
		char.SetAbilities(abilitiesConfig)
		char.Position = [2]int{-1, -1}
		// Применяем эффекты Titan Armour при инициализации
		if char.IsTitanArmour {
			char.Wrestling += 1
			char.Stamina += 1
			char.Initiative += 1
			char.Defense -= 2
			char.HP -= 5
			if char.HP < 1 {
				char.HP = 1 // Минимальное значение HP
			}
			if char.Defense < 0 {
				char.Defense = 0 // Минимальное значение защиты
			}
		}
		if char.Team == 0 {
			firstCharactersTeam = append(firstCharactersTeam, char)
		} else if char.Team == 1 {
			secondCharactersTeam = append(secondCharactersTeam, char)
		}
	}

	teams := map[int]Team{
		0: {Characters: firstCharactersTeam},
		1: {Characters: secondCharactersTeam},
	}

	game := &Game{
		Connections:     make(map[*websocket.Conn]*Client),
		GameSessionId:   uuid.New().String(),
		WeaponsConfig:   weaponsConfig,
		ShieldsConfig:   shieldsConfig,
		TeamsConfig:     teamsConfig,
		Teams:           teams,
		AbilitiesConfig: abilitiesConfig,
		CurrentTurn:     -1,
		Phase:           "setup",
		Players:         make(map[int]string),
		Board:           [16][9]int{},
	}

	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = -1
		}
	}

	for _, team := range game.Teams {
		for i := range team.Characters {
			char := &team.Characters[i]
			if shield, ok := game.ShieldsConfig[char.Shield]; ok {
				char.Defense += shield.DefenseBonus
				char.AttackMin += shield.AttackBonus
				char.AttackMax += shield.AttackBonus
			}
			if weapon, ok := game.WeaponsConfig[char.Weapon]; ok {
				char.AttackMin += weapon.AttackBonus
				char.AttackMax += weapon.AttackBonus
			}
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
			if x >= 0 && x < 16 && y >= 0 && y < 9 && game.Board[x][y] != -1 {
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
	weapon := game.WeaponsConfig[attacker.Weapon]
	shield := game.ShieldsConfig[attacker.Shield]

	baseDamage := rand.Intn(attacker.AttackMax-attacker.AttackMin+1) + attacker.AttackMin
	totalDefense := target.Defense
	for _, effect := range target.Effects {
		totalDefense += effect.DefenseMod
	}

	// Учитываем разницу между Attack и Defense
	attackDefenseDiff := attacker.Attack - totalDefense
	baseDamage += attackDefenseDiff

	// Добавляем бонусы от оружия и щита
	baseDamage += weapon.AttackBonus + shield.AttackBonus

	damage := baseDamage
	if damage < 0 {
		damage = 0
	}

	surroundingEnemies := countSurroundingEnemies(game, target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost

	if damage < 0 {
		damage = 0
	}

	log.Printf("Damage calculation: base=%d, attackDefenseDiff=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, attackDefenseDiff, totalDefense, damageBoost, damage)
	return damage
}

func calculateDamageAfterWrestle(attacker, target *Character, game *Game) int {
	baseDamage := rand.Intn(attacker.AttackMax-attacker.AttackMin+1) + attacker.AttackMin
	totalDefense := target.Defense
	for _, effect := range target.Effects {
		totalDefense += effect.DefenseMod
	}

	// Учитываем разницу между Attack и Defense
	attackDefenseDiff := attacker.Attack - totalDefense
	baseDamage += attackDefenseDiff

	damage := baseDamage
	if damage < 0 {
		damage = 0
	}

	surroundingEnemies := countSurroundingEnemies(game, target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost

	if damage < 0 {
		damage = 0
	}

	log.Printf("Damage calculation after wrestle: base=%d, attackDefenseDiff=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, attackDefenseDiff, totalDefense, damageBoost, damage)
	return damage
}

func applyWrestlingMove(game *Game, attacker, target *Character, moveName string) {
	var ability Ability
	for abilityID := range game.AbilitiesConfig {
		if abilityID == moveName {
			ability = game.AbilitiesConfig[abilityID]
			break
		}
	}

	if ability.Name == "" {
		log.Printf("Ability %s not found for character %s", moveName, attacker.Name)
		return
	}
	weapon := game.WeaponsConfig[attacker.Weapon]
	shield := game.ShieldsConfig[attacker.Shield]

	// Новая формула с учетом Wrestling
	wrestlingDiff := attacker.Wrestling - target.Wrestling
	successChance := 25 + wrestlingDiff*5 // Базовый шанс успеха увеличивается или уменьшается на 5% за каждый пункт разницы
	partialSuccessChance := 25
	nothingChance := 25
	failureChance := 15 - wrestlingDiff*2      // Уменьшаем шанс провала при высокой разнице
	totalFailureChance := 10 - wrestlingDiff*3 // Уменьшаем шанс полного провала

	// Модификаторы от роста и веса
	heightDiff := float64(attacker.Height-target.Height) / 10.0
	weightDiff := float64(attacker.Weight-target.Weight) / 10.0
	mod := int(heightDiff+weightDiff) * 5

	// Учитываем окружающих врагов
	surroundingEnemies := countSurroundingEnemies(game, target)
	successBoost := surroundingEnemies * 5

	// Добавляем бонусы от оружия и щита
	successChance += mod + successBoost + weapon.GrappleBonus + shield.GrappleBonus
	partialSuccessChance += mod + successBoost + weapon.GrappleBonus + shield.GrappleBonus
	failureChance -= mod / 2
	totalFailureChance -= mod / 2

	// Ограничиваем минимальные и максимальные значения
	if successChance < 5 {
		successChance = 5
	}
	if successChance > 90 {
		successChance = 90
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

	// Нормализуем проценты до 100
	total := successChance + partialSuccessChance + nothingChance + failureChance + totalFailureChance
	if total != 100 {
		scale := float64(100) / float64(total)
		successChance = int(float64(successChance) * scale)
		partialSuccessChance = int(float64(partialSuccessChance) * scale)
		nothingChance = int(float64(nothingChance) * scale)
		failureChance = int(float64(failureChance) * scale)
		totalFailureChance = 100 - successChance - partialSuccessChance - nothingChance - failureChance
	}

	r := rand.Intn(100)
	log.Printf("%s attempts %s on %s: success=%d%%, partial=%d%%, nothing=%d%%, failure=%d%%, totalFailure=%d%%, wrestlingDiff=%d, roll=%d",
		attacker.Name, moveName, target.Name, successChance, partialSuccessChance, nothingChance, failureChance, totalFailureChance, wrestlingDiff, r)

	switch {
	case r < successChance:
		target.HP = 0
		log.Printf("%s successfully used %s on %s, knocking them out!", attacker.Name, moveName, target.Name)
		game.Board[target.Position[0]][target.Position[1]] = -1
	case r < successChance+partialSuccessChance:
		damage := calculateDamageAfterWrestle(attacker, target, game)
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
	aliveTeams := 0
	winner := -1
	for teamID, team := range game.Teams {
		for _, char := range team.Characters {
			if char.HP > 0 {
				aliveTeams++
				winner = teamID
				break
			}
		}
	}
	if aliveTeams <= 1 {
		game.Winner = winner
		game.Phase = "finished"
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
