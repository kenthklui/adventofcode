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

type sensor struct {
	x, y, bx, by int
}

func (s sensor) manhattan() int {
	x1, x2 := intMinMax(s.x, s.bx)
	y1, y2 := intMinMax(s.y, s.by)
	return x2 - x1 + y2 - y1
}

type beaconRange struct {
	start, end int
}

func (br1 beaconRange) greaterThan(br2 beaconRange) bool { return br1.start > br2.end+1 }
func (br1 beaconRange) lesserThan(br2 beaconRange) bool  { return br1.end+1 < br2.start }
func (br1 beaconRange) overlap(br2 beaconRange) bool {
	return br1.start <= br2.end+1 && br2.start <= br1.end+1
}

type beacons struct {
	ranges []beaconRange
}

func NewBeacons() *beacons {
	return &beacons{make([]beaconRange, 0)}
}

func (b *beacons) addRange(newRange beaconRange) {
	n := len(b.ranges)
	if n == 0 {
		b.ranges = append(b.ranges, newRange)
		return
	}

	if newRange.greaterThan(b.ranges[n-1]) { // larger than everything
		b.ranges = append(b.ranges, newRange)
		return
	} else if newRange.lesserThan(b.ranges[0]) { // smaller than everything
		b.ranges = append([]beaconRange{newRange}, b.ranges...)
		return
	}

	var overlapStartIndex, overlapEndIndex int
	for i, brs := range b.ranges {
		if newRange.overlap(brs) {
			overlapStartIndex, overlapEndIndex = i, i
			for _, bre := range b.ranges[i+1:] {
				if newRange.overlap(bre) {
					overlapEndIndex++
				} else {
					break
				}
			}
			break
		}
	}

	if firstStart := b.ranges[overlapStartIndex].start; newRange.start > firstStart {
		newRange.start = firstStart
	}
	if lastEnd := b.ranges[overlapEndIndex].end; newRange.end < lastEnd {
		newRange.end = lastEnd
	}

	newRanges := append(b.ranges[:overlapStartIndex], newRange)
	newRanges = append(newRanges, b.ranges[overlapEndIndex+1:]...)
	b.ranges = newRanges
}

func buildMap(sensors []sensor, minVal, maxVal int) []*beacons {
	bm := make([]*beacons, maxVal-minVal+1)
	for i := range bm {
		bm[i] = NewBeacons()
	}

	for _, s := range sensors {
		man := s.manhattan()
		for dy := -man; dy <= man; dy++ {
			y := s.y + dy
			if y < minVal {
				continue
			} else if y > maxVal {
				break
			}

			dx := man - intAbs(dy)
			bm[y-minVal].addRange(beaconRange{s.x - dx, s.x + dx})
		}
	}

	return bm
}

func countOccupied(sensors []sensor, row int) int {
	occupied := 0
	bm := buildMap(sensors, row, row)
	for _, br := range bm[0].ranges {
		occupied += br.end - br.start + 1
	}

	// Subtract tiles occupied by beacons
	beaconX := make(map[int]byte)
	for _, s := range sensors {
		if s.by == row {
			beaconX[s.bx] = 0
		}
	}
	occupied -= len(beaconX)

	return occupied
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
	occupied := countOccupied(sensors, 2000000)

	fmt.Println(occupied)
}
