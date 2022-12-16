package main

import (
	"bufio"
	"fmt"
	"os"
)

func intMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

type sensor struct {
	x, y, bx, by int
}

func (s sensor) manhattan() int {
	x1, x2 := intMinMax(s.x, s.bx)
	y1, y2 := intMinMax(s.y, s.by)
	return x2 - x1 + y2 - y1
}

func countImpossible(sensors []sensor, row int) int {
	bm := make(map[int]bool)
	for _, s := range sensors {
		y1, y2 := intMinMax(s.y, row)
		extend := s.manhattan() - (y2 - y1)
		if extend >= 0 {
			for x := s.x - extend; x <= s.x+extend; x++ {
				bm[x] = false
			}
		}
	}

	for _, s := range sensors {
		if s.by == row {
			if _, ok := bm[s.bx]; ok {
				delete(bm, s.bx)
			}
		}
	}

	return len(bm)
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

func parseInput(input []string) []sensor {
	sensors := make([]sensor, 0, len(input))
	for _, line := range input {
		var s sensor
		n, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&s.x, &s.y, &s.bx, &s.by)
		if n != 4 {
			panic("Failed to parse all 4 coordinates")
		} else if err != nil {
			panic(err)
		}

		sensors = append(sensors, s)
	}
	return sensors
}

func main() {
	input := readInput()
	sensors := parseInput(input)
	fmt.Println(countImpossible(sensors, 2000000))
}
