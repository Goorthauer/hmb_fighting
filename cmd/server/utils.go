package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

func distanceToAbility(pos1, pos2 [2]int) int {
	dx := abs(pos1[0] - pos2[0])
	dy := abs(pos1[1] - pos2[1])
	dist := max(dx, dy)
	log.Printf("Chebyshev Distance from (%d, %d) to (%d, %d) = %d", pos1[0], pos1[1], pos2[0], pos2[1], dist)
	return dist
}

func distanceToAttack(pos1, pos2 [2]int, weapon Weapon) int {
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

// A* для поиска пути на сервере с учётом атак в догонку
func findPath(startX, startY, endX, endY, stamina int, board [16][9]int, currentCharID int, game *Game) (path [][2]int, opportunityAttacks []OpportunityAttack) {
	openList := make([]*Node, 0)
	closedList := make(map[string]bool)
	startNode := &Node{X: startX, Y: startY, G: 0, H: heuristic(startX, startY, endX, endY)}
	startNode.F = startNode.G + startNode.H
	openList = append(openList, startNode)

	currentChar := findCharacter(game, currentCharID)
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
			opportunityAttacks = checkOpportunityAttacks(game, currentChar, path)
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

var mutex = sync.Mutex{}

// Проверка атак в догонку
func checkOpportunityAttacks(game *Game, target *Character, path [][2]int) []OpportunityAttack {
	var attacks []OpportunityAttack
	startX, startY := path[0][0], path[0][1]
	endX, endY := path[len(path)-1][0], path[len(path)-1][1]

	for i := 0; i < 16; i++ {
		for j := 0; j < 9; j++ {
			if game.Board[i][j] != -1 {
				attacker := findCharacter(game, game.Board[i][j])
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
						enemies := countSurroundingEnemies(game, target)
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
							damage := calculateDamage(attacker, target, game)
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

func isInThreatZone(x, y, enemyX, enemyY int) bool {
	return abs(x-enemyX) <= 1 && abs(y-enemyY) <= 1
}
