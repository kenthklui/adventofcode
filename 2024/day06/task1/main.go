package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct {
	x, y int
}

func (v vec2) in(sizeX, sizeY int) bool {
	return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY
}

func (v vec2) add(v2 vec2) vec2 {
	return vec2{v.x + v2.x, v.y + v2.y}
}

var dirs = []vec2{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type guard struct {
	pos          vec2
	dir, covered int
	arena        [][]byte
}

func newGuard(x, y int, input []string) *guard {
	pos := vec2{x, y}
	arena := make([][]byte, len(input))
	for y, line := range input {
		arena[y] = []byte(line)
	}
	g := guard{pos, 0, 1, arena}
	return &g
}

func (g *guard) move() bool {
	next := g.pos.add(dirs[g.dir])
	if !next.in(len(g.arena[0]), len(g.arena)) {
		return false
	} else if g.arena[next.y][next.x] == '#' {
		g.dir = (g.dir + 1) % 4
		return true
	} else {
		g.pos = next
		if g.arena[g.pos.y][g.pos.x] == '.' {
			g.arena[g.pos.y][g.pos.x] = 'X'
			g.covered++
		}
		return true
	}
}

func (g *guard) run() int {
	for g.move() {
	}
	return g.covered
}

func solve(input []string) (output string) {
	var g *guard
	for y, line := range input {
		if index := strings.IndexByte(line, '^'); index != -1 {
			input[y] = strings.ReplaceAll(line, "^", "X")
			g = newGuard(index, y, input)
			break
		}
	}
	return fmt.Sprintf("%d", g.run())
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
