package main

import (
	"fmt"
	"slices"

	"github.com/kenthklui/adventofcode/util"
)

const int32mask int = (1 << 32) - 1

type position struct {
	x, y int
}

func (p position) add(dir position) position { return position{p.x + dir.x, p.y + dir.y} }
func (p position) key() int                  { return (p.x&int32mask)<<32 + (p.y & int32mask) }

var (
	n            = position{0, -1}
	s            = position{0, 1}
	w            = position{-1, 0}
	e            = position{1, 0}
	nw           = position{-1, -1}
	ne           = position{1, -1}
	sw           = position{-1, 1}
	se           = position{1, 1}
	dir          = [8]position{n, s, w, e, nw, ne, sw, se}
	clear  int64 = 0b11111111
	nValid int64 = 0b10001100
	sValid int64 = 0b01000011
	wValid int64 = 0b00101010
	eValid int64 = 0b00010101
)

type elf struct {
	pos position
}

type claim struct {
	e      *elf
	newPos position
	valid  bool
}

type proposals map[int]*claim

type floorMap struct {
	elves     []*elf
	floor     map[int]*elf
	firstCons int
	props     proposals
}

func NewFloorMap(elves []*elf) *floorMap {
	fm := floorMap{
		elves:     elves,
		floor:     make(map[int]*elf, len(elves)),
		firstCons: 0,
		props:     make(proposals, len(elves)),
	}
	for _, e := range elves {
		fm.floor[e.pos.key()] = e
	}
	return &fm
}

func (fm *floorMap) round() bool {
	fm.consider()
	moves := fm.execute()
	fm.firstCons = (fm.firstCons + 1) % 4
	return moves > 0
}

func (fm *floorMap) options(e *elf) [5]bool {
	var flag int64
	for _, d := range dir {
		pos := e.pos.add(d)
		flag <<= 1
		if _, ok := fm.floor[pos.key()]; !ok {
			flag++
		}
	}
	return [5]bool{
		flag&nValid == nValid,
		flag&sValid == sValid,
		flag&wValid == wValid,
		flag&eValid == eValid,
		flag == clear,
	}
}

func (fm *floorMap) consider() {
	for _, e := range fm.elves {
		valid := fm.options(e)
		if valid[4] {
			continue
		}

		for j := 0; j < 4; j++ {
			cons := (fm.firstCons + j) % 4
			if valid[cons] {
				newPos := e.pos.add(dir[cons])
				newPosKey := newPos.key()
				if cl, ok := fm.props[newPosKey]; ok {
					cl.valid = false
				} else {
					fm.props[newPosKey] = &claim{e, newPos, true}
				}
				break
			}
		}
	}
}

func (fm *floorMap) execute() int {
	moves := 0
	for newPosKey, cl := range fm.props {
		if cl.valid {
			delete(fm.floor, cl.e.pos.key())
			fm.floor[newPosKey] = cl.e
			cl.e.pos = cl.newPos
			moves++
		}

		delete(fm.props, newPosKey)
	}
	return moves
}

func (fm *floorMap) emptySpace() int {
	x, y := make([]int, len(fm.elves)), make([]int, len(fm.elves))
	for i, e := range fm.elves {
		x[i], y[i] = e.pos.x, e.pos.y
	}
	minX, maxX := slices.Min(x), slices.Max(x)
	minY, maxY := slices.Min(y), slices.Max(y)
	return (maxX-minX+1)*(maxY-minY+1) - len(fm.elves)
}

func parseInput(input []string) *floorMap {
	elves := make([]*elf, 0)
	for y, line := range input {
		for x, r := range line {
			if r == '#' {
				pos := position{x, y}
				e := &elf{pos}
				elves = append(elves, e)
			}
		}
	}
	return NewFloorMap(elves)
}

func main() {
	input := util.StdinReadlines()
	fm := parseInput(input)
	i := 1
	for fm.round() {
		i++
	}
	fmt.Println(i)
}
