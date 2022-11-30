package main

import (
	"bufio"
	"fmt"
	"os"
)

type cuboid struct {
	xMin, xMax, yMin, yMax, zMin, zMax int
	state                              byte
}

const minBounds = -50
const maxBounds = 50

func (c cuboid) inBounds() bool {
	outOfBounds := c.xMin < minBounds || c.yMin < minBounds || c.zMin < minBounds ||
		c.xMax > maxBounds || c.yMax > maxBounds || c.zMax > maxBounds

	return !outOfBounds
}

type reactor struct {
	cube [][][]byte
}

// 100x100x100 is fine for now, but task 2 will no doubt be infinite sized. Worry about it later!
func NewReactor() *reactor {
	cube := make([][][]byte, 101)
	for i := 0; i < 101; i++ {
		cube[i] = make([][]byte, 101)
		for j := 0; j < 101; j++ {
			cube[i][j] = make([]byte, 101)
		}
	}

	return &reactor{cube: cube}
}

func (r *reactor) applyCuboid(c cuboid) {
	for x := c.xMin; x <= c.xMax; x++ {
		for y := c.yMin; y <= c.yMax; y++ {
			for z := c.zMin; z <= c.zMax; z++ {
				r.cube[x+50][y+50][z+50] = c.state
			}
		}
	}
}

func (r *reactor) countOn() int {
	onCount := 0
	for x := 0; x < len(r.cube); x++ {
		for y := 0; y < len(r.cube[0]); y++ {
			for z := 0; z < len(r.cube[0][0]); z++ {
				if r.cube[x][y][z] == 1 {
					onCount++
				}
			}
		}
	}

	return onCount
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

func parseInput(input []string) []cuboid {
	var stateStr string
	var state byte
	var xMin, xMax, yMin, yMax, zMin, zMax int

	cuboids := make([]cuboid, 0, len(input))
	for _, line := range input {
		n, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
			&stateStr, &xMin, &xMax, &yMin, &yMax, &zMin, &zMax)
		if err != nil {
			fmt.Println(line)
			panic(err)
		} else if n != 7 {
			panic("Failed to parse cuboid box coordinates")
		}

		if stateStr == "on" {
			state = 1
		} else {
			state = 0
		}

		c := cuboid{xMin, xMax, yMin, yMax, zMin, zMax, state}
		if c.inBounds() {
			cuboids = append(cuboids, c)
		}
	}

	return cuboids
}

func main() {
	input := readInput()
	cuboids := parseInput(input)
	r := NewReactor()

	for _, c := range cuboids {
		r.applyCuboid(c)
	}

	fmt.Println(r.countOn())
}
