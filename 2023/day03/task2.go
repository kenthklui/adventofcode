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

	// Find and handle gears
	gearRatioSum := 0
	for lineNum, line := range input {
		for index, r := range line {
			if r == '*' {
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

				neighborParts := make(map[*partNum]int)
				for l := lineMin; l <= lineMax; l++ {
					for i := indexMin; i <= indexMax; i++ {
						if pn, ok := partsMap[coordinate{l, i}]; ok {
							neighborParts[pn] = 0
						}
					}
				}
				if len(neighborParts) == 2 {
					gearRatio := 1
					for pn := range neighborParts {
						gearRatio *= pn.Value
					}
					gearRatioSum += gearRatio
				}
			}
		}
	}

	return gearRatioSum
}

func main() {
	input := util.StdinReadlines()
	fmt.Println(find(input))
}
