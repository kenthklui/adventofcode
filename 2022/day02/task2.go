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
		response := int(splits[1][0] - 'X')

		var played int
		switch response {
		case 0: // lose
			played = opponent - 1
			if played == 0 {
				played += 3
			}
		case 1: // draw
			played = opponent
			score += 3
		case 2: // win
			played = opponent + 1
			if played > 3 {
				played -= 3
			}
			score += 6
		}
		score += played
	}

	return score
}

func main() {
	input := readInput()
	score := parseInput(input)

	fmt.Println(score)
}
