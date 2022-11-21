package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	// 1, 4, 7, 8 are kinda easy as the only 2, 4, 3 and 7 segment chars
	count1478 := 0

	for _, s := range input {
		splits := strings.Split(s, " | ")
		_, output := splits[0], splits[1]

		for _, p := range strings.Split(output, " ") {
			switch len(p) {
			case 2, 3, 4, 7:
				count1478++
			}
		}
	}

	return count1478
}

func main() {
	input := readInput()

	fmt.Println(parseInput(input))
}
