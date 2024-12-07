package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type equation struct {
	result int
	inputs []int
}

func newEq(line string) equation {
	ints := util.ParseLineInts(line)
	return equation{ints[0], ints[1:]}
}

func (eq equation) solveFor(value int, inputs []int) bool {
	if value > eq.result {
		return false
	}

	if len(inputs) == 0 {
		return value == eq.result
	}

	return eq.solveFor(value*inputs[0], inputs[1:]) || eq.solveFor(value+inputs[0], inputs[1:])
}

func (eq equation) isSolvable() bool {
	return eq.solveFor(eq.inputs[0], eq.inputs[1:])
}

func solve(input []string) (output string) {
	sum := 0
	for _, line := range input {
		eq := newEq(line)
		if eq.isSolvable() {
			sum += eq.result
		}
	}
	return strconv.Itoa(sum)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
