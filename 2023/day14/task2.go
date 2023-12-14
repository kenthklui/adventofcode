package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type rockMap struct {
	rm          [][]rune
	sig         string
	gridInvalid bool
}

func makeMap(input []string) *rockMap {
	grid := make([][]rune, len(input))
	for i, line := range input {
		grid[i] = make([]rune, len(line))
		for j, c := range line {
			grid[i][j] = c
		}
	}

	rm := &rockMap{
		rm:          grid,
		gridInvalid: false,
	}
	rm.updateSig()

	return rm
}

func (rm *rockMap) updateSig() {
	if rm.gridInvalid {
		panic("Updating signature when grid is invalid")
	}

	var b strings.Builder
	for _, line := range rm.rm {
		b.WriteString(string(line))
	}
	rm.sig = b.String()
}

func (rm *rockMap) restoreFromSig() {
	if rm.gridInvalid {
		for i, line := range rm.rm {
			for j := range line {
				rm.rm[i][j] = rune(rm.sig[i*len(line)+j])
			}
		}
		rm.gridInvalid = false
	}
}

func (rm *rockMap) roll() {
	rm.restoreFromSig()
	for col := range rm.rm[0] {
		openSpot := -1
		for row := range rm.rm {
			switch rm.rm[row][col] {
			case '.':
				if openSpot < 0 {
					openSpot = row
				}
			case 'O':
				if openSpot >= 0 {
					rm.rm[openSpot][col], rm.rm[row][col] = 'O', '.'
					openSpot++
				}
			case '#':
				openSpot = -1
			default:
				panic("Invalid map char")
			}
		}
	}
	rm.updateSig()
}

func (rm *rockMap) rotate() {
	rm.restoreFromSig()
	mid, max := len(rm.rm)/2, len(rm.rm)-1
	for x := 0; x < mid; x++ {
		for y := 0; y < mid; y++ {
			rm.rm[x][y], rm.rm[y][max-x], rm.rm[max-x][max-y], rm.rm[max-y][x] =
				rm.rm[max-y][x], rm.rm[x][y], rm.rm[y][max-x], rm.rm[max-x][max-y]
		}
	}
	rm.updateSig()
}

func (rm *rockMap) cycle() {
	for i := 0; i < 4; i++ {
		rm.roll()
		rm.rotate()
	}
}

func (rm *rockMap) runCycles(runs int) {
	cache := make(map[string]string)

	for cycles := 0; cycles < runs; cycles++ {
		oldSig := rm.sig
		if cached, found := cache[oldSig]; found {
			cycleLength := 1
			for s := cached; s != oldSig; s = cache[s] {
				cycleLength++
			}
			remainingRuns := (runs - cycles) % cycleLength
			for i := 0; i < remainingRuns; i++ {
				rm.sig = cache[rm.sig]
			}
			rm.gridInvalid = true
			break
		} else {
			rm.cycle()
			cache[oldSig] = rm.sig
		}
	}
}

func (rm *rockMap) load() int {
	rm.restoreFromSig()
	sum := 0
	for rowNum, row := range rm.rm {
		for _, c := range row {
			if c == 'O' {
				sum += len(rm.rm) - rowNum
			}
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	rm := makeMap(input)
	rm.runCycles(1000000000)
	fmt.Println(rm.load())
}
