package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

var numbers = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func lineValue(line string) int {
	var value, lastDigit int

	for i, c := range line {
		digit := 0

		if c >= '0' && c <= '9' {
			digit = int(c - '0')
		} else {
			for j, number := range numbers {
				if strings.HasPrefix(line[i:], number) {
					digit = j + 1
					break
				}
			}
		}

		if digit == 0 {
			continue
		}

		if value == 0 {
			value += digit * 10
		}
		lastDigit = digit
	}
	value += lastDigit

	return value
}

func calibrationValue(input []string) int {
	sum := 0
	for _, line := range input {
		sum += lineValue(line)
	}

	return sum
}

func main() {
	input := util.StdinReadlines()
	fmt.Println(calibrationValue(input))
}
