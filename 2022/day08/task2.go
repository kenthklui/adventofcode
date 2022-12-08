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

func maxScenic(trees [][]int) int {
	var maxScenic int
	for i, row := range trees {
		for j, tree := range row {
			var upCount, downCount, leftCount, rightCount int
			for up := i - 1; up >= 0; up-- {
				upCount++
				if trees[up][j] >= tree {
					break
				}
			}
			for down := i + 1; down < len(trees); down++ {
				downCount++
				if trees[down][j] >= tree {
					break
				}
			}
			for left := j - 1; left >= 0; left-- {
				leftCount++
				if trees[i][left] >= tree {
					break
				}
			}
			for right := j + 1; right < len(row); right++ {
				rightCount++
				if trees[i][right] >= tree {
					break
				}
			}

			scenic := upCount * downCount * leftCount * rightCount
			if scenic > maxScenic {
				maxScenic = scenic
			}
		}
	}

	return maxScenic
}

func main() {
	input := readInput()
	trees := parseInput(input)

	fmt.Println(maxScenic(trees))

}
