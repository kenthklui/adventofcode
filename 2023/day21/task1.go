package main

import (
	"fmt"

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
var dirs = []vec{up, down, left, right}

func (v vec) add(dir vec) vec { return vec{v.x + dir.x, v.y + dir.y} }

type garden struct {
	width, height int
	start         vec
	grid          map[vec]byte
}

func parseGarden(input []string) garden {
	width, height := len(input[0]), len(input)
	var start vec
	grid := make(map[vec]byte)
	for y, row := range input {
		for x, c := range row {
			if c == 'S' {
				start = vec{x, y}
				c = '.'
			}
			grid[vec{x, y}] = byte(c)
		}
	}
	return garden{width, height, start, grid}
}

func (g garden) inBounds(v vec) bool {
	return v.x >= 0 && v.y >= 0 && v.x < g.width && v.y < g.height
}

type agent struct {
	v     vec
	steps int
}

func (g garden) reachable(steps int) int {
	travelled := make(map[vec]int)
	queue := make([]agent, 0, g.width*g.height)
	queue = append(queue, agent{g.start, 0})
	for len(queue) > 0 {
		a := queue[0]
		queue = queue[1:]

		travelled[a.v] = 2
		if a.steps < steps {
			for _, d := range dirs {
				if n := a.v.add(d); g.inBounds(n) {
					if travelled[n] > 0 {
						continue
					}
					if g.grid[n] == '#' {
						continue
					}

					travelled[n] = 1
					queue = append(queue, agent{n, a.steps + 1})
				}
			}
		}
	}
	squareCount := 0
	mod2 := (g.start.x + g.start.y + steps) % 2
	for t := range travelled {
		if (t.x+t.y)%2 == mod2 {
			squareCount++
		}
	}
	return squareCount
}

// const maxSteps = 6

const maxSteps = 64

func main() {
	input := util.StdinReadlines()
	garden := parseGarden(input)
	fmt.Println(garden.reachable(maxSteps))
}
