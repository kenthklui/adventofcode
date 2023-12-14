package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func makeMap(input []string) [][]rune {
	rockMap := make([][]rune, len(input))
	for i, line := range input {
		rockMap[i] = make([]rune, len(line))
		for j, c := range line {
			rockMap[i][j] = c
		}
	}
	return rockMap
}

func rollNorth(rockMap [][]rune) [][]rune {
	for col := range rockMap[0] {
		openSpot := -1
		for row := range rockMap {
			switch rockMap[row][col] {
			case '.':
				if openSpot < 0 {
					openSpot = row
				}
			case 'O':
				if openSpot >= 0 {
					rockMap[openSpot][col], rockMap[row][col] = 'O', '.'
					openSpot++
				}
			case '#':
				openSpot = -1
			default:
				panic("Invalid map char")
			}
		}
	}

	return rockMap
}

func load(rockMap [][]rune) int {
	sum := 0
	for rowNum, row := range rockMap {
		for _, c := range row {
			if c == 'O' {
				sum += len(rockMap) - rowNum
			}
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	rockMap := makeMap(input)
	rolled := rollNorth(rockMap)
	fmt.Println(load(rolled))
}
