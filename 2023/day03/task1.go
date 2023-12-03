package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type coordinate struct {
	LineNum, Index int
}

type partNum struct {
	Value int
	Valid bool
}

func partNumSum(input []string) int {
	intRegex := regexp.MustCompile(`\d+`)
	partsList := make([]*partNum, 0)
	partsMap := make(map[coordinate]*partNum)

	// Find and mark numbers
	for lineNum, line := range input {
		matches := intRegex.FindAllStringIndex(line, -1)
		for _, match := range matches {
			if value, err := strconv.Atoi(line[match[0]:match[1]]); err == nil {
				pn := partNum{value, false}
				partsList = append(partsList, &pn)
				for index := match[0]; index < match[1]; index++ {
					partsMap[coordinate{lineNum, index}] = &pn
				}
			} else {
				panic(err)
			}
		}
	}

	// Find and use symbols
	for lineNum, line := range input {
		for index, r := range line {
			if r != '.' && (r < '0' || r > '9') {
				for l := lineNum - 1; l <= lineNum+1; l++ {
					for i := index - 1; i <= index+1; i++ {
						if pn, ok := partsMap[coordinate{l, i}]; ok {
							pn.Valid = true
						}
					}
				}
			}
		}
	}

	sum := 0
	for _, pn := range partsList {
		if pn.Valid {
			sum += pn.Value
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	fmt.Println(partNumSum(input))
}
