package main

import (
	"bufio"
	"fmt"
	"os"
)

type group struct {
	Yes     map[rune]int
	Members int
}

func newGroup() *group {
	yes := make(map[rune]int)
	for i := 'a'; i <= 'z'; i++ {
		yes[i] = 0
	}

	return &group{Yes: yes, Members: 0}
}

func countYes() int {
	totalYes := 0

	currentGroup := newGroup()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			for _, yesCount := range currentGroup.Yes {
				if yesCount == currentGroup.Members {
					totalYes++
				}
			}

			currentGroup = newGroup()

			continue
		}

		currentGroup.Members++
		for _, c := range line {
			currentGroup.Yes[c]++
		}
	}

	for _, yesCount := range currentGroup.Yes {
		if yesCount == currentGroup.Members {
			totalYes++
		}
	}

	return totalYes
}

func main() {
	fmt.Println(countYes())
}
