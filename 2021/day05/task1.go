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
	if l.X1 != l.X2 && l.Y1 != l.Y2 {
		return // skip non straight lines
	}

	var xStart, xStop, yStart, yStop int
	if l.X1 == l.X2 {
		xStart = l.X1
		xStop = l.X1

		if l.Y1 < l.Y2 {
			yStart, yStop = l.Y1, l.Y2
		} else {
			yStart, yStop = l.Y2, l.Y1
		}

	} else if l.Y1 == l.Y2 {
		yStart = l.Y1
		yStop = l.Y1

		if l.X1 < l.X2 {
			xStart, xStop = l.X1, l.X2
		} else {
			xStart, xStop = l.X2, l.X1
		}
	}

	for x := xStart; x <= xStop; x++ {
		for y := yStart; y <= yStop; y++ {
			pos := y*fm.dimX + x

			fm.points[pos]++
		}
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

	fmt.Println(fm.HazardPoints())
}
