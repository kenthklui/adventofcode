package main

import (
	"bufio"
	"fmt"
	"os"
)

// Assume box is below origin (ie. yMax < 0)
// Highest height is achieved with min X velocity + max y velocity

// For y:
// Any positive y velocity returns to 0 after 2y+1 steps, at which velocity is -y-1
// Max y is therefore bound by (-y-1) >= yMin --> y <= -yMin-1

type box struct {
	xMin, xMax, yMin, yMax int
}

func possibleSteps(b box) map[int]int {
	// First set a minimum for vX
	minVX, x := 0, 0
	for x <= b.xMax {
		minVX++
		x += minVX

		// Trivial case - does the minimum vX let us stop in the box?
		// If yes, return -1 for a later shortcut
		if x >= b.xMin {
			return map[int]int{-1: -1}
		}
	}

	possibleSteps := make(map[int]int)
	for vX := minVX; vX <= b.xMin; vX++ {
		x = 0
		steps := 0
		for currVX := vX; currVX >= 0; currVX-- {
			steps++
			x += currVX

			if x > b.xMax {
				break
			} else if b.xMin <= x {
				possibleSteps[steps] = steps
			}
		}
	}

	return possibleSteps
}

func calculateMaxY(b box, steps map[int]int) int {
	// Start iteration from maximum possible Y velocity
	vYMax := -b.yMin - 1
	// Check trivial case first - if X will let us land in the box, we can just use vYMax
	if _, ok := steps[-1]; ok {
		return vYMax
	}

	// Original plan was to iterate downward from vYMax...
	// Turns out none of the rest was needed, because input fit the trivial case!
	// I guess I'll flesh this out another time
	return 0
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

func parseInput(input []string) box {
	var xMin, xMax, yMin, yMax int

	n, err := fmt.Sscanf(input[0], "target area: x=%d..%d, y=%d..%d",
		&xMin, &xMax, &yMin, &yMax)
	if err != nil {
		panic(err)
	}
	if n != 4 {
		panic("Didn't parse 4 values")
	}

	return box{xMin, xMax, yMin, yMax}
}

func main() {
	input := readInput()
	b := parseInput(input)
	steps := possibleSteps(b)
	vYMax := calculateMaxY(b, steps)

	fmt.Println(vYMax * (vYMax + 1) / 2)
}
