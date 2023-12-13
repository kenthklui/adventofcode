package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func solvePuzzle(puzzle []string) int {
	// Check horizontal
	for mirror := 1; mirror < len(puzzle); mirror++ {
		mismatch := 0
		for distance := 0; distance < len(puzzle); distance++ {
			left, right := mirror-distance-1, mirror+distance
			if left < 0 || right >= len(puzzle) {
				break
			}

			for row := range puzzle[left] {
				if puzzle[left][row] != puzzle[right][row] {
					mismatch++
				}
			}
			if mismatch > 1 {
				break
			}
		}
		if mismatch == 1 {
			return mirror * 100
		}
	}

	// Check vertical
	for mirror := 1; mirror < len(puzzle[0]); mirror++ {
		mismatch := 0
		for distance := 0; distance < len(puzzle[0]); distance++ {
			left, right := mirror-distance-1, mirror+distance
			if left < 0 || right >= len(puzzle[0]) {
				break
			}

			for row := range puzzle {
				if puzzle[row][left] != puzzle[row][right] {
					mismatch++
				}
				if mismatch > 1 {
					break
				}
			}
			if mismatch > 1 {
				break
			}
		}
		if mismatch == 1 {
			return mirror
		}
	}

	panic("Missing new reflection line")
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
