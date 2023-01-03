package util

import (
	"regexp"
	"strconv"
)

func ParseInts(input []string) [][]int {
	ints := make([][]int, len(input))
	intRegex := regexp.MustCompile(`-?\d+`)
	for i, line := range input {
		matches := intRegex.FindAllString(line, -1)
		ints[i] = make([]int, 0, len(matches))
		for _, match := range matches {
			if n, err := strconv.Atoi(match); err == nil {
				ints[i] = append(ints[i], n)
			} else {
				panic(err)
			}
		}
	}
	return ints
}
