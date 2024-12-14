package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

var AREA_WIDTH, AREA_HEIGHT = 101, 103
var DURATION = 10000

type vec2 struct{ x, y int }

func (v vec2) fma(v2 vec2, s int) vec2  { return vec2{v.x + v2.x*s, v.y + v2.y*s} }
func (v vec2) in(sizeX, sizeY int) bool { return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY }
func (v vec2) neighbors() [8]vec2 {
	return [8]vec2{
		{v.x - 1, v.y - 1}, {v.x, v.y - 1}, {v.x + 1, v.y - 1},
		{v.x - 1, v.y}, {v.x + 1, v.y},
		{v.x - 1, v.y + 1}, {v.x, v.y + 1}, {v.x + 1, v.y + 1},
	}
}
func (v vec2) teleportInbound() vec2 {
	return vec2{
		(v.x%AREA_WIDTH + AREA_WIDTH) % AREA_WIDTH,
		(v.y%AREA_HEIGHT + AREA_HEIGHT) % AREA_HEIGHT,
	}
}

type robot struct {
	pos, vel vec2
}

func newRobot(line string) *robot {
	ints := util.ParseLineInts(line)
	return &robot{vec2{ints[0], ints[1]}, vec2{ints[2], ints[3]}}
}

func (r *robot) move(seconds int) {
	r.pos = r.pos.fma(r.vel, seconds)
	if !r.pos.in(AREA_WIDTH, AREA_HEIGHT) {
		r.pos = r.pos.teleportInbound()
	}
}

type floorMap struct {
	robots []*robot
	floor  [][]bool
}

func newFloorMap(robots []*robot) *floorMap {
	floor := make([][]bool, AREA_HEIGHT)
	for i := range floor {
		floor[i] = make([]bool, AREA_WIDTH)
	}
	return &floorMap{robots, floor}
}

func (fm *floorMap) move() {
	for _, robot := range fm.robots {
		fm.floor[robot.pos.y][robot.pos.x] = false
		robot.move(1)
		fm.floor[robot.pos.y][robot.pos.x] = true
	}
}

// A floor map is a potential solution if "most", ie. 2/3rds of robots have at least 1 neighbor
func (fm floorMap) treeCandidate() bool {
	hasNearby := 0
	for _, robot := range fm.robots {
		nearby := 0
		for _, neighbor := range robot.pos.neighbors() {
			if neighbor.in(AREA_WIDTH, AREA_HEIGHT) && fm.floor[neighbor.y][neighbor.x] {
				nearby++
			}
		}
		if nearby > 0 {
			hasNearby++
		}
	}
	return hasNearby >= len(fm.robots)*2/3
}

func (fm floorMap) String() string {
	return strings.Join(util.ConvertScreen(fm.floor), "\n")
}

func solve(input []string) (output string) {
	robots := make([]*robot, len(input))
	for i, line := range input {
		robots[i] = newRobot(line)
	}

	fm := newFloorMap(robots)
	for seconds := 1; seconds <= DURATION; seconds++ {
		fm.move()
		if fm.treeCandidate() {
			return fmt.Sprintf("Seconds: %d\n%s", seconds, fm)
		}
	}

	return fmt.Sprintf("Not found within %d seconds", DURATION)
}

func main() {
	// AREA_WIDTH, AREA_HEIGHT = 11, 7
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
