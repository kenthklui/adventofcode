package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

func solve(input []string) (output string) {
	ints := util.ParseInts(input)
	freq := make(map[int]int)
	for _, t := range ints {
		freq[t[1]]++
	}

	distance := 0
	for _, t := range ints {
		distance += t[0] * freq[t[0]]
	}
	return strconv.Itoa(distance)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
