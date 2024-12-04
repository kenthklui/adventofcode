package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct{ x, y int }

func (v vec2) sum(u vec2) vec2 { return vec2{v.x + u.x, v.y + u.y} }
func (v vec2) in(input []string) bool {
	return 0 <= v.x && v.x < len(input[0]) && 0 <= v.y && v.y < len(input)
}

var DIRS []vec2
var XMAS = "XMAS"

func dirs() []vec2 {
	if DIRS == nil {
		DIRS = make([]vec2, 0, 8)
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				DIRS = append(DIRS, vec2{dx, dy})
			}
		}
	}
	return DIRS
}

func checkXMAS(input []string, start, dir vec2) bool {
	for i, pos := 0, start; i < len(XMAS); i, pos = i+1, pos.sum(dir) {
		if !pos.in(input) || input[pos.y][pos.x] != XMAS[i] {
			return false
		}
	}
	return true
}

func solve(input []string) (output string) {
	times := 0
	for y, line := range input {
		for x := range line {
			for _, dir := range dirs() {
				if checkXMAS(input, vec2{x, y}, dir) {
					times++
				}
			}
		}
	}
	return fmt.Sprintf("%d", times)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
