package main

import (
	"bufio"
	"fmt"
	"os"
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

type instruction struct {
	cycles int
	val    int
}

func run(instructions []*instruction) int {
	signalSum := 0

	checkpoint := 20
	cycle := 1
	register := 1
	for _, ins := range instructions {
		if cycle+ins.cycles > checkpoint {
			signalSum += checkpoint * register
			checkpoint += 40
			// Break needed when 220 is reached?
		}

		cycle += ins.cycles
		register += ins.val
	}

	return signalSum
}

func parseInput(input []string) []*instruction {
	instructions := make([]*instruction, 0, len(input))
	for _, line := range input {
		tokens := strings.Split(line, " ")

		var ins instruction
		switch tokens[0] {
		case "noop":
			ins.cycles = 1
			ins.val = 0 // for safety
		case "addx":
			ins.cycles = 2
			ins.val, _ = strconv.Atoi(tokens[1])
		default:
			panic("Invalid operation")
		}
		instructions = append(instructions, &ins)
	}

	return instructions
}

func main() {
	input := readInput()
	instructions := parseInput(input)
	fmt.Println(run(instructions))
}
