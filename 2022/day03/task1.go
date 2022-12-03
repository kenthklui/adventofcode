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

func sumPriorities(bags []bag) int {
	sum := 0
	for _, b := range bags {
		for item := range b.c1 {
			if _, ok := b.c2[item]; ok {
				sum += itemPriority(item)
				break
			}
		}
	}

	return sum
}

type bag struct {
	c1, c2 map[rune]int
}

func parseInput(input []string) []bag {
	bags := make([]bag, 0, len(input))
	for _, line := range input {
		b := bag{make(map[rune]int), make(map[rune]int)}
		bagsize := len(line) / 2
		for _, r := range line[:bagsize] {
			if _, ok := b.c1[r]; ok {
				b.c1[r]++
			} else {
				b.c1[r] = 1
			}
		}
		for _, r := range line[bagsize:] {
			if _, ok := b.c2[r]; ok {
				b.c2[r]++
			} else {
				b.c2[r] = 1
			}
		}

		bags = append(bags, b)
	}

	return bags
}

func main() {
	input := readInput()
	bags := parseInput(input)

	fmt.Println(sumPriorities(bags))
}
