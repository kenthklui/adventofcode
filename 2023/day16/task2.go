package main

import (
	"fmt"
	"slices"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

type vec struct {
	x, y int
}

var up = vec{0, -1}
var down = vec{0, 1}
var left = vec{-1, 0}
var right = vec{1, 0}

func (v vec) add(dir vec) vec { return vec{v.x + dir.x, v.y + dir.y} }
func (v vec) inBounds(width, height int) bool {
	return v.x >= 0 && v.y >= 0 && v.x < width && v.y < height
}

type cave struct {
	width, height int
	layout        []string

	energized [][]bool
	splitted  map[vec]void
}

func makeCave(input []string) *cave {
	width, height := len(input[0]), len(input)
	e := make([][]bool, height)
	for i := range e {
		e[i] = make([]bool, width)
	}

	return &cave{
		width:     width,
		height:    height,
		layout:    input,
		energized: e,
		splitted:  make(map[vec]void),
	}
}

func (c *cave) splitterUsed(loc vec) bool {
	_, used := c.splitted[loc]
	if !used {
		c.splitted[loc] = nul
	}
	return used
}

func (c *cave) sendBeam(origin, dir vec) {
	loc := origin
	for blocked := false; !blocked && loc.inBounds(c.width, c.height); loc = loc.add(dir) {
		c.energized[loc.y][loc.x] = true
		switch c.layout[loc.y][loc.x] {
		case '.':
			continue
		case '/':
			switch dir {
			case up:
				dir = right
			case down:
				dir = left
			case left:
				dir = down
			case right:
				dir = up
			}
		case '\\':
			switch dir {
			case up:
				dir = left
			case down:
				dir = right
			case left:
				dir = up
			case right:
				dir = down
			}
		case '|':
			if dir == left || dir == right {
				if !c.splitterUsed(loc) {
					c.sendBeam(loc, up)
					c.sendBeam(loc, down)
				}
				blocked = true
			}
		case '-':
			if dir == up || dir == down {
				if !c.splitterUsed(loc) {
					c.sendBeam(loc, left)
					c.sendBeam(loc, right)
				}
				blocked = true
			}
		}
	}
}

func (c *cave) countEnergized() int {
	sum := 0
	for _, row := range c.energized {
		for _, b := range row {
			if b {
				sum++
			}
		}
	}
	return sum
}

func (c *cave) reset() {
	for y, row := range c.energized {
		for x := range row {
			c.energized[y][x] = false
		}
	}
	c.splitted = make(map[vec]void)
}

func (c *cave) maxEnergized() int {
	energized := make([]int, 0, (c.width+c.height)*2)
	for i := 0; i < c.height; i++ {
		c.sendBeam(vec{0, i}, right)
		energized = append(energized, c.countEnergized())
		c.reset()

		c.sendBeam(vec{c.width - 1, i}, left)
		energized = append(energized, c.countEnergized())
		c.reset()
	}
	for i := 0; i < c.width; i++ {
		c.sendBeam(vec{i, 0}, down)
		energized = append(energized, c.countEnergized())
		c.reset()

		c.sendBeam(vec{i, c.height - 1}, up)
		energized = append(energized, c.countEnergized())
		c.reset()
	}
	return slices.Max(energized)
}

func main() {
	input := util.StdinReadlines()
	c := makeCave(input)
	fmt.Println(c.maxEnergized())
}
