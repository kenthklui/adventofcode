package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

func recurseArrangements(springs string, damaged []int, minLength int) int {
	springs = strings.Trim(springs, ".")

	if len(springs) == 0 {
		if len(damaged) == 0 {
			return 1
		} else {
			return 0
		}
	}
	if len(damaged) == 0 {
		if strings.Index(springs, "#") == -1 {
			return 1
		} else {
			return 0
		}
	}
	if len(springs) < minLength {
		return 0
	}

	count := 0

	// Match current group
	next := damaged[0]
	newMinLength := minLength - next - 1
	before, after, found := strings.Cut(springs, ".")
	for i, r := range before {
		endIndex := i + next
		if endIndex > len(before) { // Too long for current group
			break
		} else if endIndex == len(before) { // Barely fit current group
			count += recurseArrangements(after, damaged[1:], newMinLength)
			break
		}

		if before[endIndex] == '?' {
			count += recurseArrangements(springs[endIndex+1:], damaged[1:], newMinLength)
		}

		if r == '#' {
			break
		}
	}

	// Skip current group if all ???s
	if found && strings.Index(before, "#") == -1 {
		count += recurseArrangements(after, damaged, minLength)
	}

	return count
}

func arrangements(line string) int {
	springs, damagedStr, _ := strings.Cut(line, " ")
	springs = strings.Trim(springs, ".")
	damaged := util.ParseLineInts(damagedStr)
	minLength := len(damaged) - 1
	for _, d := range damaged {
		minLength += d
	}

	return recurseArrangements(springs, damaged, minLength)
}

func main() {
	sum := 0
	for _, line := range util.StdinReadlines() {
		sum += arrangements(line)
	}
	fmt.Println(sum)
}
