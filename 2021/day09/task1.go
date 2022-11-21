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
	heightMap := make([][]int, len(input))

	for i, s := range input {
		heightMap[i] = make([]int, len(s))
		for j, c := range s {
			heightMap[i][j] = int(c - '0')
		}
	}

	return heightMap
}

func lowPoints(heightMap [][]int) [][]bool {
	isLow := make([][]bool, len(heightMap))

	for i, row := range heightMap {
		isLow[i] = make([]bool, len(row))
		for j := range isLow[i] {
			isLow[i][j] = true
		}
	}

	// compare right
	for i, row := range heightMap {
		for j, val := range row[:len(row)-1] {
			if val > row[j+1] {
				isLow[i][j] = false
			} else if val < row[j+1] {
				isLow[i][j+1] = false
			} else {
				isLow[i][j] = false
				isLow[i][j+1] = false
			}
		}
	}

	// compare down
	for i, row := range heightMap[:len(heightMap)-1] {
		for j, val := range row {
			if val > heightMap[i+1][j] {
				isLow[i][j] = false
			} else if val < heightMap[i+1][j] {
				isLow[i+1][j] = false
			} else {
				isLow[i][j] = false
				isLow[i+1][j] = false
			}
		}
	}

	return isLow
}

func sumRisk(heightMap [][]int, lows [][]bool) int {
	sum := 0

	for i, row := range heightMap {
		for j, val := range row {
			if lows[i][j] {
				sum += val + 1
			}
		}
	}

	return sum
}

func main() {
	input := readInput()
	heightMap := parseInput(input)
	lows := lowPoints(heightMap)

	fmt.Println(sumRisk(heightMap, lows))
}
