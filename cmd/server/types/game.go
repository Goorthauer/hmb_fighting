package types

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func (g *Game) FindCharacter(id int) *Character {
	for i := range g.Teams {
		for j := range g.Teams[i].Characters {
			if g.Teams[i].Characters[j].ID == id {
				return &g.Teams[i].Characters[j]
			}
		}
	}
	return nil
}

// A* для поиска пути на сервере с учётом атак в догонку
func (g *Game) FindPath(startX, startY, endX, endY, stamina int, board [16][9]int, currentCharID int) (path [][2]int, opportunityAttacks []OpportunityAttack) {
	openList := make([]*Node, 0)
	closedList := make(map[string]bool)
	startNode := &Node{X: startX, Y: startY, G: 0, H: heuristic(startX, startY, endX, endY)}
	startNode.F = startNode.G + startNode.H
	openList = append(openList, startNode)

	currentChar := g.FindCharacter(currentCharID)
	if currentChar == nil {
		return nil, nil
	}

	for len(openList) > 0 {
		currentIdx := 0
		for i, node := range openList {
			if node.F < openList[currentIdx].F {
				currentIdx = i
			}
		}
		current := openList[currentIdx]
		openList = append(openList[:currentIdx], openList[currentIdx+1:]...)

		key := fmt.Sprintf("%d,%d", current.X, current.Y)
		if closedList[key] {
			continue
		}
		closedList[key] = true

		if current.X == endX && current.Y == endY {
			path = make([][2]int, 0)
			node := current
			for node != nil {
				path = append([][2]int{[2]int{node.X, node.Y}}, path...)
				node = node.Parent
			}
			if len(path)-1 > stamina {
				return nil, nil
			}

			// Проверка атак в догонку
			opportunityAttacks = g.CheckOpportunityAttacks(currentChar, path)
			return path, opportunityAttacks
		}

		neighbors := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, n := range neighbors {
			newX := current.X + n[0]
			newY := current.Y + n[1]

			if newX < 0 || newX >= 16 || newY < 0 || newY >= 9 {
				continue
			}
			if board[newX][newY] != -1 && (newX != endX || newY != endY) {
				continue // Препятствие, кроме конечной точки
			}

			g := current.G + 1
			if g > stamina {
				continue
			}

			h := heuristic(newX, newY, endX, endY)
			f := g + h
			neighbor := &Node{X: newX, Y: newY, G: g, H: h, F: f, Parent: current}
			neighborKey := fmt.Sprintf("%d,%d", newX, newY)

			if !closedList[neighborKey] {
				openList = append(openList, neighbor)
			}
		}
	}
	return nil, nil // Путь не найден
}

func heuristic(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}
func isInThreatZone(x, y, enemyX, enemyY int) bool {
	return abs(x-enemyX) <= 1 && abs(y-enemyY) <= 1
}

// Проверка атак в догонку
func (g *Game) CheckOpportunityAttacks(target *Character, path [][2]int) []OpportunityAttack {
	var attacks []OpportunityAttack
	startX, startY := path[0][0], path[0][1]
	endX, endY := path[len(path)-1][0], path[len(path)-1][1]

	for i := 0; i < 16; i++ {
		for j := 0; j < 9; j++ {
			if g.Board[i][j] != -1 {
				attacker := g.FindCharacter(g.Board[i][j])
				if attacker != nil && attacker.TeamID != target.TeamID && attacker.HP > 0 {
					startInThreat := isInThreatZone(startX, startY, attacker.Position[0], attacker.Position[1])
					endInThreat := isInThreatZone(endX, endY, attacker.Position[0], attacker.Position[1])
					entersAndExits := false
					for _, p := range path[1 : len(path)-1] { // Проверяем промежуточные точки
						if isInThreatZone(p[0], p[1], attacker.Position[0], attacker.Position[1]) {
							entersAndExits = true
							break
						}
					}

					if (startInThreat && !endInThreat) || (entersAndExits && !endInThreat) {
						pathLength := len(path) - 1
						enemies := g.CountSurroundingEnemies(target)
						wrestlingDiff := attacker.Wrestling - target.Wrestling

						tripChance := 15 + wrestlingDiff*2 + pathLength*3 + enemies*5
						attackChance := 45 + pathLength*2 + enemies*3
						if tripChance < 5 {
							tripChance = 5
						}
						if tripChance > 90 {
							tripChance = 90
						}
						if attackChance > 90-tripChance {
							attackChance = 90 - tripChance
						}

						roll := rand.Intn(100)
						log.Printf("Opportunity Attack by %s on %s: Trip=%d%%, Attack=%d%%, Roll=%d", attacker.Name, target.Name, tripChance, attackChance, roll)

						if roll < tripChance {
							attacks = append(attacks, OpportunityAttack{
								AttackerID: attacker.ID,
								Type:       "trip",
								Damage:     target.HP,
							})
						} else if roll < tripChance+attackChance {
							damage := g.CalculateDamage(attacker, target)
							attacks = append(attacks, OpportunityAttack{
								AttackerID: attacker.ID,
								Type:       "attack",
								Damage:     damage,
							})
						}
					}
				}
			}
		}
	}
	return attacks
}

func (g *Game) DistanceToAbility(pos1, pos2 [2]int) int {
	dx := abs(pos1[0] - pos2[0])
	dy := abs(pos1[1] - pos2[1])
	dist := max(dx, dy)
	log.Printf("Chebyshev Distance from (%d, %d) to (%d, %d) = %d", pos1[0], pos1[1], pos2[0], pos2[1], dist)
	return dist
}

func (g *Game) DistanceToAttack(pos1, pos2 [2]int, weapon Weapon) int {
	dist := max(abs(pos1[0]-pos2[0]), abs(pos1[1]-pos2[1]))
	log.Printf("Attack Distance from (%d, %d) to (%d, %d) = %d, weapon: %s", pos1[0], pos1[1], pos2[0], pos2[1], dist, weapon.Name)
	return dist
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (g *Game) CountSurroundingEnemies(char *Character) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x, y := char.Position[0]+dx, char.Position[1]+dy
			if x >= 0 && x < 16 && y >= 0 && y < 9 && g.Board[x][y] != -1 {
				target := g.FindCharacter(g.Board[x][y])
				if target != nil && target.TeamID != char.TeamID && target.HP > 0 {
					count++
				}
			}
		}
	}
	log.Printf("Surrounding enemies for %s at (%d, %d): %d", char.Name, char.Position[0], char.Position[1], count)
	return count
}

func (g *Game) CalculateDamage(attacker, target *Character) int {
	weapon := g.WeaponsConfig[attacker.Weapon]
	shield := g.ShieldsConfig[attacker.Shield]

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

	surroundingEnemies := g.CountSurroundingEnemies(target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost

	if damage < 0 {
		damage = 0
	}

	log.Printf("Damage calculation: base=%d, attackDefenseDiff=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, attackDefenseDiff, totalDefense, damageBoost, damage)
	return damage
}

func (g *Game) CalculateDamageAfterWrestle(attacker, target *Character) int {
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

	surroundingEnemies := g.CountSurroundingEnemies(target)
	damageBoost := surroundingEnemies * 2
	damage += damageBoost

	if damage < 0 {
		damage = 0
	}

	log.Printf("Damage calculation after wrestle: base=%d, attackDefenseDiff=%d, defense=%d, surroundingBoost=%d, total=%d", baseDamage, attackDefenseDiff, totalDefense, damageBoost, damage)
	return damage
}

func (g *Game) ApplyWrestlingMove(attacker, target *Character, moveName string) {
	var ability Ability
	for abilityID := range g.AbilitiesConfig {
		if abilityID == moveName {
			ability = g.AbilitiesConfig[abilityID]
			break
		}
	}

	if ability.Name == "" {
		log.Printf("Ability %s not found for character %s", moveName, attacker.Name)
		return
	}
	weapon := g.WeaponsConfig[attacker.Weapon]
	shield := g.ShieldsConfig[attacker.Shield]

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
	surroundingEnemies := g.CountSurroundingEnemies(target)
	boostFromSurrounding := surroundingEnemies * 5

	// Добавляем бонусы от оружия и щита
	successChance += mod + boostFromSurrounding + weapon.GrappleBonus + shield.GrappleBonus
	partialSuccessChance += mod + boostFromSurrounding + weapon.GrappleBonus + shield.GrappleBonus
	failureChance -= (mod + boostFromSurrounding) / 2
	totalFailureChance -= (mod + boostFromSurrounding) / 2

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
		g.SetBattleLog(
			fmt.Sprintf("%s успешно применил %s к %s и поверг его!",
				attacker.Name, moveName, target.Name))
		g.Board[target.Position[0]][target.Position[1]] = -1
	case r < successChance+partialSuccessChance:
		damage := g.CalculateDamageAfterWrestle(attacker, target)
		target.HP -= damage
		g.SetBattleLog(
			fmt.Sprintf("%s применил %s  к %s и нанес %d урона!",
				attacker.Name, moveName, target.Name, damage))
		if target.HP <= 0 {
			g.Board[target.Position[0]][target.Position[1]] = -1
		}
	case r < successChance+partialSuccessChance+nothingChance:
		g.SetBattleLog(
			fmt.Sprintf("%s попытался сделать %s на %s и не получилось, видимо плохо подготовил прием!",
				attacker.Name, moveName, target.Name))
	case r < successChance+partialSuccessChance+nothingChance+failureChance:
		attacker.HP = 0
		target.HP = 0
		g.SetBattleLog(
			fmt.Sprintf("%s попытался сделать %s и, %s уже летя вниз утянул его с собой",
				attacker.Name, moveName, target.Name))
		g.Board[attacker.Position[0]][attacker.Position[1]] = -1
		g.Board[target.Position[0]][target.Position[1]] = -1
	default:
		attacker.HP = 0
		g.SetBattleLog(
			fmt.Sprintf("%s попытался сделать %s на %s и, запутавшись в ногах, упал как мешок",
				attacker.Name, moveName, target.Name))
		g.Board[attacker.Position[0]][attacker.Position[1]] = -1
	}
}

func (g *Game) NextTurn() {
	// Если порядок инициативы ещё не установлен, инициализируем его
	if len(g.InitialOrder) == 0 {
		g.InitTurnOrder()
	}

	// Проверяем, остались ли живые команды
	aliveTeamsOne := 0
	aliveTeamsTwo := 0
	for _, team := range g.Teams {
		for _, char := range team.Characters {
			if char.HP > 0 {
				if char.TeamID == 0 {
					aliveTeamsOne++
				} else if char.TeamID == 1 {
					aliveTeamsTwo++
				}
			}
		}
	}
	if aliveTeamsOne < 1 {
		g.Winner = 1
		g.Phase = "finished"
		log.Printf("Team 1 wins!")
		return
	}
	if aliveTeamsTwo < 1 {
		g.Winner = 0
		g.Phase = "finished"
		log.Printf("Team 0 wins!")
		return
	}

	// Находим текущий индекс в порядке инициативы
	currentIndex := -1
	for i, id := range g.InitialOrder {
		if id == g.CurrentTurn {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 { // Если текущего хода нет в списке (например, первый ход)
		currentIndex = len(g.InitialOrder) - 1 // Устанавливаем в конец, чтобы начать с начала
	}

	// Ищем следующего живого персонажа
	nextIndex := (currentIndex + 1) % len(g.InitialOrder)
	startIndex := nextIndex
	for {
		nextChar := g.FindCharacter(g.InitialOrder[nextIndex])
		if nextChar != nil && nextChar.HP > 0 {
			g.CurrentTurn = nextChar.ID
			g.Phase = "move"
			log.Printf("Next turn: %s (ID: %d)", nextChar.Name, nextChar.ID)
			break
		}
		nextIndex = (nextIndex + 1) % len(g.InitialOrder)
		if nextIndex == startIndex { // Если обошли весь круг и никого не нашли
			g.Phase = "finished"
			log.Printf("No alive characters left!")
			return
		}
	}

	// Обновляем эффекты
	for i := range g.Teams {
		for j := range g.Teams[i].Characters {
			char := &g.Teams[i].Characters[j]
			for k := len(char.Effects) - 1; k >= 0; k-- {
				char.Effects[k].Duration--
				if char.Effects[k].Duration <= 0 {
					char.Effects = append(char.Effects[:k], char.Effects[k+1:]...)
				}
			}
		}
	}
}

func (g *Game) SetBattleLog(action string) {
	bl := Battlelog{
		Time:   time.Now().Format(time.TimeOnly),
		Action: action,
	}
	g.Battlelog = append(g.Battlelog, bl)
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

// Инициализация порядка хода в начале игры (вызывать один раз)
func (g *Game) InitTurnOrder() {
	liveChars := []Character{}
	for _, team := range g.Teams {
		for _, char := range team.Characters {
			liveChars = append(liveChars, char)
		}
	}
	sortCharactersByInitiative(liveChars)
	g.InitialOrder = make([]int, len(liveChars))
	for i, char := range liveChars {
		g.InitialOrder[i] = char.ID
	}
	log.Printf("Initial turn order set: %v", g.InitialOrder)
}
