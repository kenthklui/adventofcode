package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func intMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

type point struct {
	x, y, z int
}

func (p point) neighbors() []point {
	return []point{
		point{p.x - 1, p.y, p.z}, point{p.x + 1, p.y, p.z},
		point{p.x, p.y - 1, p.z}, point{p.x, p.y + 1, p.z},
		point{p.x, p.y, p.z - 1}, point{p.x, p.y, p.z + 1},
	}
}

type cube struct {
	val                    [][][]bool
	startX, startY, startZ int
}

func NewCube(startX, startY, startZ, endX, endY, endZ int) *cube {
	sizeX, sizeY, sizeZ := endX-startX+1, endY-startY+1, endZ-startZ+1
	val := make([][][]bool, sizeX)
	for x := range val {
		val[x] = make([][]bool, sizeY)
		for y := range val[x] {
			val[x][y] = make([]bool, sizeZ)
		}
	}

	return &cube{val: val, startX: startX, startY: startY, startZ: startZ}
}

func (c *cube) get(p point) (bool, error) {
	cx, cy, cz := p.x-c.startX, p.y-c.startY, p.z-c.startZ
	if cx < 0 || cx >= len(c.val) ||
		cy < 0 || cy >= len(c.val[cx]) ||
		cz < 0 || cz >= len(c.val[cx][cy]) {
		return false, fmt.Errorf("Out of bounds")
	}
	return c.val[cx][cy][cz], nil
}

func (c *cube) set(p point, val bool) error {
	cx, cy, cz := p.x-c.startX, p.y-c.startY, p.z-c.startZ
	if cx < 0 || cx >= len(c.val) ||
		cy < 0 || cy >= len(c.val[cx]) ||
		cz < 0 || cz >= len(c.val[cx][cy]) {
		return fmt.Errorf("Out of bounds")
	}

	c.val[cx][cy][cz] = val
	return nil
}

func surfaceArea(lava *cube, points []point) int {
	sides := 0
	for _, p := range points {
		sides += 6
		for _, n := range p.neighbors() {
			if hit, err := lava.get(n); err == nil && hit {
				sides--
			}
		}
	}

	return sides
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

func parseInput(input []string) (*cube, []point) {
	minX, minY, minZ := 1000, 1000, 1000
	maxX, maxY, maxZ := -1000, -1000, -1000

	points := make([]point, 0, len(input))
	for _, text := range input {
		var p point
		splits := strings.Split(text, ",")
		p.x, _ = strconv.Atoi(splits[0])
		p.y, _ = strconv.Atoi(splits[1])
		p.z, _ = strconv.Atoi(splits[2])
		points = append(points, p)

		minX, _ = intMinMax(p.x, minX)
		minY, _ = intMinMax(p.y, minY)
		minZ, _ = intMinMax(p.z, minZ)
		_, maxX = intMinMax(p.x, maxX)
		_, maxY = intMinMax(p.y, maxY)
		_, maxZ = intMinMax(p.z, maxZ)
	}

	lava := NewCube(minX, minY, minZ, maxX, maxY, maxZ)
	for _, p := range points {
		if err := lava.set(p, true); err != nil {
			panic(err)
		}
	}

	return lava, points
}

func main() {
	input := readInput()
	lava, points := parseInput(input)
	area := surfaceArea(lava, points)
	fmt.Println(area)
}
