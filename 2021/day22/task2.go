package main

import (
	"bufio"
	"fmt"
	"os"
)

func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

type point struct {
	x, y, z int
}

type cuboid struct {
	xMin, xMax, yMin, yMax, zMin, zMax int
}

func (c *cuboid) String() string {
	return fmt.Sprintf("{(%d,%d,%d),(%d,%d,%d)}", c.xMin, c.yMin, c.zMin, c.xMax, c.yMax, c.zMax)
}

func NewCuboid(xMin, xMax, yMin, yMax, zMin, zMax int) *cuboid {
	if xMin > xMax || yMin > yMax || zMin > zMax {
		return nil
	}

	return &cuboid{xMin, xMax, yMin, yMax, zMin, zMax}
}

func (c *cuboid) contains(p *point) bool {
	return p.x >= c.xMin && p.y >= c.yMin && p.z >= c.zMin &&
		p.x <= c.xMax && p.y <= c.yMax && p.z <= c.zMax
}

func (c *cuboid) volume() int {
	return (c.xMax - c.xMin + 1) * (c.yMax - c.yMin + 1) * (c.zMax - c.zMin + 1)
}

func (c1 *cuboid) intersect(c2 *cuboid) *cuboid {
	if c1.xMin > c2.xMax || c1.xMax < c2.xMin ||
		c1.yMin > c2.yMax || c1.yMax < c2.yMin ||
		c1.zMin > c2.zMax || c1.zMax < c2.zMin {
		return nil
	}

	xMin, yMin, zMin := intMax(c1.xMin, c2.xMin), intMax(c1.yMin, c2.yMin), intMax(c1.zMin, c2.zMin)
	xMax, yMax, zMax := intMin(c1.xMax, c2.xMax), intMin(c1.yMax, c2.yMax), intMin(c1.zMax, c2.zMax)
	c := NewCuboid(xMin, xMax, yMin, yMax, zMin, zMax)

	return c
}

func (c1 *cuboid) subtract(c2 *cuboid) []*cuboid {
	intersection := c1.intersect(c2)
	if intersection == nil {
		return []*cuboid{c1}
	}

	remainder := make([]*cuboid, 0)
	xMins := []int{c1.xMin, intersection.xMin, intersection.xMax + 1}
	yMins := []int{c1.yMin, intersection.yMin, intersection.yMax + 1}
	zMins := []int{c1.zMin, intersection.zMin, intersection.zMax + 1}
	xMaxs := []int{intersection.xMin - 1, intersection.xMax, c1.xMax}
	yMaxs := []int{intersection.yMin - 1, intersection.yMax, c1.yMax}
	zMaxs := []int{intersection.zMin - 1, intersection.zMax, c1.zMax}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				// Skip intersection cuboid
				if i == 1 && j == 1 && k == 1 {
					continue
				}

				c := NewCuboid(xMins[i], xMaxs[i], yMins[j], yMaxs[j], zMins[k], zMaxs[k])
				if c != nil {
					remainder = append(remainder, c)
				}
			}
		}
	}

	return remainder
}

type instruction struct {
	cube  *cuboid
	state byte
}

type reactor struct {
	cubesOn []*cuboid
}

func NewReactor() *reactor {
	return &reactor{cubesOn: make([]*cuboid, 0)}
}

func (r *reactor) applyInstruction(i instruction) {
	if i.state == 1 {
		cubeAddition := []*cuboid{i.cube}

		for _, cubeOn := range r.cubesOn {
			newAddition := make([]*cuboid, 0)
			for _, cubeAdd := range cubeAddition {
				newAddition = append(newAddition, cubeAdd.subtract(cubeOn)...)
			}

			cubeAddition = newAddition
		}

		r.cubesOn = append(r.cubesOn, cubeAddition...)
	} else {
		newCubesOn := make([]*cuboid, 0)
		for _, cubeOn := range r.cubesOn {
			newCubesOn = append(newCubesOn, cubeOn.subtract(i.cube)...)
		}
		r.cubesOn = newCubesOn
	}
}

func (r *reactor) countOn() int {
	onCount := 0
	for _, c := range r.cubesOn {
		onCount += c.volume()
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

func parseInput(input []string) []instruction {
	var stateStr string
	var state byte
	var xMin, xMax, yMin, yMax, zMin, zMax int

	instructions := make([]instruction, 0, len(input))
	for _, line := range input {
		n, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
			&stateStr, &xMin, &xMax, &yMin, &yMax, &zMin, &zMax)
		if err != nil {
			fmt.Println(line)
			panic(err)
		} else if n != 7 {
			panic("Failed to parse instruction values")
		}

		if stateStr == "on" {
			state = 1
		} else {
			state = 0
		}

		c := NewCuboid(xMin, xMax, yMin, yMax, zMin, zMax)
		if c != nil {
			i := instruction{cube: c, state: state}
			instructions = append(instructions, i)
		}
	}

	return instructions
}

func main() {
	input := readInput()
	instructions := parseInput(input)
	r := NewReactor()

	for _, instruction := range instructions {
		r.applyInstruction(instruction)
	}

	fmt.Println(r.countOn())
}
