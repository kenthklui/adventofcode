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

func illegalScore(r rune) int {
	switch r {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	default:
		panic("Unknown illegal char")
	}
}

func syntaxErrorScoreLine(s string) int {
	stack := make([]rune, 0)

	for _, r := range s {
		switch r {
		case '(', '[', '{', '<':
			stack = append(stack, r)
		case ')':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '(' {
				return illegalScore(r)
			}
		case ']':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '[' {
				return illegalScore(r)
			}
		case '}':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '{' {
				return illegalScore(r)
			}
		case '>':
			pop := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if pop != '<' {
				return illegalScore(r)
			}
		}
	}

	return 0
}

func main() {
	input := readInput()

	score := 0
	for _, s := range input {
		score += syntaxErrorScoreLine(s)
	}

	fmt.Println(score)
}
