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

func find(input []string) int {
	intRegex := regexp.MustCompile(`\d+`)
	partsMap := make(map[coordinate]*partNum)
	partsList := make([]*partNum, 0)

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
				lineMin := lineNum - 1
				if lineMin < 0 {
					lineMin++
				}
				lineMax := lineNum + 1
				if lineMax == len(input) {
					lineMax--
				}

				indexMin := index - 1
				if indexMin < 0 {
					indexMin++
				}
				indexMax := index + 1
				if indexMax == len(input) {
					indexMax--
				}

				for l := lineMin; l <= lineMax; l++ {
					for i := indexMin; i <= indexMax; i++ {
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
	fmt.Println(find(input))
}
