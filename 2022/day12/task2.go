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

type hillmap struct {
	elev                               [][]int
	startRow, startCol, endRow, endCol int
}

type point struct {
	row, col, steps int
}

func (hm *hillmap) traverse(startRow, startCol int) int {
	// status: 0 = untouched, 1 = queued, 2 = status
	status := make([][]uint8, len(hm.elev))
	for i, row := range hm.elev {
		status[i] = make([]uint8, len(row))
	}
	status[startRow][startCol] = 2

	queue := make([]point, 0, len(hm.elev))
	queue = append(queue, point{startRow, startCol, 0})

	for index := 0; index < len(queue); index++ {
		curr := queue[index]
		status[curr.row][curr.col] = 2
		if curr.row == hm.endRow && curr.col == hm.endCol {
			return curr.steps
		}

		for _, next := range hm.options(curr, status) {
			status[next.row][next.col] = 1
			queue = append(queue, next)
		}
	}

	return -1
}

func (hm *hillmap) fastestPath() int {
	maxPath := len(hm.elev) * len(hm.elev[0])
	for row, line := range hm.elev {
		for col, e := range line {
			if e == 0 {
				path := hm.traverse(row, col)
				if path != -1 && path < maxPath {
					maxPath = path
				}
			}
		}
	}

	return maxPath
}

func (hm *hillmap) options(p point, status [][]uint8) []point {
	opts := make([]point, 0, 4)
	if p.row != 0 {
		if status[p.row-1][p.col] == 0 {
			diff := hm.elev[p.row][p.col] - hm.elev[p.row-1][p.col]
			if diff >= -1 {
				opts = append(opts, point{p.row - 1, p.col, p.steps + 1})
			}
		}
	}
	if p.row < len(hm.elev)-1 {
		if status[p.row+1][p.col] == 0 {
			diff := hm.elev[p.row][p.col] - hm.elev[p.row+1][p.col]
			if diff >= -1 {
				opts = append(opts, point{p.row + 1, p.col, p.steps + 1})
			}
		}
	}
	if p.col != 0 {
		if status[p.row][p.col-1] == 0 {
			diff := hm.elev[p.row][p.col] - hm.elev[p.row][p.col-1]
			if diff >= -1 {
				opts = append(opts, point{p.row, p.col - 1, p.steps + 1})
			}
		}
	}
	if p.col < len(hm.elev[p.row])-1 {
		if status[p.row][p.col+1] == 0 {
			diff := hm.elev[p.row][p.col] - hm.elev[p.row][p.col+1]
			if diff >= -1 {
				opts = append(opts, point{p.row, p.col + 1, p.steps + 1})
			}
		}
	}
	return opts
}

func NewHillmap(rowCount, colCount int) *hillmap {
	hm := new(hillmap)
	hm.elev = make([][]int, rowCount)
	for i := range hm.elev {
		hm.elev[i] = make([]int, colCount)
	}
	return hm
}

func parseInput(input []string) *hillmap {
	hm := NewHillmap(len(input), len(input[0]))
	for row, line := range input {
		for col, r := range line {
			switch r {
			case 'S':
				hm.elev[row][col] = 0
				hm.startRow, hm.startCol = row, col
			case 'E':
				hm.elev[row][col] = 25
				hm.endRow, hm.endCol = row, col
			default:
				hm.elev[row][col] = int(r - 'a')
			}
		}
	}

	return hm
}

func main() {
	input := readInput()
	hm := parseInput(input)
	fmt.Println(hm.fastestPath())
}
