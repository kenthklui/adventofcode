package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type node struct {
	val  int8
	next []*node
}

func (n *node) countTrails() int {
	if n.val == int8(9) {
		return 1
	}

	count := 0
	for _, next := range n.next {
		count += next.countTrails()
	}
	return count
}

type nodeMap struct {
	heads []*node
}

func parse(input []string) nodeMap {
	heads := []*node{}
	nodes := make([][]*node, len(input))
	for y, line := range input {
		nodes[y] = make([]*node, len(line))
		for x, char := range line {
			val := int8(char - '0')
			curr := &node{val, []*node{}}

			if val == 0 {
				heads = append(heads, curr)
			}

			if x > 0 {
				prevX := nodes[y][x-1]
				switch val - prevX.val {
				case 1:
					prevX.next = append(prevX.next, curr)
				case -1:
					curr.next = append(curr.next, prevX)
				}
			}
			if y > 0 {
				prevY := nodes[y-1][x]
				switch val - prevY.val {
				case 1:
					prevY.next = append(prevY.next, curr)
				case -1:
					curr.next = append(curr.next, prevY)
				}
			}

			nodes[y][x] = curr
		}
	}
	return nodeMap{heads}
}

func (nm nodeMap) trailheads() int {
	count := 0
	for _, head := range nm.heads {
		count += head.countTrails()
	}
	return count
}

func solve(input []string) (output string) {
	nm := parse(input)
	return strconv.Itoa(nm.trailheads())
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
