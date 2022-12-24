package main

import (
	"bufio"
	"fmt"
	"os"
)

type position struct {
	x, y int
}

func intAbsDiff(a, b uint) uint {
	if a < b {
		return b - a
	} else {
		return a - b
	}
}

type blizzardLine []int8

func NewBlizzardLine(size uint) blizzardLine {
	return make([]int8, size)
}

func (b blizzardLine) clear(coordinate, minute uint) bool {
	length := uint(len(b))
	minute %= length
	if forward := (coordinate + minute) % length; b[forward] == -1 {
		return false
	}
	if backward := (coordinate + length - minute) % length; b[backward] == 1 {
		return false
	}
	return true
}

type flag struct {
	x, y, minute uint
}

func (f flag) next(height, width uint) []flag {
	if f.y == 0 {
		return []flag{flag{f.x, f.y + 1, f.minute + 1}, flag{f.x, f.y, f.minute + 1}}
	} else if f.y == height-1 {
		return []flag{flag{f.x, f.y - 1, f.minute + 1}, flag{f.x, f.y, f.minute + 1}}
	}

	n := make([]flag, 0, 5)
	if f.x > 1 {
		n = append(n, flag{f.x - 1, f.y, f.minute + 1})
	}
	if f.x < width-2 {
		n = append(n, flag{f.x + 1, f.y, f.minute + 1})
	}
	if f.y > 1 {
		n = append(n, flag{f.x, f.y - 1, f.minute + 1})
	}
	if f.y < height-2 {
		n = append(n, flag{f.x, f.y + 1, f.minute + 1})
	}
	n = append(n, flag{f.x, f.y, f.minute + 1})
	return n
}

func (f flag) key(height, width uint) uint {
	k := f.minute
	k = k*height + f.y
	k = k*width + f.x
	return k
}

type mountain struct {
	rowBlizz, colBlizz []blizzardLine
	height, width      uint
}

func NewMountain(height, width int) *mountain {
	m := mountain{
		rowBlizz: make([]blizzardLine, height-2),
		colBlizz: make([]blizzardLine, width-2),
		height:   uint(height),
		width:    uint(width),
	}

	for i := range m.rowBlizz {
		m.rowBlizz[i] = NewBlizzardLine(m.width - 2)
	}
	for i := range m.colBlizz {
		m.colBlizz[i] = NewBlizzardLine(m.height - 2)
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
	case '#':
		// Not a blizzard
	default:
		panic("Invalid rune")
	}
}

func (m *mountain) clear(f flag) bool {
	if f.y == 0 || f.y == m.height-1 {
		return true // start and end are always clear
	}
	return m.rowBlizz[f.y-1].clear(f.x-1, f.minute) && m.colBlizz[f.x-1].clear(f.y-1, f.minute)
}

type void struct{}

var empty void

func (m *mountain) traverse(startTime uint, forward bool) uint {
	queued := make(map[uint]void)
	queue := make([]flag, 0, 1<<20)

	var startX, startY, endX, endY uint
	if forward {
		startX, startY, endX, endY = 1, 0, m.width-2, m.height-1
	} else {
		startX, startY, endX, endY = m.width-2, m.height-1, 1, 0
	}
	queue = append(queue, flag{startX, startY, startTime})

	for len(queue) > 0 {
		currFlag := queue[0]
		queue = queue[1:]

		if currFlag.x == endX && intAbsDiff(currFlag.y, endY) == 1 {
			return currFlag.minute + 1
		}

		for _, nextFlag := range currFlag.next(m.height, m.width) {
			if m.clear(nextFlag) {
				key := nextFlag.key(m.height, m.width)
				if _, ok := queued[key]; !ok {
					queue = append(queue, nextFlag)
					queued[key] = empty
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
	m := NewMountain(len(input), len(input[0]))
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

	done1st := mount.traverse(0, true)
	done2nd := mount.traverse(done1st, false)
	done3rd := mount.traverse(done2nd, true)

	fmt.Println(done3rd)
}
