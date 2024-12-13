package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct{ x, y int }

func (v vec2) sum(u vec2) vec2         { return vec2{v.x + u.x, v.y + u.y} }
func (v vec2) sumK(u vec2, k int) vec2 { return vec2{v.x + u.x*k, v.y + u.y*k} }
func (v vec2) in(input []string) bool {
	return 0 <= v.x && v.x < len(input[0]) && 0 <= v.y && v.y < len(input)
}

type requirements struct {
	layouts [][]vec2
	chars   []byte
}

var XMAS requirements

func xmas() requirements {
	if XMAS.layouts == nil {
		layouts := make([][]vec2, 0, 4)
		axes := [][]vec2{
			{vec2{1, 0}, vec2{0, 1}},   // right, down
			{vec2{0, 1}, vec2{-1, 0}},  // down, left
			{vec2{-1, 0}, vec2{0, -1}}, // left, up
			{vec2{0, -1}, vec2{1, 0}},  // up, right
		}
		zero := vec2{0, 0}

		for _, axis := range axes {
			layouts = append(layouts, []vec2{
				// 0 1
				//  2
				// 3 4
				zero, zero.sumK(axis[0], 2),
				zero.sum(axis[0]).sum(axis[1]),
				zero.sumK(axis[1], 2), zero.sumK(axis[0], 2).sumK(axis[1], 2),
			})
		}

		XMAS = requirements{layouts, []byte{'M', 'S', 'A', 'M', 'S'}}
	}
	return XMAS
}

func checkXMAS(input []string, start vec2, layout []vec2) bool {
	if !start.sum(layout[4]).in(input) {
		return false
	}

	for i, pos := range layout {
		v := start.sum(pos)
		if input[v.y][v.x] != xmas().chars[i] {
			return false
		}
	}
	return true
}

func solve(input []string) (output string) {
	times := 0
	for y, line := range input {
		for x := range line {
			for _, layout := range xmas().layouts {
				if checkXMAS((input), vec2{x, y}, layout) {
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
