package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type location struct {
	Name        string
	Left, Right *location
}

func parseNetwork(input []string) map[string]*location {
	network := make(map[string]*location)

	var source, left, right *location
	var exists bool
	for _, line := range input {
		sourceName, leftName, rightName := line[0:3], line[7:10], line[12:15]
		if source, exists = network[sourceName]; !exists {
			source = &location{sourceName, nil, nil}
			network[sourceName] = source
		}
		if left, exists = network[leftName]; !exists {
			left = &location{leftName, nil, nil}
			network[leftName] = left
		}
		if right, exists = network[rightName]; !exists {
			right = &location{rightName, nil, nil}
			network[rightName] = right
		}
		source.Left, source.Right = left, right
	}

	return network
}

func traverse(instruction string, network map[string]*location) int {
	curr, dest := network["AAA"], network["ZZZ"]

	steps := 0
	for curr != dest {
		move := instruction[steps%len(instruction)]
		switch move {
		case 'L':
			curr = curr.Left
		case 'R':
			curr = curr.Right
		}
		steps++
	}

	return steps
}

func main() {
	input := util.StdinReadlines()

	instruction := input[0]
	network := parseNetwork(input[2:])

	fmt.Println(traverse(instruction, network))
}
