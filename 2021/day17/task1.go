package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Assume box is below origin (ie. yMax < 0)
// Highest height is achieved with min X velocity + max y velocity

type box struct {
	xMin, xMax, yMin, yMax int
}

// Assume positive x box
func possibleVXs(b box, yStepMap map[int][]int) map[int][]int {
	stepMap := make(map[int][]int)

	// min X velocity is achieved when vX(vX+1)/2 == xMin
	minVX := int(math.Ceil(math.Sqrt(float64(b.xMin)*2.0+0.25) - 0.5))
	// max X velocity is achieved when you hit xMax in 1 step
	maxVX := b.xMax

	maxSteps := 0
	for step := range yStepMap {
		if step > maxSteps {
			maxSteps = step
		}
	}

	for vX := maxVX; vX >= minVX; vX-- {
		x := 0
		currVX := vX

		for steps := 1; steps <= maxSteps; steps++ {
			x += currVX

			if x > b.xMax {
				break
			}

			if x >= b.xMin {
				if _, ok := yStepMap[steps]; !ok {
					continue
				}

				if stepList, ok := stepMap[steps]; ok {
					stepMap[steps] = append(stepList, vX)
				} else {
					stepMap[steps] = []int{vX}
				}
			}

			if currVX > 0 {
				currVX--
			}
		}
	}

	return stepMap
}

// Assume negative y box
func possibleVYs(b box) map[int][]int {
	stepMap := make(map[int][]int)

	// First constrain by y box
	minVY := b.yMin
	maxVY := -b.yMin - 1

	for vY := maxVY; vY >= minVY; vY-- {
		y := 0
		steps := 0

		for currVY := vY; currVY >= minVY; currVY-- {
			steps++
			y += currVY

			if y < b.yMin {
				break
			}

			if y <= b.yMax {
				if stepList, ok := stepMap[steps]; ok {
					stepMap[steps] = append(stepList, vY)
				} else {
					stepMap[steps] = []int{vY}
				}
			}
		}
	}

	return stepMap
}

func findMaxHeight(xStepMap, yStepMap map[int][]int) int {
	maxVY := 0

	for step := range xStepMap {
		if yVelocities, ok := yStepMap[step]; ok {
			if yVelocities[0] > maxVY {
				maxVY = yVelocities[0]
			}
		}
	}

	return maxVY * (maxVY + 1) / 2
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

	yStepMap := possibleVYs(b)
	xStepMap := possibleVXs(b, yStepMap)

	fmt.Println(findMaxHeight(xStepMap, yStepMap))
}
