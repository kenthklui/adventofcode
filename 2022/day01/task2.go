package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
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

func parseInput(input []string) []int {
	elves := make([]int, 0)

	elf := 0
	for _, line := range input {
		if line == "" {
			if elf != 0 {
				elves = append(elves, elf)
				elf = 0
			}
		} else {
			calorie, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}

			elf += calorie
		}
	}

	if elf != 0 {
		elves = append(elves, elf)
		elf = 0
	}
	return elves
}

func sumMax3(elves []int) int {
	sort.Ints(elves)
	l := len(elves)
	return elves[l-1] + elves[l-2] + elves[l-3]
}

func main() {
	input := readInput()
	elves := parseInput(input)

	fmt.Println(sumMax3(elves))
}
