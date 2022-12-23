package main

import (
	"bufio"
	"fmt"
	"os"
)

func intMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

type position struct {
	x, y int
}

func (p position) north() position { return position{p.x, p.y + 1} }
func (p position) south() position { return position{p.x, p.y - 1} }
func (p position) west() position  { return position{p.x - 1, p.y} }
func (p position) east() position  { return position{p.x + 1, p.y} }

type floorMap struct {
	elves     []*elf
	floor     map[position]*elf
	firstCons int
}

func NewFloorMap() *floorMap {
	return &floorMap{
		elves:     make([]*elf, 0),
		floor:     make(map[position]*elf),
		firstCons: 0,
	}
}

type elf struct {
	pos position
}

func (fm *floorMap) round() bool {
	proposals := fm.consider()
	fm.firstCons = (fm.firstCons + 1) % 4
	return fm.execute(proposals)
}

func (fm *floorMap) consider() map[position][]*elf {
	proposals := make(map[position][]*elf)

	for _, e := range fm.elves {
		if fm.considerEmpty(e.pos.x, e.pos.y) {
			continue
		}
		cons := fm.firstCons
	Loop:
		for i := 0; i < 4; i++ {
			switch (cons + i) % 4 {
			case 0:
				if fm.considerNorth(e.pos.x, e.pos.y) {
					val, ok := proposals[e.pos.north()]
					if !ok {
						val = []*elf{}
					}
					proposals[e.pos.north()] = append(val, e)
					break Loop
				}
			case 1:
				if fm.considerSouth(e.pos.x, e.pos.y) {
					val, ok := proposals[e.pos.south()]
					if !ok {
						val = []*elf{}
					}
					proposals[e.pos.south()] = append(val, e)
					break Loop
				}
			case 2:
				if fm.considerWest(e.pos.x, e.pos.y) {
					val, ok := proposals[e.pos.west()]
					if !ok {
						val = []*elf{}
					}
					proposals[e.pos.west()] = append(val, e)
					break Loop
				}
			case 3:
				if fm.considerEast(e.pos.x, e.pos.y) {
					val, ok := proposals[e.pos.east()]
					if !ok {
						val = []*elf{}
					}
					proposals[e.pos.east()] = append(val, e)
					break Loop
				}
			}
		}
	}

	return proposals
}

func (fm *floorMap) execute(proposals map[position][]*elf) bool {
	moved := false
	for pos, elves := range proposals {
		if len(elves) > 1 {
			continue
		}

		fm.floor[elves[0].pos] = nil
		fm.floor[pos] = elves[0]
		elves[0].pos = pos
		moved = true
	}
	return moved
}

func (fm *floorMap) considerEmpty(x, y int) bool {
	n := position{x, y + 1}
	ne := position{x + 1, y + 1}
	nw := position{x - 1, y + 1}
	w := position{x - 1, y}
	e := position{x + 1, y}
	s := position{x, y - 1}
	se := position{x + 1, y - 1}
	sw := position{x - 1, y - 1}
	for _, pos := range []position{n, ne, nw, w, e, s, se, sw} {
		if val, ok := fm.floor[pos]; ok && val != nil {
			return false
		}
	}
	return true
}

func (fm *floorMap) considerNorth(x, y int) bool {
	n := position{x, y + 1}
	ne := position{x + 1, y + 1}
	nw := position{x - 1, y + 1}
	for _, pos := range []position{n, ne, nw} {
		if val, ok := fm.floor[pos]; ok && val != nil {
			return false
		}
	}
	return true
}

func (fm *floorMap) considerSouth(x, y int) bool {
	s := position{x, y - 1}
	se := position{x + 1, y - 1}
	sw := position{x - 1, y - 1}
	for _, pos := range []position{s, se, sw} {
		if val, ok := fm.floor[pos]; ok && val != nil {
			return false
		}
	}
	return true
}

func (fm *floorMap) considerWest(x, y int) bool {
	w := position{x - 1, y}
	nw := position{x - 1, y + 1}
	sw := position{x - 1, y - 1}
	for _, pos := range []position{w, nw, sw} {
		if val, ok := fm.floor[pos]; ok && val != nil {
			return false
		}
	}
	return true
}

func (fm *floorMap) considerEast(x, y int) bool {
	e := position{x + 1, y}
	ne := position{x + 1, y + 1}
	se := position{x + 1, y - 1}
	for _, pos := range []position{e, ne, se} {
		if val, ok := fm.floor[pos]; ok && val != nil {
			return false
		}
	}
	return true
}

func (fm *floorMap) emptySpace() int {
	var minX, minY, maxX, maxY int
	for _, e := range fm.elves {
		minX, _ = intMinMax(e.pos.x, minX)
		minY, _ = intMinMax(e.pos.y, minY)
		_, maxX = intMinMax(e.pos.x, maxX)
		_, maxY = intMinMax(e.pos.y, maxY)
	}
	return (maxX-minX+1)*(maxY-minY+1) - len(fm.elves)
}

func parseInput(input []string) *floorMap {
	fm := NewFloorMap()
	for dy, line := range input {
		for x, r := range line {
			if r == '#' {
				pos := position{x, -dy}
				e := &elf{pos}
				fm.floor[pos] = e
				fm.elves = append(fm.elves, e)
			}
		}
	}
	return fm
}

func readInput() []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return lines
}

func main() {
	input := readInput()
	fm := parseInput(input)
	for i := 0; true; i++ {
		if !fm.round() {
			fmt.Println(i + 1)
			return
		}
	}
}
