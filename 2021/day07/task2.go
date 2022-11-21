package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
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

func parseInput(input []string) []int {
	crabs := make([]int, 0)

	for _, s := range strings.Split(input[0], ",") {
		if n, err := strconv.Atoi(s); err == nil {
			crabs = append(crabs, n)
		}
	}

	// sorting the crabs will help later
	sort.Ints(crabs)

	return crabs
}

func alignCrabs(crabs []int) (int, int) {
	lowestAlign, lowestFuel := math.MaxInt32, math.MaxInt32

	rightmostCrab := crabs[len(crabs)-1]
	for alignCandidate := crabs[0]; alignCandidate <= rightmostCrab; alignCandidate++ {
		fuel := 0

		for _, c := range crabs {
			diff := c - alignCandidate
			if diff < 0 {
				diff = -diff
			}

			fuel += diff * (diff + 1) / 2
		}

		if fuel < lowestFuel {
			lowestAlign, lowestFuel = alignCandidate, fuel
		}
	}

	return lowestAlign, lowestFuel
}

func main() {
	input := readInput()
	crabs := parseInput(input)

	align, fuel := alignCrabs(crabs)

	fmt.Println(align, fuel)
}
