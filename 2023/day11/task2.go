package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

type loc struct {
	row, col int
}

func parseGalaxies(input []string) (map[loc]int, []bool, []bool) {
	galaxyId := 1
	galaxies := make(map[loc]int)

	rowOccupied := make([]bool, len(input))
	colOccupied := make([]bool, len(input[0]))

	for row, line := range input {
		for col, r := range line {
			if r == '#' {
				galaxies[loc{row, col}] = galaxyId
				galaxyId++

				rowOccupied[row] = true
				colOccupied[col] = true
			}
		}
	}

	return galaxies, rowOccupied, colOccupied
}

func distance(val1, val2 int, occupied []bool) int {
	if val1 > val2 {
		val1, val2 = val2, val1
	}

	steps := 0
	for i := val1 + 1; i <= val2; i++ {
		if occupied[i] {
			steps++
		} else {
			steps += 1000000
		}
	}
	return steps
}

func distanceSum(galaxies map[loc]int, rowOccupied, colOccupied []bool) int {
	sum := 0
	for l1, galaxy1 := range galaxies {
		for l2, galaxy2 := range galaxies {
			if galaxy1 > galaxy2 {
				continue
			}

			sum += distance(l1.row, l2.row, rowOccupied)
			sum += distance(l1.col, l2.col, colOccupied)
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	galaxies, rowOccupied, colOccupied := parseGalaxies(input)
	fmt.Println(distanceSum(galaxies, rowOccupied, colOccupied))
}
