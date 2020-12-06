package main

import (
	"bufio"
	"fmt"
	"os"
)

type group struct {
	Yes map[rune]bool
}

func countYes() int {
	totalYes := 0

	currentGroup := &group{Yes: make(map[rune]bool)}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			totalYes += len(currentGroup.Yes)
			currentGroup = &group{Yes: make(map[rune]bool)}

			continue
		}

		for _, c := range line {
			currentGroup.Yes[c] = true
		}
	}

	totalYes += len(currentGroup.Yes)

	return totalYes
}

func main() {
	fmt.Println(countYes())
}
