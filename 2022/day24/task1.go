package main

import (
	"bufio"
	"fmt"
	"os"
)

type position struct {
	x, y int
}

type blizzardLine []int8

func NewBlizzardLine(size uint) blizzardLine {
	return make([]int8, size)
}

func (b blizzardLine) clear(coordinate, minute uint) bool {
	minute %= uint(len(b))
	if forward := (coordinate + minute) % uint(len(b)); b[forward] == -1 {
		return false
	}
	if backward := (coordinate + uint(len(b)) - minute) % uint(len(b)); b[backward] == 1 {
		return false
	}
	return true
}

type flag struct {
	x, y, minute uint
}

func (f flag) next(height, width uint) []flag {
	n := make([]flag, 0, 5)
	if f.x > 0 {
		n = append(n, flag{f.x - 1, f.y, f.minute + 1})
	}
	if f.x < width-1 {
		n = append(n, flag{f.x + 1, f.y, f.minute + 1})
	}
	if f.y > 0 {
		n = append(n, flag{f.x, f.y - 1, f.minute + 1})
	}
	if f.y < height-1 {
		n = append(n, flag{f.x, f.y + 1, f.minute + 1})
	}
	n = append(n, flag{f.x, f.y, f.minute + 1})
	return n
}

type mountain struct {
	rowBlizz, colBlizz []blizzardLine
	height, width      uint
}

func NewMountain(height, width uint) *mountain {
	m := mountain{
		rowBlizz: make([]blizzardLine, height),
		colBlizz: make([]blizzardLine, width),
		height:   height,
		width:    width,
	}

	for i := range m.rowBlizz {
		m.rowBlizz[i] = NewBlizzardLine(width)
	}

	for i := range m.colBlizz {
		m.colBlizz[i] = NewBlizzardLine(height)
	}

	return &m
}

func (m *mountain) addBlizzard(x, y int, r rune) {
	switch r {
	case '>':
		m.rowBlizz[y][x] = 1
	case 'v':
		m.colBlizz[x][y] = 1
	case '<':
		m.rowBlizz[y][x] = -1
	case '^':
		m.colBlizz[x][y] = -1
	case '.':
		// Not a blizzard
	default:
		panic("Invalid rune")
	}
}

func (m *mountain) clear(f flag) bool {
	return m.rowBlizz[f.y].clear(f.x, f.minute) && m.colBlizz[f.x].clear(f.y, f.minute)
}

type void struct{}

var empty void

func (m *mountain) traverse() uint {
	queue := make([]flag, 0, 1<<12)
	// Assume we enter the storm on first turn - maybe not actually safe to do?
	queue = append(queue, flag{0, 0, 1})

	queued := make([]bool, m.height*m.width)
	var currMin uint = 1
	for len(queue) > 0 {
		currFlag := queue[0]
		queue = queue[1:]

		if currFlag.x == m.width-1 && currFlag.y == m.height-1 {
			return currFlag.minute + 1
		}

		if currFlag.minute != currMin {
			for i := range queued {
				queued[i] = false
			}
			currMin = currFlag.minute
		}

		for _, nextFlag := range currFlag.next(m.height, m.width) {
			if m.clear(nextFlag) {
				key := m.width*nextFlag.y + nextFlag.x
				if !queued[key] {
					queue = append(queue, nextFlag)
					queued[key] = true
				}

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
	height := uint(len(input) - 2)
	width := uint(len(input[0]) - 2)
	m := NewMountain(height, width)
	for y, line := range input[1 : len(input)-1] {
		for x, r := range line[1 : len(line)-1] {
			m.addBlizzard(x, y, r)
		}
	}

	return m
}

func main() {
	input := readInput()
	mount := parseInput(input)
	fmt.Println(mount.traverse())
}
