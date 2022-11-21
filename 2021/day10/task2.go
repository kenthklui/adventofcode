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

func charScore(r rune) int {
	switch r {
	case '(':
		return 1
	case '[':
		return 2
	case '{':
		return 3
	case '<':
		return 4
	default:
		panic("Unknown legal char")
	}
}

func scoreLine(s string) int {
	stack := make([]rune, 0)

	for _, r := range s {
		switch r {
		case '(', '[', '{', '<':
			stack = append(stack, r)
		case ')':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '(' {
				return 0
			}
		case ']':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '[' {
				return 0
			}
		case '}':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '{' {
				return 0
			}
		case '>':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '<' {
				return 0
			}
		}
	}

	score := 0
	for len(stack) > 0 {
		score *= 5

		pop := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		score += charScore(pop)
	}

	return score
}

func main() {
	input := readInput()

	scores := make([]int, 0)
	for _, s := range input {
		score := scoreLine(s)
		if score != 0 {
			scores = append(scores, score)
		}
	}
	sort.Ints(scores)
	median := scores[len(scores)/2]

	fmt.Println(median)
}
