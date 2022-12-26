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

func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattan(x1, y1, x2, y2 int) int {
	x1, x2 = intMinMax(x1, x2)
	y1, y2 = intMinMax(y1, y2)
	return x2 - x1 + y2 - y1
}

type sensor struct {
	x, y, bx, by, detectRange int
}

type sensorGroup []*sensor

func (s *sensor) detectionRange() int {
	if s.detectRange == 0 {
		s.detectRange = manhattan(s.x, s.y, s.bx, s.by)
	}
	return s.detectRange
}

func (sg sensorGroup) inRange(x, y int) bool {
	for _, s := range sg {
		if manhattan(s.x, s.y, x, y) <= s.detectionRange() {
			return true
		}
	}
	return false
}

func (sg sensorGroup) findBeacon(minVal, maxVal int) (int, int) {
	// Check edge of sensors
	for _, s := range sg {
		sensorEdge := s.detectionRange() + 1
		for dy := -sensorEdge; dy <= sensorEdge; dy++ {
			y := s.y + dy
			if y < minVal {
				continue
			} else if y > maxVal {
				break
			}

			dx := sensorEdge - intAbs(dy)
			if x := s.x - dy; x >= 0 && !sg.inRange(x, y) {
				return x, y
			}
			if x := s.x + dx; x <= maxVal && !sg.inRange(x, y) {
				return x, y
			}
		}
	}

	return -1, -1
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

func parseInput(input []string) sensorGroup {
	sg := make(sensorGroup, 0, len(input))
	for _, line := range input {
		var s sensor
		n, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&s.x, &s.y, &s.bx, &s.by)
		if n != 4 {
			panic("Failed to parse all 4 coordinates")
		} else if err != nil {
			panic(err)
		}

		sg = append(sg, &s)
	}
	return sg
}

func main() {
	input := readInput()
	sg := parseInput(input)
	x, y := sg.findBeacon(0, 4000000)

	if x != -1 && y != -1 {
		fmt.Println(x*4000000 + y)
	} else {
		fmt.Println("Failed to find beacon")
	}
}
