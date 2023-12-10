package util

import (
	"regexp"
	"strconv"
)

var intRegex = regexp.MustCompile(`-?\d+`)

func ParseLineInts(line string) []int {
	matches := intRegex.FindAllString(line, -1)
	ints := make([]int, 0, len(matches))
	for _, match := range matches {
		if n, err := strconv.Atoi(match); err == nil {
			ints = append(ints, n)
		} else {
			panic(err)
		}
	}
	return ints
}

func ParseInts(input []string) [][]int {
	ints := make([][]int, len(input))
	for i, line := range input {
		ints[i] = ParseLineInts(line)
	}
	return ints
}
