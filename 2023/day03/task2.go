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

func gearRatioSum(input []string) int {
	intRegex := regexp.MustCompile(`\d+`)
	partsMap := make(map[coordinate]*partNum)

	// Find and mark numbers
	for lineNum, line := range input {
		matches := intRegex.FindAllStringIndex(line, -1)
		for _, match := range matches {
			if value, err := strconv.Atoi(line[match[0]:match[1]]); err == nil {
				pn := partNum{value, false}
				for index := match[0]; index < match[1]; index++ {
					partsMap[coordinate{lineNum, index}] = &pn
				}
			} else {
				panic(err)
			}
		}
	}

	// Find and handle gears
	sum := 0
	for lineNum, line := range input {
		for index, r := range line {
			if r == '*' {
				neighborParts := make(map[*partNum]int)
				for l := lineNum - 1; l <= lineNum+1; l++ {
					for i := index - 1; i <= index+1; i++ {
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
					sum += gearRatio
				}
			}
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	fmt.Println(gearRatioSum(input))
}
