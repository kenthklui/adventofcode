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

type line struct {
	X1, Y1, X2, Y2 int
}

type floorMap struct {
	dimX, dimY int
	points     []int
}

func NewFloorMap(xMax, yMax int) *floorMap {
	dimX := xMax + 1
	dimY := yMax + 1
	size := dimX * dimY

	return &floorMap{dimX, dimY, make([]int, size)}
}

func (fm *floorMap) PlotLine(l line) {
	var xInc, yInc int

	if l.Y1 == l.Y2 {
		yInc = 0
	} else if l.Y1 < l.Y2 {
		yInc = 1
	} else {
		yInc = -1
	}

	if l.X1 == l.X2 {
		xInc = 0
	} else if l.X1 < l.X2 {
		xInc = 1
	} else {
		xInc = -1
	}

	for x, y := l.X1, l.Y1; ; {
		pos := y*fm.dimX + x

		fm.points[pos]++

		if x == l.X2 && y == l.Y2 {
			break
		}

		x += xInc
		y += yInc
	}
}

func (fm *floorMap) HazardPoints() int {
	points := 0
	for _, p := range fm.points {
		if p > 1 {
			points++
		}
	}

	return points
}

func (fm *floorMap) Print() {
	fmt.Println(fm.points)
}

func parseInput(input []string) ([]line, int, int) {
	var x1, y1, x2, y2 int
	var xMax, yMax int
	lines := make([]line, len(input))

	for i, s := range input {
		n, err := fmt.Sscanf(s, "%d,%d -> %d,%d", &x1, &y1, &x2, &y2)
		if err != nil {
			fmt.Println(s)
			panic(err)
		} else if n != 4 {
			panic("Failed to parse 4 coordinates")
		}

		lines[i] = line{x1, y1, x2, y2}

		if x1 > xMax {
			xMax = x1
		}
		if x2 > xMax {
			xMax = x2
		}
		if y1 > yMax {
			yMax = y1
		}
		if y2 > yMax {
			yMax = y2
		}
	}

	return lines, xMax, yMax
}

func main() {
	input := readInput()
	lines, xMax, yMax := parseInput(input)

	fm := NewFloorMap(xMax, yMax)
	for _, l := range lines {
		fm.PlotLine(l)
	}

	// fm.Print()

	fmt.Println(fm.HazardPoints())
}
