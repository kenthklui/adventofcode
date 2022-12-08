package main

import (
	"bufio"
	"fmt"
	"os"
)

func readInput() []string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

func parseInput(input []string) [][]int {
	trees := make([][]int, len(input))

	for i, line := range input {
		trees[i] = make([]int, len(line))
		for j, r := range line {
			trees[i][j] = int(r - '0')
		}
	}

	return trees
}

func countVisible(trees [][]int) int {
	visible := make([][]bool, len(trees))
	for i, line := range trees {
		visible[i] = make([]bool, len(line))
	}

	for i, row := range trees {
		// From left edge
		maxHeight := -1
		for j := range row {
			tree := trees[i][j]
			if tree > maxHeight {
				visible[i][j] = true
				maxHeight = tree
			}
		}
		// From right edge
		maxHeight = -1
		for j := len(row) - 1; j >= 0; j-- {
			tree := trees[i][j]
			if tree > maxHeight {
				visible[i][j] = true
				maxHeight = tree
			}
		}
	}
	for j := range trees[0] {
		// From top edge
		maxHeight := -1
		for i := range trees {
			tree := trees[i][j]
			if tree > maxHeight {
				visible[i][j] = true
				maxHeight = tree
			}
		}
		// From bottom edge
		maxHeight = -1
		for i := len(trees) - 1; i >= 0; i-- {
			tree := trees[i][j]
			if tree > maxHeight {
				visible[i][j] = true
				maxHeight = tree
			}
		}
	}

	visibleCount := 0
	for _, row := range visible {
		for _, v := range row {
			if v {
				visibleCount++
			}
		}
	}
	return visibleCount
}

func main() {
	input := readInput()
	trees := parseInput(input)

	fmt.Println(countVisible(trees))

}
