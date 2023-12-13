package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func solvePuzzle(puzzle []string) int {
	var mirrorRow, mirrorCol int

	// Check horizontal
	for mirror := 1; mirror < len(puzzle); mirror++ {
		match := true
		for distance := 0; match; distance++ {
			left, right := mirror-distance-1, mirror+distance
			if left < 0 || right >= len(puzzle) {
				break
			}
			if puzzle[left] != puzzle[right] {
				match = false
			}
		}
		if match {
			mirrorRow = mirror
			break
		}
	}

	// Check vertical
	for mirror := 1; mirror < len(puzzle[0]); mirror++ {
		match := true
		for distance := 0; match; distance++ {
			left, right := mirror-distance-1, mirror+distance
			if left < 0 || right >= len(puzzle[0]) {
				break
			}
			for row := range puzzle {
				if puzzle[row][left] != puzzle[row][right] {
					match = false
					break
				}
			}
		}
		if match {
			mirrorCol = mirror
			break
		}
	}

	return mirrorRow*100 + mirrorCol
}

func main() {
	input := util.StdinReadlines()

	sum, start := 0, 0
	for i, line := range input {
		if len(line) == 0 {
			sum += solvePuzzle(input[start:i])
			start = i + 1
		}
	}
	sum += solvePuzzle(input[start:])
	fmt.Println(sum)
}
