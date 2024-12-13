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
	arenaDir     [][][4]bool
}

func newGuard(x, y int, input []string) *guard {
	pos := vec2{x, y}
	arena := make([][]byte, len(input))
	arenaDir := make([][][4]bool, len(input))
	for y, line := range input {
		arena[y] = []byte(line)
		arenaDir[y] = make([][4]bool, len(line))
	}
	arenaDir[y][x][0] = true
	g := guard{pos, 0, 1, arena, arenaDir}
	return &g
}

func newGuardWithObstacle(x, y, obsX, obsY int, input []string) *guard {
	g := newGuard(x, y, input)
	g.arena[obsY][obsX] = '#'
	return g
}

// 2: OOB, 1: loop, 0: moved
func (g *guard) checkLoop() int {
	if g.arenaDir[g.pos.y][g.pos.x][g.dir] {
		return 1
	} else {
		g.arenaDir[g.pos.y][g.pos.x][g.dir] = true
		return 0
	}
}

func (g *guard) move() int {
	next := g.pos.add(dirs[g.dir])
	if !next.in(len(g.arena[0]), len(g.arena)) {
		return 2
	} else if g.arena[next.y][next.x] == '#' {
		g.dir = (g.dir + 1) % 4
		return g.checkLoop()
	} else {
		g.pos = next
		if g.arena[g.pos.y][g.pos.x] == '.' {
			g.arena[g.pos.y][g.pos.x] = 'X'
			g.covered++
		}
		return g.checkLoop()
	}
}

func (g *guard) run() int {
	m := 0
	for m == 0 {
		m = g.move()
	}
	return m
}

func solve(input []string) (output string) {
	var guardX, guardY int
	for y, line := range input {
		if index := strings.IndexByte(line, '^'); index != -1 {
			input[y] = strings.ReplaceAll(line, "^", "X")
			guardX, guardY = index, y
		}
	}

	options := 0
	for obsY, line := range input {
		for obsX := range line {
			if input[obsY][obsX] == '.' {
				g := newGuardWithObstacle(guardX, guardY, obsX, obsY, input)
				if g.run() == 1 {
					options++
				}
			}
		}
	}
	return fmt.Sprintf("%d", options)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
