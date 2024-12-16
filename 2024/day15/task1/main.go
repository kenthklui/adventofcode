package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct{ x, y int }

func (v vec2) add(v2 vec2) vec2 { return vec2{v.x + v2.x, v.y + v2.y} }

var dirs = []vec2{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type object int

const (
	wall object = iota
	boulder
)

type warehouse struct {
	robot   vec2
	objects map[vec2]object
}

func newWarehouse(input []string) *warehouse {
	warehouse := warehouse{
		objects: make(map[vec2]object),
	}
	for y, line := range input {
		for x, char := range line {
			switch char {
			case '#':
				warehouse.objects[vec2{x, y}] = wall
			case 'O':
				warehouse.objects[vec2{x, y}] = boulder
			case '@':
				warehouse.robot = vec2{x, y}
			}
		}
	}
	return &warehouse
}

func (wh *warehouse) moveRobot(dir int) bool {
	newPos := wh.robot.add(dirs[dir])
	if wh.moveObject(newPos, dir) {
		wh.robot = newPos
		return true
	} else {
		return false
	}
}

func (wh *warehouse) moveObject(pos vec2, dir int) bool {
	if obj, ok := wh.objects[pos]; ok {
		switch obj {
		case wall:
			return false
		case boulder:
			newPos := pos.add(dirs[dir])
			if wh.moveObject(newPos, dir) {
				wh.objects[newPos] = boulder
				delete(wh.objects, pos)
				return true
			} else {
				return false
			}
		default:
			panic("unknown object")
		}
	} else {
		return true
	}
}

func (wh *warehouse) gps() int {
	gps := 0
	for pos, obj := range wh.objects {
		if obj == boulder {
			gps += pos.x + pos.y*100
		}
	}
	return gps
}

func parseMoves(input []string) (moves []int) {
	for _, line := range input {
		for _, char := range line {
			switch char {
			case '^':
				moves = append(moves, 0)
			case '>':
				moves = append(moves, 1)
			case 'v':
				moves = append(moves, 2)
			case '<':
				moves = append(moves, 3)
			}
		}
	}
	return
}

func solve(input []string) (output string) {
	var wh *warehouse
	var moves []int

	for y, line := range input {
		if len(line) == 0 {
			wh = newWarehouse(input[:y])
			moves = parseMoves(input[y+1:])
			break
		}
	}

	for _, move := range moves {
		wh.moveRobot(move)
	}

	return strconv.Itoa(wh.gps())
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
