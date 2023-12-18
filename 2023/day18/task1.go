package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type vec struct {
	x, y int
}

func (v vec) move(x, y int) vec { return vec{v.x + x, v.y + y} }

type color string // handle this later

type lagoon struct {
	trenches               map[vec]color
	minX, maxX, minY, maxY int
}

func (l *lagoon) setColor(v vec, c color) {
	if v.x < l.minX {
		l.minX = v.x
	}
	if v.y < l.minY {
		l.minY = v.y
	}
	if v.x > l.maxX {
		l.maxX = v.x
	}
	if v.y > l.maxY {
		l.maxY = v.y
	}
	l.trenches[v] = c
}

func (l lagoon) volume() int {
	left, right, top, bottom := l.minX-1, l.maxX+1, l.minY-1, l.maxY+1
	width, height := right-left+1, bottom-top+1
	whole := width * height

	terrain := make([]bool, whole)
	start := vec{left, top}
	terrain[0] = true

	queue := make([]vec, 0, whole)
	queue = append(queue, start)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.x > left {
			n := curr.move(-1, 0)
			_, isWall := l.trenches[n]
			val := (n.x - left) + (n.y-top)*width
			checked := terrain[val]
			if !isWall && !checked {
				terrain[val] = true
				queue = append(queue, n)
			}
		}
		if curr.x < right {
			n := curr.move(1, 0)
			_, isWall := l.trenches[n]
			val := (n.x - left) + (n.y-top)*width
			checked := terrain[val]
			if !isWall && !checked {
				terrain[val] = true
				queue = append(queue, n)
			}
		}
		if curr.y > top {
			n := curr.move(0, -1)
			_, isWall := l.trenches[n]
			val := (n.x - left) + (n.y-top)*width
			checked := terrain[val]
			if !isWall && !checked {
				terrain[val] = true
				queue = append(queue, n)
			}
		}
		if curr.y < bottom {
			n := curr.move(0, 1)
			_, isWall := l.trenches[n]
			val := (n.x - left) + (n.y-top)*width
			checked := terrain[val]
			if !isWall && !checked {
				terrain[val] = true
				queue = append(queue, n)
			}
		}
	}
	for _, b := range terrain {
		if b {
			whole--
		}
	}
	return whole
}

func parseLagoon(input []string) *lagoon {
	var dir, colorStr string
	var meters int
	curr := vec{0, 0}

	l := &lagoon{trenches: make(map[vec]color)}
	l.setColor(curr, color(""))

	for _, line := range input {
		if n, err := fmt.Sscanf(line, "%s %d %s", &dir, &meters, &colorStr); err != nil {
			panic(err)
		} else if n != 3 {
			panic("Failed to parse 3 params")
		}
		col := color(colorStr[1:7])
		switch dir {
		case "U":
			for i := 1; i <= meters; i++ {
				curr = curr.move(0, -1)
				l.setColor(curr, col)
			}
		case "D":
			for i := 1; i <= meters; i++ {
				curr = curr.move(0, 1)
				l.setColor(curr, col)
			}
		case "L":
			for i := 1; i <= meters; i++ {
				curr = curr.move(-1, 0)
				l.setColor(curr, col)
			}
		case "R":
			for i := 1; i <= meters; i++ {
				curr = curr.move(1, 0)
				l.setColor(curr, col)
			}
		default:
			panic("Invalid direction")
		}
	}

	return l
}

func main() {
	input := util.StdinReadlines()
	l := parseLagoon(input)
	fmt.Println(l.volume())
}
