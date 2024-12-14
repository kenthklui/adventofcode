package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

var AREA_WIDTH, AREA_HEIGHT = 101, 103
var DURATION = 100

type vec2 struct{ x, y int }

func (v vec2) fma(v2 vec2, s int) vec2  { return vec2{v.x + v2.x*s, v.y + v2.y*s} }
func (v vec2) in(sizeX, sizeY int) bool { return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY }
func (v vec2) teleportInbound() vec2 {
	return vec2{
		(v.x%AREA_WIDTH + AREA_WIDTH) % AREA_WIDTH,
		(v.y%AREA_HEIGHT + AREA_HEIGHT) % AREA_HEIGHT,
	}
}
func (v vec2) quadrant() int {
	midX, midY := AREA_WIDTH/2, AREA_HEIGHT/2
	if v.x == midX || v.y == midY {
		return -1
	}
	quadrant := 0
	if v.x > midX {
		quadrant += 1
	}
	if v.y > midY {
		quadrant += 2
	}
	return quadrant
}

type robot struct {
	pos, vel vec2
}

func (r robot) quadrant() int { return r.pos.quadrant() }

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

func solve(input []string) (output string) {
	quadrantCount := make([]int, 4)
	for _, line := range input {
		robot := newRobot(line)
		robot.move(DURATION)
		if robot.quadrant() != -1 {
			quadrantCount[robot.quadrant()]++
		}
	}
	safetyFactor := quadrantCount[0] * quadrantCount[1] * quadrantCount[2] * quadrantCount[3]
	return strconv.Itoa(safetyFactor)
}

func main() {
	// AREA_WIDTH, AREA_HEIGHT = 11, 7
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
