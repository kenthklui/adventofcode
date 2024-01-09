package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
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

func run(instructions []*instruction) []int {
	registerValues := make([]int, 0, 240)

	register := 1
	for _, ins := range instructions {
		for i := 0; i < ins.cycles; i++ {
			registerValues = append(registerValues, register)
		}

		register += ins.val
	}

	return registerValues
}

const (
	rows    = 6
	columns = 40
)

func printScreen(registerValues []int) {
	screen := make([][]bool, rows)
	for row := range screen {
		screen[row] = make([]bool, columns)
	}

	for pixel, register := range registerValues {
		row, column := pixel/columns, pixel%columns
		diff := column - register
		screen[row][column] = (diff*diff <= 1)
	}
	util.PrintScreen(screen)
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
	registerValues := run(instructions)
	printScreen(registerValues)
}
