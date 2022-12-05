package main

import (
	"bufio"
	"fmt"
	"os"
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
	count, from, to int
}

type crates [][]byte

func (cc *crates) execute(ins instruction) {
	c := *cc
	n := len(c[ins.from]) - ins.count
	c[ins.to] = append(c[ins.to], c[ins.from][n:]...)
	c[ins.from] = c[ins.from][:n]
}

func (cc *crates) printTop() {
	c := *cc
	var b strings.Builder
	for _, stack := range c {
		if len(stack) == 0 {
			continue
		}
		n := len(stack) - 1
		b.WriteByte(stack[n])
	}

	fmt.Println(b.String())
}

func parseInput(input []string) (*crates, []instruction) {
	linebreak := 0
	for i, line := range input {
		if len(line) == 0 {
			linebreak = i
			break
		}
	}

	maxStackCount := 9
	stacks := make(crates, maxStackCount)
	for i := range stacks {
		stacks[i] = make([]byte, 0)
	}

	for i := linebreak - 2; i >= 0; i-- {
		line := input[i]
		for j := 0; j < maxStackCount; j++ {
			index := j*4 + 1
			if index >= len(line) {
				break
			}

			if line[index] != ' ' {
				stacks[j] = append(stacks[j], line[index])
			}
		}
	}

	instructions := make([]instruction, 0)
	var count, from, to int
	for _, line := range input[linebreak+1:] {
		n, err := fmt.Sscanf(line, "move %d from %d to %d", &count, &from, &to)
		if err != nil {
			panic(err)
		} else if n != 3 {
			panic("Failed to parse instruction")
		}

		instructions = append(instructions, instruction{count, from - 1, to - 1})
	}

	return &stacks, instructions
}

func main() {
	input := readInput()
	stacks, instructions := parseInput(input)

	for _, ins := range instructions {
		stacks.execute(ins)
	}

	stacks.printTop()
}
