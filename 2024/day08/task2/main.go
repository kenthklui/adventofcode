package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func gcd(a, b int) int {
	if a < b {
		a, b = b, a
	}
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type vec2 struct {
	x, y int
}

func (v vec2) in(sizeX, sizeY int) bool { return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY }
func (v vec2) add(v2 vec2) vec2         { return vec2{v.x + v2.x, v.y + v2.y} }
func (v vec2) sub(v2 vec2) vec2         { return vec2{v.x - v2.x, v.y - v2.y} }

func (v vec2) reduce() vec2 {
	if v.x == 0 {
		v.y /= abs(v.y)
	} else if v.y == 0 {
		v.x /= abs(v.x)
	} else {
		gcd := gcd(abs(v.x), abs(v.y))
		v.x, v.y = v.x/gcd, v.y/gcd
	}
	return v
}

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
	nm.antinodes[v1], nm.antinodes[v2] = true, true
	diff := v2.sub(v1).reduce()

	for cand := v1.sub(diff); cand.in(nm.sizeX, nm.sizeY); cand = cand.sub(diff) { // Backwards from v1
		nm.antinodes[cand] = true
	}
	for cand := v1.add(diff); cand.x != v2.x; cand = cand.add(diff) { // Between v1 and v2
		nm.antinodes[cand] = true
	}
	for cand := v2.add(diff); cand.in(nm.sizeX, nm.sizeY); cand = cand.add(diff) { // Forwards from v2
		nm.antinodes[cand] = true
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
