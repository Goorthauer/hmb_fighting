package main

import (
	"fmt"
	"log"
)

func distance(pos1, pos2 [2]int) int {
	dist := abs(pos1[0]-pos2[0]) + abs(pos1[1]-pos2[1])
	log.Printf("Distance from (%d, %d) to (%d, %d) = %d", pos1[0], pos1[1], pos2[0], pos2[1], dist)
	return dist
}

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

// A* для поиска пути на сервере
func findPath(startX, startY, endX, endY, stamina int, board [16][9]int, currentCharID int) [][2]int {
	openList := make([]*Node, 0)
	closedList := make(map[string]bool)
	startNode := &Node{X: startX, Y: startY, G: 0, H: heuristic(startX, startY, endX, endY)}
	startNode.F = startNode.G + startNode.H
	openList = append(openList, startNode)

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
			path := make([][2]int, 0)
			node := current
			for node != nil {
				path = append([][2]int{[2]int{node.X, node.Y}}, path...)
				node = node.Parent
			}
			if len(path)-1 > stamina {
				return nil
			}
			return path
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
	return nil // Путь не найден
}

func heuristic(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}
