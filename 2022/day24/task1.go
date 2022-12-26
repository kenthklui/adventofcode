package main

import (
	"bufio"
	"fmt"
	"os"
)

type mountain struct {
	blizzards     [][]byte
	height, width int
}

func NewMountain(height, width int) *mountain {
	blizzards := make([][]byte, height)
	for i := range blizzards {
		blizzards[i] = make([]byte, width)
	}

	return &mountain{
		blizzards: blizzards,
		height:    height,
		width:     width,
	}
}

func (m *mountain) clear(x, y, minute int) bool {
	rowOffset, colOffset := minute%m.height, minute%m.width
	if up := (y + rowOffset) % m.height; m.blizzards[up][x] == '^' {
		return false
	}
	if down := (y - rowOffset + m.height) % m.height; m.blizzards[down][x] == 'v' {
		return false
	}
	if left := (x + colOffset) % m.width; m.blizzards[y][left] == '<' {
		return false
	}
	if right := (x - colOffset + m.width) % m.width; m.blizzards[y][right] == '>' {
		return false
	}
	return true
}

func (m *mountain) next(x, y int) [][2]int {
	neighbors := make([][2]int, 0, 5)
	neighbors = append(neighbors, [2]int{x, y})
	if x > 0 {
		neighbors = append(neighbors, [2]int{x - 1, y})
	}
	if x < m.width-1 {
		neighbors = append(neighbors, [2]int{x + 1, y})
	}
	if y > 0 {
		neighbors = append(neighbors, [2]int{x, y - 1})
	}
	if y < m.height-1 {
		neighbors = append(neighbors, [2]int{x, y + 1})
	}
	return neighbors
}

func (m *mountain) traverse() int {
	prevReached, reached := make([][2]int, 0, m.height*m.width), make([][2]int, 0, m.height*m.width)
	checked := make([][]bool, m.height)
	for i := range checked {
		checked[i] = make([]bool, m.width)
	}

	for min := 1; true; min++ {
		checked[0][0] = true
		if m.clear(0, 0, min) {
			reached = append(reached, [2]int{0, 0})
		}

		for _, r := range prevReached {
			if r[0] == m.width-1 && r[1] == m.height-1 {
				return min
			}

			for _, n := range m.next(r[0], r[1]) {
				if !checked[n[1]][n[0]] {
					checked[n[1]][n[0]] = true
					if m.clear(n[0], n[1], min) {
						reached = append(reached, n)
					}
				}

			}
		}

		prevReached, reached = reached, prevReached[:0]
		for y, row := range checked {
			for x := range row {
				checked[y][x] = false
			}
		}
	}

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

func parseInput(input []string) *mountain {
	height := len(input) - 2
	width := len(input[0]) - 2
	m := NewMountain(height, width)
	for y, line := range input[1 : len(input)-1] {
		for x, r := range line[1 : len(line)-1] {
			m.blizzards[y][x] = byte(r)
		}
	}

	return m
}

func main() {
	input := readInput()
	mount := parseInput(input)
	done := mount.traverse()
	fmt.Println(done)
}
