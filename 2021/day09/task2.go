package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	heightMap := make([][]int, len(input))

	for i, s := range input {
		heightMap[i] = make([]int, len(s))
		for j, c := range s {
			heightMap[i][j] = int(c - '0')
		}
	}

	return heightMap
}

func basinSizes(heightMap [][]int) []int {
	counted := make([][]bool, len(heightMap))
	for i, row := range heightMap {
		counted[i] = make([]bool, len(row))
	}

	basins := make([]int, 0)
	for i, row := range heightMap {
		for j := range row {
			size := markBasin(heightMap, counted, i, j)
			if size > 0 {
				basins = append(basins, size)
			}
		}
	}

	return basins
}

func markBasin(heightMap [][]int, counted [][]bool, i, j int) int {
	// OOB check
	if i < 0 || i >= len(heightMap) {
		return 0
	}
	if j < 0 || j >= len(heightMap[0]) {
		return 0
	}

	// Counted check
	if counted[i][j] {
		return 0
	}

	counted[i][j] = true

	// 9 check
	if heightMap[i][j] == 9 {
		return 0
	}

	size := 1

	size += markBasin(heightMap, counted, i-1, j) // up
	size += markBasin(heightMap, counted, i+1, j) // down
	size += markBasin(heightMap, counted, i, j-1) // left
	size += markBasin(heightMap, counted, i, j+1) // right

	return size
}

func multiplyTopBasins(basins []int) int {
	basinsCopy := make([]int, len(basins))
	copy(basinsCopy, basins)

	sort.Ints(basinsCopy)

	product := 1
	for i := 1; i <= 3; i++ {
		product *= basinsCopy[len(basinsCopy)-i]
	}

	return product
}

func main() {
	input := readInput()
	heightMap := parseInput(input)
	basins := basinSizes(heightMap)

	fmt.Println(multiplyTopBasins(basins))
}
