package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func lineValue(line string) int {
	var value, lastDigit int
	for _, c := range line {
		if c >= '0' && c <= '9' {
			digit := int(c - '0')

			if value == 0 {
				value += digit * 10
			}
			lastDigit = digit
		}
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
