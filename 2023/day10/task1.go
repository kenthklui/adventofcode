package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type loc struct {
	x, y int
}

func (l loc) up() loc    { return loc{l.x, l.y - 1} }
func (l loc) down() loc  { return loc{l.x, l.y + 1} }
func (l loc) left() loc  { return loc{l.x - 1, l.y} }
func (l loc) right() loc { return loc{l.x + 1, l.y} }

type tube struct {
	r        rune
	loc      loc
	neighbor [4]*tube // up, down, left, right
}

func (t *tube) exit(entrance *tube) *tube {
	for _, out := range t.neighbor {
		if out != nil && out != entrance {
			return out
		}
	}
	return nil
}

func parseLoop(input []string) *tube {
	var start *tube

	loop := make(map[loc]*tube)
	for lineNum, line := range input {
		for charNum, r := range line {
			l := loc{charNum, lineNum}
			loop[l] = &tube{r: r, loc: l}
		}
	}

	for l, t := range loop {
		switch t.r {
		case '|':
			t.neighbor[0] = loop[l.up()]
			t.neighbor[1] = loop[l.down()]
		case '-':
			t.neighbor[2] = loop[l.left()]
			t.neighbor[3] = loop[l.right()]
		case 'L':
			t.neighbor[0] = loop[l.up()]
			t.neighbor[3] = loop[l.right()]
		case 'J':
			t.neighbor[0] = loop[l.up()]
			t.neighbor[2] = loop[l.left()]
		case '7':
			t.neighbor[1] = loop[l.down()]
			t.neighbor[2] = loop[l.left()]
		case 'F':
			t.neighbor[1] = loop[l.down()]
			t.neighbor[3] = loop[l.right()]
		case 'S':
			start = t
		case '.':
			delete(loop, l)
		default:
			panic("Invalid char")
		}
	}

	up, ok := loop[start.loc.up()]
	if ok && up.neighbor[1] == start {
		start.neighbor[0] = up
	}
	down, ok := loop[start.loc.down()]
	if ok && down.neighbor[0] == start {
		start.neighbor[1] = down
	}
	left, ok := loop[start.loc.left()]
	if ok && left.neighbor[3] == start {
		start.neighbor[2] = left
	}
	right, ok := loop[start.loc.right()]
	if ok && right.neighbor[2] == start {
		start.neighbor[3] = right
	}

	return start
}

func (t *tube) farthestSteps() int {
	length := 1
	for prev, curr := t, t.exit(nil); curr != t; prev, curr = curr, curr.exit(prev) {
		length++
	}
	return length / 2
}

func main() {
	input := util.StdinReadlines()
	start := parseLoop(input)
	fmt.Println(start.farthestSteps())
}
