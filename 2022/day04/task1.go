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

func parseInput(input []string) int {
	var aMin, aMax, bMin, bMax int
	encapsule := 0
	for _, line := range input {
		_, err := fmt.Sscanf(line, "%d-%d,%d-%d", &aMin, &aMax, &bMin, &bMax)
		if err != nil {
			panic(err)
		}

		if aMin <= bMin && aMax >= bMax {
			encapsule++
		} else if aMin >= bMin && aMax <= bMax {
			encapsule++
		}
	}
	return encapsule
}

func main() {
	input := readInput()
	encapsule := parseInput(input)
	fmt.Println(encapsule)
}
