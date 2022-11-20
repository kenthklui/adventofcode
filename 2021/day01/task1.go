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

func countIncreases(input []string) int {
	increases := 0
	last := 100000000 // set arbitrary large to avoid counting first increase

	for _, s := range input {
		if current, err := strconv.Atoi(s); err == nil {
			if current > last {
				increases++
			}
			last = current
		}
	}

	return increases
}

func main() {
	input := readInput()
	fmt.Println(countIncreases(input))
}
