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
		startRow := 0
		rockCount := 0
		for row := range rm.rm {
			switch rm.rm[row][col] {
			case '.':
				continue
			case 'O':
				rockCount++
				continue
			case '#':
				endRow := row
				for i := startRow; i < endRow; i++ {
					if rockCount > 0 {
						rm.rm[i][col] = 'O'
						rockCount--
					} else {
						rm.rm[i][col] = '.'
					}
				}

				startRow = row + 1
			default:
				panic("Invalid map char")
			}
		}

		if rockCount > 0 {
			for i := startRow; i < len(rm.rm); i++ {
				if rockCount > 0 {
					rm.rm[i][col] = 'O'
					rockCount--
				} else {
					rm.rm[i][col] = '.'
				}
			}
		}
	}
	rm.updateSig()
}

func (rm *rockMap) rotate() {
	rm.restoreFromSig()
	width := len(rm.rm[0])
	rotated := make([][]rune, width)
	for row := range rotated {
		rotated[row] = make([]rune, len(rm.rm))
		for col := range rotated[row] {
			rotated[row][col] = rm.rm[width-1-col][row]
		}
	}
	rm.rm = rotated
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
