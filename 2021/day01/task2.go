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
	windows := make([]int, len(input)+2)
	for i, s := range input {
		if n, err := strconv.Atoi(s); err == nil {
			windows[i] += n
			windows[i+1] += n
			windows[i+2] += n
		}
	}

	increases := 0
	for i, n := range windows[3:] {
		if n > windows[i+2] {
			increases++
		}
	}

	return increases
}

func main() {
	input := readInput()
	fmt.Println(countIncreases(input))
}
