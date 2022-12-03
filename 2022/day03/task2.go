package main

import (
	"bufio"
	"fmt"
	"os"
)

func readInput() []string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

func itemPriority(r rune) int {
	if r >= 'a' {
		return int(r-'a') + 1
	} else {
		return int(r-'A') + 27
	}
}

func findBadge(bags []bag) int {
	for item := range bags[0].c {
		_, ok1 := bags[1].c[item]
		_, ok2 := bags[2].c[item]
		if ok1 && ok2 {
			return itemPriority(item)
		}
	}

	return 0
}

func sumBadges(bags []bag) int {
	sum := 0
	for i := 0; i < len(bags); i += 3 {
		sum += findBadge(bags[i : i+3])
	}

	return sum
}

type bag struct {
	c map[rune]int
}

func parseInput(input []string) []bag {
	bags := make([]bag, 0, len(input))
	for _, line := range input {
		b := bag{make(map[rune]int)}
		for _, r := range line {
			if _, ok := b.c[r]; ok {
				b.c[r]++
			} else {
				b.c[r] = 1
			}
		}

		bags = append(bags, b)
	}

	return bags
}

func main() {
	input := readInput()
	bags := parseInput(input)

	fmt.Println(sumBadges(bags))
}
