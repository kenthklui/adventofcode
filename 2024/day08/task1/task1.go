package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct {
	x, y int
}

func (v vec2) in(sizeX, sizeY int) bool { return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY }
func (v vec2) add(v2 vec2) vec2         { return vec2{v.x + v2.x, v.y + v2.y} }
func (v vec2) sub(v2 vec2) vec2         { return vec2{v.x - v2.x, v.y - v2.y} }

type nodeMap struct {
	sizeX, sizeY int
	antennas     map[string][]vec2
	antinodes    map[vec2]bool
}

func newNodeMap(input []string) nodeMap {
	antennas := make(map[string][]vec2)
	antinodes := make(map[vec2]bool)
	sizeX, sizeY := len(input[0]), len(input)
	for y, line := range input {
		for x, char := range line {
			if char != '.' {
				antennas[string(char)] = append(antennas[string(char)], vec2{x, y})
			}
		}
	}
	return nodeMap{sizeX, sizeY, antennas, antinodes}
}

func (nm nodeMap) markAntinodes(v1, v2 vec2) {
	diff := v2.sub(v1)
	an1, an2 := v1.sub(diff), v2.add(diff)
	if an1.in(nm.sizeX, nm.sizeY) {
		nm.antinodes[an1] = true
	}
	if an2.in(nm.sizeX, nm.sizeY) {
		nm.antinodes[an2] = true
	}
}

func (nm nodeMap) markAllAntinodes() {
	for _, antennas := range nm.antennas {
		for i := range antennas[1:] {
			for j := range antennas[i+1:] {
				nm.markAntinodes(antennas[i], antennas[i+j+1])
			}
		}
	}
}

func solve(input []string) (output string) {
	nm := newNodeMap(input)
	nm.markAllAntinodes()
	return fmt.Sprint(len(nm.antinodes))
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
