package main

import (
	"bufio"
	"fmt"
	"os"
)

type sc struct {
	x, y int
}

type seaMap struct {
	width, length     int
	floor             [][]*sc
	easties, southies []*sc
}

func NewSeaMap(width, length int) *seaMap {
	floor := make([][]*sc, length)
	for i := range floor {
		floor[i] = make([]*sc, width)
	}

	sm := seaMap{
		width:    width,
		length:   length,
		floor:    floor,
		easties:  []*sc{},
		southies: []*sc{},
	}

	return &sm
}

func (sm *seaMap) addEastie(x, y int) {
	if sm.floor[y][x] != nil {
		panic("Adding to an occupied square")
	}
	e := &sc{x, y}
	sm.floor[y][x] = e
	sm.easties = append(sm.easties, e)
}

func (sm *seaMap) canMoveEast(cu *sc) bool {
	newX := (cu.x + 1) % sm.width
	return sm.floor[cu.y][newX] == nil
}

func (sm *seaMap) moveEast(cu *sc) {
	newX := (cu.x + 1) % sm.width
	sm.floor[cu.y][cu.x] = nil
	sm.floor[cu.y][newX] = cu
	cu.x = newX
}

func (sm *seaMap) moveEasties() int {
	canMove := make([]bool, len(sm.easties))
	for i, e := range sm.easties {
		canMove[i] = sm.canMoveEast(e)
	}

	moved := 0
	for i, e := range sm.easties {
		if canMove[i] {
			sm.moveEast(e)
			moved++
		}
	}
	return moved
}

func (sm *seaMap) addSouthie(x, y int) {
	if sm.floor[y][x] != nil {
		panic("Adding to an occupied square")
	}
	e := &sc{x, y}
	sm.floor[y][x] = e
	sm.southies = append(sm.southies, e)
}

func (sm *seaMap) canMoveSouth(cu *sc) bool {
	newY := (cu.y + 1) % sm.length
	return sm.floor[newY][cu.x] == nil
}

func (sm *seaMap) moveSouth(cu *sc) {
	newY := (cu.y + 1) % sm.length
	sm.floor[cu.y][cu.x] = nil
	sm.floor[newY][cu.x] = cu
	cu.y = newY
}

func (sm *seaMap) moveSouthies() int {
	canMove := make([]bool, len(sm.southies))
	for i, s := range sm.southies {
		canMove[i] = sm.canMoveSouth(s)
	}

	moved := 0
	for i, s := range sm.southies {
		if canMove[i] {
			sm.moveSouth(s)
			moved++
		}
	}
	return moved
}

func (sm *seaMap) step() int {
	moved := 0
	moved += sm.moveEasties()
	moved += sm.moveSouthies()
	return moved
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

func parseInput(input []string) *seaMap {
	length := len(input)
	width := len(input[0])

	sm := NewSeaMap(width, length)
	for y, line := range input {
		for x, r := range line {
			switch r {
			case '.':
				continue
			case '>':
				sm.addEastie(x, y)
			case 'v':
				sm.addSouthie(x, y)
			}
		}
	}

	return sm
}

func main() {
	input := readInput()
	sm := parseInput(input)

	steps := 1
	moved := sm.step()
	for moved > 0 {
		steps++
		moved = sm.step()
	}

	fmt.Println(steps)
}
