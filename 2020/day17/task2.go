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

type hypercube [][][][]int // order: z, x, y

func (c hypercube) CountActive() int {
	var actives int
	for _, cube := range c {
		for _, plane := range cube {
			for _, row := range plane {
				for _, active := range row {
					actives += active
				}
			}
		}
	}

	return actives
}

func (c hypercube) CountActiveNeighbors(l, k, i, j int) int {
	var activeNeighbors int

	lMin, lMax, kMin, kMax, iMin, iMax, jMin, jMax := l-1, l+1, k-1, k+1, i-1, i+1, j-1, j+1

	if l == 0 {
		lMin = l
	} else if l == len(c)-1 {
		lMax = l
	}
	if k == 0 {
		kMin = k
	} else if k == len(c[0])-1 {
		kMax = k
	}
	if i == 0 {
		iMin = i
	} else if i == len(c[0][0])-1 {
		iMax = i
	}
	if j == 0 {
		jMin = j
	} else if j == len(c[0][0][0])-1 {
		jMax = j
	}

	activeNeighbors -= c[l][k][i][j]
	for l := lMin; l <= lMax; l++ {
		for k := kMin; k <= kMax; k++ {
			for i := iMin; i <= iMax; i++ {
				for j := jMin; j <= jMax; j++ {
					activeNeighbors += c[l][k][i][j]
				}
			}
		}
	}

	// fmt.Println(activeNeighbors)
	return activeNeighbors
}

func (c *hypercube) Cycle() int {
	nextHypercubes := make(hypercube, len(*c))

	for l, cube := range *c {
		nextHypercubes[l] = make([][][]int, len(cube))
		for k, plane := range cube {
			nextHypercubes[l][k] = make([][]int, len(plane))
			for i, row := range plane {
				nextHypercubes[l][k][i] = make([]int, len(row))
				for j, active := range row {
					nextHypercubes[l][k][i][j] = (*c)[l][k][i][j]
					activeNeighbors := c.CountActiveNeighbors(l, k, i, j)

					if active == 1 {
						if activeNeighbors != 2 && activeNeighbors != 3 {
							nextHypercubes[l][k][i][j] = 0
						}
					} else {
						if activeNeighbors == 3 {
							nextHypercubes[l][k][i][j] = 1
						}
					}
				}
			}
		}
	}

	*c = nextHypercubes

	return c.CountActive()
}

func (c *hypercube) String() string {
	str := "\n"

	for _, cube := range *c {
		for _, plane := range cube {
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
		str += "\n"
	}

	return str
}

func newHypercube(initState [][]int, expectedCycles int) *hypercube {
	xDim := len(initState) + (expectedCycles * 2)
	yDim := len(initState[0]) + (expectedCycles * 2)
	zDim := 1 + (expectedCycles * 2)
	wDim := 1 + (expectedCycles * 2)

	c := make(hypercube, wDim)
	for l := 0; l < wDim; l++ {
		c[l] = make([][][]int, zDim)
		for k := 0; k < zDim; k++ {
			c[l][k] = make([][]int, xDim)
			for i := 0; i < xDim; i++ {
				c[l][k][i] = make([]int, yDim)
			}
		}
	}

	for i, row := range initState {
		for j, active := range row {
			c[expectedCycles][expectedCycles][i+expectedCycles][j+expectedCycles] = active
		}
	}

	return &c
}

func initHypercube(input []string, expectedCycles int) *hypercube {
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

	return newHypercube(initState, expectedCycles)
}

func main() {
	cycles := 6

	input := readInput()

	hypercube := initHypercube(input, cycles)

	for i := 0; i < cycles; i++ {
		hypercube.Cycle()
	}

	fmt.Println(hypercube.CountActive())
}
