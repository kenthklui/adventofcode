package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type loc struct {
	x, y int
}

func (l loc) up() loc    { return loc{l.x, l.y - 1} }
func (l loc) down() loc  { return loc{l.x, l.y + 1} }
func (l loc) left() loc  { return loc{l.x - 1, l.y} }
func (l loc) right() loc { return loc{l.x + 1, l.y} }
func (l loc) neighbors(xLimit, yLimit int) []loc {
	list := make([]loc, 0, 4)
	if l.x > 1 {
		list = append(list, l.left())
	}
	if l.x < xLimit {
		list = append(list, l.right())
	}
	if l.y > 1 {
		list = append(list, l.up())
	}
	if l.y < yLimit {
		list = append(list, l.down())
	}
	return list
}
func (l loc) corners() []loc { return []loc{l, l.right(), l.down(), l.right().down()} }

type tube struct {
	r        rune
	loc      loc
	neighbor [4]*tube // up, down, left, right
}

func (t *tube) exit(entrance *tube) *tube {
	for _, out := range t.neighbor {
		if out != nil && out != entrance {
			return out
		}
	}
	return nil
}

func parseTubes(input []string) (*tube, map[loc]*tube) {
	var start *tube
	tubeMap := make(map[loc]*tube)

	for lineNum, line := range input {
		for charNum, r := range line {
			l := loc{charNum, lineNum}
			tubeMap[l] = &tube{r: r, loc: l}
		}
	}

	for l, t := range tubeMap {
		switch t.r {
		case '|':
			t.neighbor[0] = tubeMap[l.up()]
			t.neighbor[1] = tubeMap[l.down()]
		case '-':
			t.neighbor[2] = tubeMap[l.left()]
			t.neighbor[3] = tubeMap[l.right()]
		case 'L':
			t.neighbor[0] = tubeMap[l.up()]
			t.neighbor[3] = tubeMap[l.right()]
		case 'J':
			t.neighbor[0] = tubeMap[l.up()]
			t.neighbor[2] = tubeMap[l.left()]
		case '7':
			t.neighbor[1] = tubeMap[l.down()]
			t.neighbor[2] = tubeMap[l.left()]
		case 'F':
			t.neighbor[1] = tubeMap[l.down()]
			t.neighbor[3] = tubeMap[l.right()]
		case 'S':
			start = t
		case '.':
			delete(tubeMap, l)
		default:
			panic("Invalid char")
		}
	}

	up, ok := tubeMap[start.loc.up()]
	if ok && up.neighbor[1] == start {
		start.neighbor[0] = up
	}
	down, ok := tubeMap[start.loc.down()]
	if ok && down.neighbor[0] == start {
		start.neighbor[1] = down
	}
	left, ok := tubeMap[start.loc.left()]
	if ok && left.neighbor[3] == start {
		start.neighbor[2] = left
	}
	right, ok := tubeMap[start.loc.right()]
	if ok && right.neighbor[2] == start {
		start.neighbor[3] = right
	}

	return start, tubeMap
}

func (t *tube) pipe(tubeMap map[loc]*tube) map[loc]*tube {
	p := make(map[loc]*tube)
	p[t.loc] = t

	for prev, curr := t, t.exit(nil); curr != t; prev, curr = curr, curr.exit(prev) {
		p[curr.loc] = curr
	}

	return p
}

type cornerMap struct {
	width, height int
	grid          map[loc]bool // 0: new, 1: filled
	pipe          map[loc]*tube
}

func makeCornerMap(width, height int, pipe map[loc]*tube) *cornerMap {
	grid := make(map[loc]bool)
	for x := 0; x <= width; x++ {
		for y := 0; y <= height; y++ {
			grid[loc{x, y}] = false
		}
	}
	return &cornerMap{width, height, grid, pipe}
}

func (cm *cornerMap) recursiveFill(l loc) {
	if cm.grid[l] {
		return
	}
	cm.grid[l] = true

	for _, n := range l.neighbors(cm.width, cm.height) {
		if !cm.pipeBlocked(l, n) {
			cm.recursiveFill(n)
		}
	}
}

func (cm *cornerMap) pipeBlocked(source, dest loc) bool {
	if source.x != dest.x {
		leftX, rightX := source.x, dest.x
		if leftX > rightX {
			leftX, rightX = rightX, leftX
		}

		if source.y > 0 {
			topPipe, ok := cm.pipe[loc{leftX, source.y - 1}]
			return ok && topPipe.neighbor[1] != nil
		} else {
			bottomPipe, ok := cm.pipe[loc{leftX, source.y}]
			return ok && bottomPipe.neighbor[0] != nil
		}
	}
	if source.y != dest.y {
		upY, downY := source.y, dest.y
		if upY > downY {
			upY, downY = downY, upY
		}

		if source.x > 0 {
			leftPipe, ok := cm.pipe[loc{source.x - 1, upY}]
			return ok && leftPipe.neighbor[3] != nil
		} else {
			rightPipe, ok := cm.pipe[loc{source.x, upY}]
			return ok && rightPipe.neighbor[2] != nil
		}
	}

	panic("Same location")
}

func enclosed(width, height int, pipe map[loc]*tube) int {
	cm := makeCornerMap(width, height, pipe)
	cm.recursiveFill(loc{0, 0})

	count := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			l := loc{x, y}

			// Is it part of the loop?
			if _, ok := pipe[l]; ok {
				continue
			}

			// Does it border a filled corner?
			filled := false
			for _, c := range l.corners() {
				filled = filled || cm.grid[c]
			}

			// If not, it is enclosed by the loop
			if !filled {
				count++
			}
		}
	}
	return count
}

func main() {
	input := util.StdinReadlines()
	start, tubeMap := parseTubes(input)
	pipe := start.pipe(tubeMap)
	enclosedCount := enclosed(len(input), len(input[0]), pipe)
	fmt.Println(enclosedCount)
}
