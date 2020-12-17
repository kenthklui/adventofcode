package main

import (
	"bufio"
	"fmt"
	"os"
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

type cube [][][]int // order: z, x, y

func (c cube) CountActive() int {
	var actives int
	for _, plane := range c {
		for _, row := range plane {
			for _, active := range row {
				actives += active
			}
		}
	}

	return actives
}

func (c cube) CountActiveNeighbors(k, i, j int) int {
	var activeNeighbors int

	kMin, kMax, iMin, iMax, jMin, jMax := k-1, k+1, i-1, i+1, j-1, j+1

	if k == 0 {
		kMin = k
	} else if k == len(c)-1 {
		kMax = k
	}
	if i == 0 {
		iMin = i
	} else if i == len(c[0])-1 {
		iMax = i
	}
	if j == 0 {
		jMin = j
	} else if j == len(c[0][0])-1 {
		jMax = j
	}

	activeNeighbors -= c[k][i][j]
	for k := kMin; k <= kMax; k++ {
		for i := iMin; i <= iMax; i++ {
			for j := jMin; j <= jMax; j++ {
				activeNeighbors += c[k][i][j]
			}
		}
	}

	// fmt.Println(activeNeighbors)
	return activeNeighbors
}

func (c *cube) Cycle() int {
	nextCubes := make(cube, len(*c))

	for k, plane := range *c {
		nextCubes[k] = make([][]int, len(plane))
		for i, row := range plane {
			nextCubes[k][i] = make([]int, len(row))
			for j, active := range row {
				nextCubes[k][i][j] = (*c)[k][i][j]
				activeNeighbors := c.CountActiveNeighbors(k, i, j)

				if active == 1 {
					if activeNeighbors != 2 && activeNeighbors != 3 {
						nextCubes[k][i][j] = 0
					}
				} else {
					if activeNeighbors == 3 {
						nextCubes[k][i][j] = 1
					}
				}
			}
		}
	}

	*c = nextCubes

	return c.CountActive()
}

func (c *cube) String() string {
	str := "\n"

	for _, plane := range *c {
		for _, row := range plane {
			for _, active := range row {
				if active == 1 {
					str += "#"
				} else {
					str += "."
				}
			}
			str += "\n"
		}
		str += "\n"
	}

	return str
}

func newCube(initState [][]int, expectedCycles int) *cube {
	xDim := len(initState) + (expectedCycles * 4)
	yDim := len(initState[0]) + (expectedCycles * 4)
	zDim := 1 + (expectedCycles * 4)

	c := make(cube, zDim)
	for k := 0; k < len(c); k++ {
		c[k] = make([][]int, xDim)
		for i := 0; i < xDim; i++ {
			c[k][i] = make([]int, yDim)
		}
	}

	for i, row := range initState {
		for j, active := range row {
			c[expectedCycles][i+expectedCycles][j+expectedCycles] = active
		}
	}

	return &c
}

func initCube(input []string, expectedCycles int) *cube {
	initState := make([][]int, len(input))

	for i, line := range input {
		initState[i] = make([]int, len(line))
		for j, r := range line {
			if r == '#' {
				initState[i][j] = 1
			} else if r == '.' {
				initState[i][j] = 0
			} else {
				panic("Invalid state")
			}
		}
	}

	return newCube(initState, expectedCycles)
}

func main() {
	cycles := 6

	input := readInput()

	cube := initCube(input, cycles)

	for i := 0; i < cycles; i++ {
		cube.Cycle()
	}

	fmt.Println(cube.CountActive())
}
