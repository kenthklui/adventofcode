package main

import (
	"bufio"
	"fmt"
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
	var lowestAlign, lowestFuel, lastFuel, lastMarker, marker int

	alignCandidate := crabs[0]
	for marker < len(crabs) && crabs[marker] <= alignCandidate {
		lastMarker++
		marker++
	}

	rightmostCrab := crabs[len(crabs)-1]
	for _, n := range crabs {
		lastFuel += n - alignCandidate
	}

	lowestAlign, lowestFuel = alignCandidate, lastFuel

	for alignCandidate < rightmostCrab {
		alignCandidate++
		for marker < len(crabs) && crabs[marker] <= alignCandidate {
			marker++
		}

		newFuel := lastFuel
		// Everyone who was already passed on the last marker is +1
		newFuel += lastMarker
		// Everyone we passed this time is -1
		difference := marker - lastMarker
		newFuel -= difference
		// Everyone we haven't passed yet is -1
		remaining := len(crabs) - marker
		newFuel -= remaining

		if newFuel < lowestFuel {
			lowestAlign, lowestFuel = alignCandidate, newFuel
		}

		lastFuel = newFuel
		lastMarker = marker
	}

	return lowestAlign, lowestFuel
}

func main() {
	input := readInput()
	crabs := parseInput(input)

	align, fuel := alignCrabs(crabs)

	fmt.Println(align, fuel)
}
