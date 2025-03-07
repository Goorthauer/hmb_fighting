package main

import "log"

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
