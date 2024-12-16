package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var empty void

type vec2 struct{ x, y int }

func (v vec2) add(v2 vec2) vec2 { return vec2{v.x + v2.x, v.y + v2.y} }

var dirs = []vec2{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type objectType int

const (
	wall objectType = iota
	boulder
)

type object struct {
	typ objectType
	pos [2]vec2
}

type warehouse struct {
	robot     vec2
	objects   []*object
	objectMap map[vec2]*object
}

func newWarehouse(input []string) *warehouse {
	wh := &warehouse{
		objects:   make([]*object, 0),
		objectMap: make(map[vec2]*object),
	}
	for y, line := range input {
		for x, char := range line {
			switch char {
			case '#':
				wh.addObj(wall, x, y)
			case 'O':
				wh.addObj(boulder, x, y)
			case '@':
				wh.robot = vec2{2 * x, y}
			}
		}
	}
	return wh
}

func (wh *warehouse) addObj(typ objectType, x, y int) {
	obj := &object{typ: typ, pos: [2]vec2{{x * 2, y}, {x*2 + 1, y}}}
	wh.objects = append(wh.objects, obj)
	wh.objectMap[vec2{x * 2, y}] = obj
	wh.objectMap[vec2{x*2 + 1, y}] = obj
}

func (wh *warehouse) moveRobot(dirIndex int) bool {
	dir := dirs[dirIndex]
	newPos := wh.robot.add(dir)

	objMoveList := make(map[*object]void)
	if obj, ok := wh.objectMap[newPos]; ok {
		if !wh.canMoveObject(obj, dir, objMoveList) {
			return false
		}
	}

	wh.moveObjects(objMoveList, dir)
	wh.robot = newPos
	return true
}

func (wh *warehouse) canMoveObject(obj *object, dir vec2, objMoveList map[*object]void) bool {
	if _, seen := objMoveList[obj]; seen {
		return true
	}

	switch obj.typ {
	case wall:
		return false
	case boulder:
		for _, pos := range obj.pos {
			if nextObj, ok := wh.objectMap[pos.add(dir)]; ok && nextObj != obj {
				if !wh.canMoveObject(nextObj, dir, objMoveList) {
					return false
				}
			}
		}
		objMoveList[obj] = empty
		return true
	default:
		panic("unknown object")
	}
}

func (wh *warehouse) moveObjects(objMoveList map[*object]void, dir vec2) {
	for obj := range objMoveList {
		delete(wh.objectMap, obj.pos[0])
		delete(wh.objectMap, obj.pos[1])
	}

	for obj := range objMoveList {
		nextPos := [2]vec2{obj.pos[0].add(dir), obj.pos[1].add(dir)}
		obj.pos = nextPos
		wh.objectMap[nextPos[0]] = obj
		wh.objectMap[nextPos[1]] = obj
	}
}

func (wh *warehouse) gps() int {
	gps := 0
	for _, obj := range wh.objects {
		if obj.typ == boulder {
			gps += obj.pos[0].x + obj.pos[1].y*100
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
