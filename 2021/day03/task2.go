package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	Pop    int
	Ones   *node
	Zeroes *node
}

func (n *node) Insert(s string, i int) {
	if i == len(s) {
		return
	}

	n.Pop++

	switch s[i] {
	case '0':
		if n.Zeroes == nil {
			n.Zeroes = &node{1, nil, nil}
		}
		n.Zeroes.Insert(s, i+1)
	case '1':
		if n.Ones == nil {
			n.Ones = &node{1, nil, nil}
		}
		n.Ones.Insert(s, i+1)
	default:
		panic("Non binary digit")
	}
}

func (n *node) Oxygen() string {
	if n.Pop == 1 {
		return ""
	}

	ones := 0
	if n.Ones != nil {
		ones = n.Ones.Pop
	}

	zeroes := 0
	if n.Zeroes != nil {
		zeroes = n.Zeroes.Pop
	}

	if ones >= zeroes {
		return "1" + n.Ones.Oxygen()
	} else {
		return "0" + n.Zeroes.Oxygen()
	}
}

func (n *node) CO2() string {
	if n.Pop == 1 {
		return ""
	}

	ones := 0
	if n.Ones != nil {
		ones = n.Ones.Pop
	}

	zeroes := 0
	if n.Zeroes != nil {
		zeroes = n.Zeroes.Pop
	}

	if zeroes < ones {
		if zeroes > 0 {
			return "0" + n.Zeroes.CO2()
		} else {
			return "1" + n.Ones.CO2()
		}
	} else if zeroes == ones {
		// pop must > 1
		return "0" + n.Zeroes.CO2()
	} else { // zeroes > ones
		if ones > 0 {
			return "1" + n.Ones.CO2()
		} else {
			return "0" + n.Zeroes.CO2()
		}
	}
}

func intFromBinary(binary string) int {
	var n int
	for _, c := range binary {
		n *= 2

		if c == '1' {
			n++
		}
	}

	return n
}

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

func buildTree(input []string) *node {
	tree := node{}

	for _, s := range input {
		tree.Insert(s, 0)
	}

	return &tree
}

func main() {
	input := readInput()
	tree := buildTree(input)

	oxygen, co2 := tree.Oxygen(), tree.CO2()

	fmt.Println(intFromBinary(oxygen) * intFromBinary(co2))
}
