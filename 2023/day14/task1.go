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
		startRow := 0
		rockCount := 0
		for row := range rockMap {
			switch rockMap[row][col] {
			case '.':
				continue
			case 'O':
				rockCount++
				continue
			case '#':
				endRow := row
				for i := startRow; i < endRow; i++ {
					if rockCount > 0 {
						rockMap[i][col] = 'O'
						rockCount--
					} else {
						rockMap[i][col] = '.'
					}
				}

				startRow = row + 1
			default:
				panic("Invalid map char")
			}
		}

		if rockCount > 0 {
			for i := startRow; i < len(rockMap); i++ {
				if rockCount > 0 {
					rockMap[i][col] = 'O'
					rockCount--
				} else {
					rockMap[i][col] = '.'
				}
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
