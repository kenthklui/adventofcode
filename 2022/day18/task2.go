package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

type void struct{}
type grid map[int]plane
type plane map[int]line
type line map[int]void

func (g grid) get(x, y, z int) bool {
	if p, ok1 := g[x]; ok1 {
		if l, ok2 := p[y]; ok2 {
			if _, ok3 := l[z]; ok3 {
				return true
			}
		}
	}
	return false
}

var empty void

func parseInput(input []string) grid {
	g := make(grid)

	for _, text := range input {
		splits := strings.Split(text, ",")
		x, _ := strconv.Atoi(splits[0])
		y, _ := strconv.Atoi(splits[1])
		z, _ := strconv.Atoi(splits[2])
		if _, ok := g[x]; !ok {
			g[x] = make(plane)
		}
		if _, ok := g[x][y]; !ok {
			g[x][y] = make(line)
		}
		g[x][y][z] = empty
	}

	return g
}

type cube struct {
	val                    [][][]bool
	startX, startY, startZ int
}

func NewCube(startX, startY, startZ, sizeX, sizeY, sizeZ int) cube {
	var c cube
	c.val = make([][][]bool, sizeX)
	for x := range c.val {
		c.val[x] = make([][]bool, sizeY)
		for y := range c.val[x] {
			c.val[x][y] = make([]bool, sizeZ)
		}
	}
	c.startX, c.startY, c.startZ = startX, startY, startZ
	return c
}
func (c cube) get(x, y, z int) bool {
	cx, cy, cz := x-c.startX, y-c.startY, z-c.startZ
	if cx < 0 || cx >= len(c.val) ||
		cy < 0 || cy >= len(c.val[cx]) ||
		cz < 0 || cz >= len(c.val[cx][cy]) {
		return true
	}
	return c.val[cx][cy][cz]
}
func (c cube) set(x, y, z int, val bool) {
	cx, cy, cz := x-c.startX, y-c.startY, z-c.startZ
	c.val[cx][cy][cz] = val
}

func externalSurfaceArea(g grid) int {
	minX, minY, minZ := 1000, 1000, 1000
	maxX, maxY, maxZ := -1000, -1000, -1000

	for x, p := range g {
		if x > maxX {
			maxX = x
		}
		if x < minX {
			minX = x
		}
		for y, l := range p {
			if y > maxY {
				maxY = y
			}
			if y < minY {
				minY = y
			}
			for z := range l {
				if z > maxZ {
					maxZ = z
				}
				if z < minZ {
					minZ = z
				}
			}
		}
	}

	// Give the cube extra space so we can manuever around the edges
	sizeX := maxX - minX + 3
	sizeY := maxY - minY + 3
	sizeZ := maxZ - minZ + 3
	startX := minX - 1
	startY := minY - 1
	startZ := minZ - 1

	c := NewCube(startX, startY, startZ, sizeX, sizeY, sizeZ)
	return recurseCountArea(c, g, startX, startY, startZ)
}

func recurseCountArea(c cube, g grid, x, y, z int) int {
	c.set(x, y, z, true)

	sum := 0
	if !c.get(x-1, y, z) {
		if g.get(x-1, y, z) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x-1, y, z)
		}
	}
	if !c.get(x+1, y, z) {
		if g.get(x+1, y, z) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x+1, y, z)
		}
	}
	if !c.get(x, y-1, z) {
		if g.get(x, y-1, z) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x, y-1, z)
		}
	}
	if !c.get(x, y+1, z) {
		if g.get(x, y+1, z) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x, y+1, z)
		}
	}
	if !c.get(x, y, z-1) {
		if g.get(x, y, z-1) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x, y, z-1)
		}
	}
	if !c.get(x, y, z+1) {
		if g.get(x, y, z+1) {
			sum++
		} else {
			sum += recurseCountArea(c, g, x, y, z+1)
		}
	}

	return sum
}

func main() {
	input := readInput()
	object := parseInput(input)
	fmt.Println(externalSurfaceArea(object))
}
