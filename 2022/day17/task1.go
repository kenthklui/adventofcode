package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const rockCycle = 5
const caveWidth = 7

type pnt struct {
	x, y int
}

func pkgPnts(coords ...int) []pnt {
	pnts := make([]pnt, 0, len(coords)/2)
	for i := 0; i < len(coords); i += 2 {
		pnts = append(pnts, pnt{coords[i], coords[i+1]})
	}
	return pnts
}

type cave struct {
	rockCount, jetCycle, height int
	currRock                    *rock
	space                       [][caveWidth]*rock
	pattern                     []int
}

func NewCave() *cave {
	return &cave{
		rockCount: 0,
		jetCycle:  0,
		height:    0,
		currRock:  nil,
		space:     make([][caveWidth]*rock, 0),
		pattern:   make([]int, 0),
	}
}

func (c *cave) addRock() {
	c.currRock = NewRock(c.rockCount%5, c.height)
	c.height = c.currRock.pos.y + 1
	c.rockCount++

	for len(c.space) <= c.height {
		var emptyRow [caveWidth]*rock
		c.space = append(c.space, emptyRow)
	}
	for _, p := range c.currRock.pnts() {
		c.space[p.y][p.x] = c.currRock
	}
}

func (c *cave) dropOneRock() {
	c.addRock()
	for falling := true; falling; falling = c.rockFall() {
		c.applyJet()
	}
}

func (c *cave) applyJet() {
	direction := c.pattern[c.jetCycle]
	c.jetCycle = (c.jetCycle + 1) % len(c.pattern)

	var add, rem []pnt
	switch direction {
	case -1:
		add, rem = c.currRock.leftDelta()
	case 1:
		add, rem = c.currRock.rightDelta()
	}
	for _, p := range add {
		if p.x < 0 || p.x >= caveWidth {
			return
		}
		if c.space[p.y][p.x] != nil {
			return
		}
	}

	for _, p := range add {
		c.space[p.y][p.x] = c.currRock
	}
	for _, p := range rem {
		c.space[p.y][p.x] = nil
	}

	c.currRock.pos.x += direction
}

func (c *cave) rockFall() bool {
	add, rem := c.currRock.downDelta()
	for _, p := range add {
		if p.y < 0 {
			return false
		}
		if c.space[p.y][p.x] != nil {
			return false
		}
	}

	for _, p := range add {
		c.space[p.y][p.x] = c.currRock
	}
	for _, p := range rem {
		c.space[p.y][p.x] = nil
	}

	heightDrop := true
	for _, r := range c.space[c.currRock.pos.y] {
		heightDrop = heightDrop && (r == nil)
	}
	if heightDrop {
		c.height--
	}

	c.currRock.pos.y--
	return true
}

func (c cave) Print() {
	var b strings.Builder
	for n := len(c.space) - 1; n >= 0; n-- {
		b.WriteRune('|')
		for _, s := range c.space[n] {
			if s == nil {
				b.WriteRune('.')
			} else if s == c.currRock {
				b.WriteRune('@')
			} else {
				b.WriteRune('#')
			}
		}
		b.WriteString("|\n")
	}
	b.WriteString("+-------+\n")
	fmt.Println(b.String())
}

func (c *cave) dropRocks(rocks int) {
	for ; rocks > 0; rocks-- {
		c.dropOneRock()
	}
}

type rock struct {
	pos pnt
	rt  rockType
}

func NewRock(rt, caveHeight int) *rock {
	r := new(rock)
	switch rt {
	case 0:
		r.rt = longRock{}
	case 1:
		r.rt = crossRock{}
	case 2:
		r.rt = lRock{}
	case 3:
		r.rt = tallRock{}
	case 4:
		r.rt = sqrRock{}
	default:
		panic("Invalid rock type")
	}
	r.pos.x, r.pos.y = 2, caveHeight+r.height()+2
	return r
}

func (r *rock) height() int                { return r.rt.height() }
func (r *rock) pnts() []pnt                { return r.rt.pnts(r.pos.x, r.pos.y) }
func (r *rock) leftDelta() ([]pnt, []pnt)  { return r.rt.leftDelta(r.pos.x, r.pos.y) }
func (r *rock) rightDelta() ([]pnt, []pnt) { return r.rt.rightDelta(r.pos.x, r.pos.y) }
func (r *rock) downDelta() ([]pnt, []pnt)  { return r.rt.downDelta(r.pos.x, r.pos.y) }

type rockType interface {
	height() int
	pnts(int, int) []pnt
	leftDelta(int, int) ([]pnt, []pnt)
	rightDelta(int, int) ([]pnt, []pnt)
	downDelta(int, int) ([]pnt, []pnt)
}
type longRock struct{}
type crossRock struct{}
type lRock struct{}
type tallRock struct{}
type sqrRock struct{}

func (lr longRock) height() int { return 1 }
func (lr longRock) pnts(x, y int) []pnt {
	return pkgPnts(x, y, x+1, y, x+2, y, x+3, y)
}
func (lr longRock) leftDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x-1, y), pkgPnts(x+3, y)
}
func (lr longRock) rightDelta(x, y int) ([]pnt, []pnt) {
	return []pnt{pnt{x + 4, y}}, []pnt{pnt{x, y}}
}
func (lr longRock) downDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y-1, x+1, y-1, x+2, y-1, x+3, y-1), pkgPnts(x, y, x+1, y, x+2, y, x+3, y)
}
func (cr crossRock) height() int { return 3 }
func (cr crossRock) pnts(x, y int) []pnt {
	return pkgPnts(x+1, y, x, y-1, x+1, y-1, x+2, y-1, x+1, y-2)
}
func (cr crossRock) leftDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y, x-1, y-1, x, y-2), pkgPnts(x+1, y, x+2, y-1, x+1, y-2)
}
func (cr crossRock) rightDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x+2, y, x+3, y-1, x+2, y-2), pkgPnts(x+1, y, x, y-1, x+1, y-2)
}
func (cr crossRock) downDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y-2, x+1, y-3, x+2, y-2), pkgPnts(x, y-1, x+1, y, x+2, y-1)
}
func (lr lRock) height() int { return 3 }
func (lr lRock) pnts(x, y int) []pnt {
	return pkgPnts(x+2, y, x+2, y-1, x, y-2, x+1, y-2, x+2, y-2)
}
func (lr lRock) leftDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x+1, y, x+1, y-1, x-1, y-2), pkgPnts(x+2, y, x+2, y-1, x+2, y-2)
}
func (lr lRock) rightDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x+3, y, x+3, y-1, x+3, y-2), pkgPnts(x+2, y, x+2, y-1, x, y-2)
}
func (lr lRock) downDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y-3, x+1, y-3, x+2, y-3), pkgPnts(x+2, y, x, y-2, x+1, y-2)
}
func (tr tallRock) height() int { return 4 }
func (tr tallRock) pnts(x, y int) []pnt {
	return pkgPnts(x, y, x, y-1, x, y-2, x, y-3)
}
func (tr tallRock) leftDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x-1, y, x-1, y-1, x-1, y-2, x-1, y-3), pkgPnts(x, y, x, y-1, x, y-2, x, y-3)
}
func (tr tallRock) rightDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x+1, y, x+1, y-1, x+1, y-2, x+1, y-3), pkgPnts(x, y, x, y-1, x, y-2, x, y-3)
}
func (tr tallRock) downDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y-4), pkgPnts(x, y)
}
func (sr sqrRock) height() int { return 2 }
func (sr sqrRock) pnts(x, y int) []pnt {
	return pkgPnts(x, y, x+1, y, x, y-1, x+1, y-1)
}
func (sr sqrRock) leftDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x-1, y, x-1, y-1), pkgPnts(x+1, y, x+1, y-1)
}
func (sr sqrRock) rightDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x+2, y, x+2, y-1), pkgPnts(x, y, x, y-1)
}
func (sr sqrRock) downDelta(x, y int) ([]pnt, []pnt) {
	return pkgPnts(x, y-2, x+1, y-2), pkgPnts(x, y, x+1, y)
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

func parseInput(input []string) *cave {
	pattern := make([]int, len(input[0]))
	for i, r := range input[0] {
		switch r {
		case '<':
			pattern[i] = -1
		case '>':
			pattern[i] = 1
		default:
			panic("Invalid jet direction")
		}
	}

	c := NewCave()
	c.pattern = pattern
	return c
}

const goal = 2022

func main() {
	input := readInput()
	c := parseInput(input)

	c.dropRocks(goal)
	fmt.Println(c.height)
}
