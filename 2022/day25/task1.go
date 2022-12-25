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

func snafuToInt(snafu string) int {
	sum := 0
	for _, r := range snafu {
		sum *= 5
		switch r {
		case '-':
			sum += -1
		case '=':
			sum += -2
		default:
			sum += int(r - '0')
		}
	}
	return sum
}

func intToSnafu(i int) string {
	digits := make([]rune, 0)

	var carry bool
	for i > 0 {
		val := i % 5
		if carry {
			val++
			carry = false
		}
		if val >= 3 {
			val -= 5
			carry = true
		}
		switch val {
		case -1:
			digits = append(digits, rune('-'))
		case -2:
			digits = append(digits, rune('='))
		default:
			digits = append(digits, rune('0'+val))
		}

		i /= 5
	}

	var b strings.Builder
	for n := len(digits) - 1; n >= 0; n-- {
		b.WriteRune(digits[n])
	}

	return b.String()
}

func parseInput(input []string) string {
	sum := 0
	for _, line := range input {
		sum += snafuToInt(line)
	}
	return intToSnafu(sum)
}

func main() {
	input := readInput()
	snafu := parseInput(input)
	fmt.Println(snafu)
}
