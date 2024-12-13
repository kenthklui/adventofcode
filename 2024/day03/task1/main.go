package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

var mulRegex = regexp.MustCompile(`mul\(\d+,\d+\)`)

func solve(input []string) (output string) {
	text := strings.Join(input, "\n")
	sum := 0
	var a, b int
	matches := mulRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if n, err := fmt.Sscanf(match[0], "mul(%d,%d)", &a, &b); n != 2 || err != nil {
			panic("invalid input")
		}
		if a < 1000 && b < 1000 {
			sum += a * b
		}
	}
	return strconv.Itoa(sum)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
