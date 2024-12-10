package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type node struct {
	val     int8
	next    []*node
	reached bool
}

func (n *node) followTrails() {
	if n.val == 9 {
		n.reached = true
		return
	}

	for _, next := range n.next {
		next.followTrails()
	}
}

type nodeMap struct {
	heads, ends []*node
}

func parse(input []string) nodeMap {
	heads, ends := []*node{}, []*node{}
	nodes := make([][]*node, len(input))
	for y, line := range input {
		nodes[y] = make([]*node, len(line))
		for x, char := range line {
			val := int8(char - '0')
			curr := &node{val, []*node{}, false}

			if val == 0 {
				heads = append(heads, curr)
			} else if val == 9 {
				ends = append(ends, curr)
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
	return nodeMap{heads, ends}
}

func (nm nodeMap) trailheads() int {
	count := 0
	for _, head := range nm.heads {
		head.followTrails()
		for _, end := range nm.ends {
			if end.reached {
				count++
				end.reached = false
			}
		}
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
