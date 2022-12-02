package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func parseInput(input []string) int {
	score := 0
	for _, line := range input {
		splits := strings.Split(line, " ")

		opponent := int(splits[0][0]-'A') + 1
		response := int(splits[1][0]-'X') + 1

		diff := response - opponent
		if diff < 0 {
			diff += 3
		}

		score += response
		switch diff {
		case 0: // draw
			score += 3
		case 1: // win
			score += 6
		case 2: // loss
		}
	}

	return score
}

func main() {
	input := readInput()
	score := parseInput(input)

	fmt.Println(score)
}
