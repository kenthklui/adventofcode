package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

func parse(line string) []int {
	totalSize := 0
	for _, c := range line {
		size := int(c - '0')
		totalSize += size
	}

	register := make([]int, totalSize)
	id, ptr := 0, 0
	for i, c := range line {
		size := int(c - '0')
		if i%2 == 0 {
			for j := 0; j < size; j++ {
				register[ptr+j] = id
			}
			id++
		} else {
			for j := 0; j < size; j++ {
				register[ptr+j] = -1
			}
		}
		ptr += size
	}

	return register
}

func checksum(register []int) int {
	checksum := 0
	for i, id := range register {
		checksum += i * id
	}
	return checksum
}

func solve(input []string) (output string) {
	register := parse(input[0])
	ptr := len(register) - 1
	for i, id := range register {
		if id == -1 {
			for ptr > i && register[ptr] == -1 {
				ptr--
			}
			if ptr <= i {
				break
			}
			register[i], register[ptr] = register[ptr], register[i]
		}
	}
	register = register[:ptr]
	return strconv.Itoa(checksum(register))
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
