package main

import (
	"fmt"
	"slices"

	"github.com/kenthklui/adventofcode/util"
)

type vec struct{ x, y int }

var up = vec{0, -1}
var right = vec{1, 0}
var down = vec{0, 1}
var left = vec{-1, 0}

func (v vec) add(d vec) vec { return vec{v.x + d.x, v.y + d.y} }
func (v vec) inBounds(width, height int) bool {
	return v.x >= 0 && v.y >= 0 && v.x < width && v.y < height
}

type area struct {
	start, end    vec
	width, height int
	aMap          []string
}

func parseArea(input []string) area {
	return area{
		start:  vec{1, 0},
		end:    vec{len(input[0]) - 2, len(input) - 1},
		width:  len(input[0]),
		height: len(input),
		aMap:   input,
	}
}

func (a area) longestHike() int {
	travelled := make([][]bool, a.height)
	for i := range travelled {
		travelled[i] = make([]bool, a.width)
	}
	if reached, steps := a.recurseHike(a.start, 0, travelled); reached {
		return steps
	} else {
		panic("Maze cannot be completed")
	}
}

func (a area) recurseHike(curr vec, steps int, travelled [][]bool) (bool, int) {
	if curr == a.end {
		return true, steps
	}

	travelled[curr.y][curr.x] = true

	var dirs []vec
	switch a.aMap[curr.y][curr.x] {
	case '^':
		dirs = []vec{up}
	case '>':
		dirs = []vec{right}
	case 'v':
		dirs = []vec{down}
	case '<':
		dirs = []vec{left}
	default:
		dirs = []vec{up, right, down, left}
	}

	longestHikes := make([]int, 0, len(dirs))
	for _, dir := range dirs {
		next := curr.add(dir)
		if !next.inBounds(a.width, a.height) {
			continue
		}

		if travelled[next.y][next.x] {
			continue
		}

		switch a.aMap[next.y][next.x] {
		case '#':
			continue
		case '.':
		case '^':
			if dir == down {
				continue
			}
		case '>':
			if dir == left {
				continue
			}
		case 'v':
			if dir == up {
				continue
			}
		case '<':
			if dir == right {
				continue
			}
		default:
			panic("Invalid map character")
		}
		if reached, steps := a.recurseHike(next, steps+1, travelled); reached {
			longestHikes = append(longestHikes, steps)
		}
	}
	travelled[curr.y][curr.x] = false

	if len(longestHikes) > 0 {
		return true, slices.Max(longestHikes)
	} else {
		return false, steps
	}
}

func main() {
	input := util.StdinReadlines()
	a := parseArea(input)
	fmt.Println(a.longestHike())
}
