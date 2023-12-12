package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type key struct {
	charsRemain, groupsRemain int
}

type memoizer struct {
	data map[key]int
}

func newMemoizer() *memoizer                  { return &memoizer{make(map[key]int)} }
func (m *memoizer) key(s string, i []int) key { return key{len(s), len(i)} }

func recurseArrangements(springs string, damaged []int, minLength int, m *memoizer) int {
	springs = strings.Trim(springs, ".")

	if v, cached := m.data[m.key(springs, damaged)]; cached {
		return v
	}

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
			count += recurseArrangements(after, damaged[1:], newMinLength, m)
			break
		}

		if before[endIndex] == '?' {
			count += recurseArrangements(springs[endIndex+1:], damaged[1:], newMinLength, m)
		}

		if r == '#' {
			break
		}
	}

	// Skip current group if all ???s
	if found && strings.Index(before, "#") == -1 {
		count += recurseArrangements(after, damaged, minLength, m)
	}

	m.data[m.key(springs, damaged)] = count

	return count
}

func arrangements(line string) int {
	springs, damagedStr, _ := strings.Cut(line, " ")
	damaged := util.ParseLineInts(damagedStr)

	springsExt := strings.Join([]string{springs, springs, springs, springs, springs}, "?")
	springsExt = strings.Trim(springsExt, ".")
	damagedExt := make([]int, 0, len(damaged)*5)
	for i := 0; i < 5; i++ {
		damagedExt = append(damagedExt, damaged...)
	}

	minLength := len(damagedExt) - 1
	for _, c := range damagedExt {
		minLength += c
	}

	return recurseArrangements(springsExt, damagedExt, minLength, newMemoizer())
}

func main() {
	sum := 0
	for _, line := range util.StdinReadlines() {
		sum += arrangements(line)
	}
	fmt.Println(sum)
}
