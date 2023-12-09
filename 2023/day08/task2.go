package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func gcd(intA, intB int) int {
	if intA < intB {
		intA, intB = intB, intA
	}
	for intB > 0 {
		intA %= intB
		intA, intB = intB, intA
	}
	return intA
}

func lcm(intA, intB int) int {
	return intA / gcd(intA, intB) * intB
}

type location struct {
	Name        string
	Left, Right *location
}

type locationInstruction struct {
	loc              *location
	instructionIndex int
}

type cycle struct {
	head, length, zPos int
}

func traverse(instruction string, network map[string]*location) int {
	sources := make([]*location, 0)
	for name, loc := range network {
		if name[2] == 'A' {
			sources = append(sources, loc)
		}
	}

	cycles := make([]cycle, 0, len(sources))
	for _, source := range sources {
		steps := 0
		curr := source
		travelled := make(map[locationInstruction]int)
		var head, length, zSteps int
		for {
			instructionIndex := steps % len(instruction)
			switch instruction[instructionIndex] {
			case 'L':
				curr = curr.Left
			case 'R':
				curr = curr.Right
			}
			steps++

			if curr.Name[2] == 'Z' {
				// For whatever reason, assume each ghost only visits a single Z node in a cycle?
				zSteps = steps
			}

			key := locationInstruction{curr, instructionIndex}
			if lastSeen, found := travelled[key]; found {
				head = lastSeen
				length = steps - lastSeen
				break
			} else {
				travelled[key] = steps
			}
		}
		cycles = append(cycles, cycle{head, length, zSteps - head})
	}

	// The default input is a special case that can be solved by taking LCM of all the zSteps
	// This solution addresses when Zs don't all magically appear at the cycle length
	increment, steps := cycles[0].length, cycles[0].head+cycles[0].zPos
	for _, c := range cycles[1:] {
		for {
			// This can get stuck in an infinite loop if cycles don't agree
			if ghostOnZ := ((steps-c.head)%c.length == c.zPos); ghostOnZ {
				increment = lcm(increment, c.length)
				break
			} else {
				steps += increment
			}
		}
	}
	return steps
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

func main() {
	input := util.StdinReadlines()
	instruction := input[0]
	network := parseNetwork(input[2:])
	fmt.Println(traverse(instruction, network))
}
