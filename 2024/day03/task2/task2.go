package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

var keyRegex = regexp.MustCompile(`(mul\(\d+,\d+\)|do\(\)|don't\(\))`)

func solve(input []string) (output string) {
	text := strings.Join(input, "\n")
	sum, enabled := 0, true
	var a, b int
	matches := keyRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if strings.Contains(match[0], "mul") {
			if enabled {
				if n, err := fmt.Sscanf(match[0], "mul(%d,%d)", &a, &b); n != 2 || err != nil {
					panic("invalid input")
				}
				if a < 1000 && b < 1000 {
					sum += a * b
				}
			}
		} else if strings.Contains(match[0], "don't") {
			enabled = false
		} else {
			enabled = true
		}
	}
	return strconv.Itoa(sum)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
