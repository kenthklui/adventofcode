package main

import (
	"bufio"
	"fmt"
	"os"
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

func parseInput(input []string) int {
	maxCalorie := 0

	elf := 0
	for _, line := range input {
		if line == "" {
			if elf > maxCalorie {
				maxCalorie = elf
			}
			elf = 0
			continue
		} else {
			calorie, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}

			elf += calorie
		}
	}
	return maxCalorie
}

func main() {
	input := readInput()
	maxCalorie := parseInput(input)

	fmt.Println(maxCalorie)
}
