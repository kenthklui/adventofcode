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

func uniq(s string) bool {
	for i, r1 := range s {
		for _, r2 := range s[i+1:] {
			if r1 == r2 {
				return false
			}
		}

	}
	return true
}

func firstStartOfPacket(input string) int {
	windowSize := 4
	for i := range input[windowSize:] {
		if uniq(input[i : i+windowSize]) {
			return i + windowSize
		}
	}

	return -1
}

func main() {
	input := readInput()
	for _, line := range input {
		fmt.Println(firstStartOfPacket(line))
	}
}
