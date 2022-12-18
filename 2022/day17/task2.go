package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const rockCycle = 5
const caveWidth = 7

type pnt struct {
	x, y int
}

// Helper for creating []pnt slices
func pkgPnts(coords ...int) []pnt {
	pnts := make([]pnt, 0, len(coords)/2)
	for i := 0; i < len(coords); i += 2 {
		pnts = append(pnts, pnt{coords[i], coords[i+1]})
	}
	return pnts
}

//
// Cave stuff
//

type cave struct {
	rockCount, jetCycle, floor, height int
	currRock                           *rock
	space                              [][caveWidth]*rock
	pattern                            []int
}

func NewCave() *cave {
	return &cave{
		rockCount: 0,
		jetCycle:  0,
		floor:     0,
		height:    0,
		currRock:  nil,
		space:     make([][caveWidth]*rock, 0),
		pattern:   make([]int, 0),
	}
}

func (c *cave) setSpace(p pnt, val *rock)          { c.space[p.y-c.floor][p.x] = val }
func (c *cave) getSpace(p pnt) *rock               { return c.space[p.y-c.floor][p.x] }
func (c *cave) getSpaceRow(y int) [caveWidth]*rock { return c.space[y-c.floor] }

func (c *cave) addRock() {
	c.currRock = NewRock(c.rockCount%5, c.height)
	c.height = c.currRock.pos.y + 1
	c.rockCount++

	for len(c.space) <= c.height-c.floor {
		var emptyRow [caveWidth]*rock
		c.space = append(c.space, emptyRow)
	}
	for _, p := range c.currRock.pnts() {
		c.setSpace(p, c.currRock)
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
	// Collision check
	for _, p := range add {
		if p.x < 0 || p.x >= caveWidth {
			return
		}
		if c.getSpace(p) != nil {
			return
		}
	}

	for _, p := range add {
		c.setSpace(p, c.currRock)
	}
	for _, p := range rem {
		c.setSpace(p, nil)
	}

	c.currRock.pos.x += direction
}

func (c *cave) rockFall() bool {
	add, rem := c.currRock.downDelta()
	// Collision check
	for _, p := range add {
		if p.y < c.floor {
			return false
		}
		if c.getSpace(p) != nil {
			return false
		}
	}

	for _, p := range add {
		c.setSpace(p, c.currRock)
	}
	for _, p := range rem {
		c.setSpace(p, nil)
	}

	heightDrop := true
	for _, r := range c.getSpaceRow(c.currRock.pos.y) {
		heightDrop = heightDrop && (r == nil)
	}
	if heightDrop {
		c.height--
	}

	c.currRock.pos.y--
	return true
}

// Top-down, row-by-row algorithm to find a safe cutoff
// More or less simulates rain falling and see how far down it reaches
func (c *cave) findCutoff() int {
	topRowNum := c.height - c.floor
	var prevFill [caveWidth]bool
	for i, v := range c.space[topRowNum] {
		prevFill[i] = (v == nil) // fill top row empty spaces
	}

	for rowNum := topRowNum - 1; rowNum >= 0; rowNum-- {
		row := c.space[rowNum]
		// Fill each cell if it is empty and above was filled
		var newFill [caveWidth]bool
		cellFilled := false
		for x, r := range row {
			newFill[x] = (r == nil) && prevFill[x]
			cellFilled = cellFilled || newFill[x]
		}
		if !cellFilled { // Terminate if no cell was filled on this row
			return rowNum
		}

		// Fill sideways neighbors
		for x, filled := range newFill {
			if filled {
				for l := x - 1; l >= 0 && row[l] == nil && !newFill[l]; l-- {
					newFill[l] = true
				}
				for r := x + 1; r < caveWidth && row[r] == nil && !newFill[r]; r++ {
					newFill[r] = true
				}
			}
		}
		prevFill = newFill
	}

	return 0
}

func (c *cave) truncateOld() {
	smallest := c.findCutoff()
	c.floor += smallest
	c.space = c.space[smallest:]
}

func (c *cave) signature() string {
	var b strings.Builder

	jetStr := strconv.FormatInt(int64(c.jetCycle), 32)
	b.WriteString(jetStr)
	b.WriteString("-")

	for _, row := range c.space {
		var i byte
		for _, r := range row {
			i <<= 1
			if r != nil {
				i++
			}
		}
		if i == 0 {
			break
		}

		b.WriteByte(i)
	}

	return b.String()
}

func (c *cave) restore(sig string) {
	splits := strings.SplitN(sig, "-", 2)

	if jetCycle, err := strconv.ParseInt(splits[0], 32, 64); err == nil {
		c.jetCycle = int(jetCycle)
	} else {
		panic(err)
	}

	spaceStr := splits[1]
	fakeRock := NewRock(0, 0)
	c.space = make([][caveWidth]*rock, len(spaceStr))
	for i, r := range spaceStr {
		for j := 0; j < caveWidth; j++ {
			if (r>>j)&1 == 1 {
				c.space[i][6-j] = fakeRock
			}
		}
	}
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

type snapshot struct {
	sig                 string
	heightInc, floorInc int
}

func (c *cave) dropRocks(rocks int) {
	cycleLength := rockCycle * rockCycle

	// Super cycle: exponentially cache larger jump sizes
	if rocks >= cycleLength*rockCycle {
		snaps := make(map[string]snapshot)
		sig := c.signature()
		for rocks >= cycleLength {
			if _, cached := snaps[sig]; cached {
				c.cachedDropRocks(&rocks, cycleLength, snaps)
				break
			}

			storedHeight := c.height
			storedFloor := c.floor
			for i := 0; i < cycleLength; i++ {
				c.dropOneRock()
				rocks--
			}
			c.truncateOld()
			snap := snapshot{
				sig:       c.signature(),
				heightInc: c.height - storedHeight,
				floorInc:  c.floor - storedFloor,
			}
			snaps[sig] = snap
			sig = snap.sig
		}
	}

	for ; rocks > 0; rocks-- {
		c.dropOneRock()
	}
}

func (c *cave) cachedDropRocks(rocks *int, cycleLength int, snaps map[string]snapshot) {
	sig := c.signature()

	// Super cycle: exponentially cache larger jump sizes
	if *rocks >= cycleLength*rockCycle {
		newSnaps := make(map[string]snapshot)
		newCycle := cycleLength * rockCycle
		for *rocks >= newCycle {
			if _, cached := newSnaps[sig]; cached {
				c.restore(sig)
				c.cachedDropRocks(rocks, newCycle, newSnaps)
				sig = c.signature()
				break
			}

			prevSig := sig
			storedHeight := c.height
			storedFloor := c.floor
			for i := 0; i < rockCycle; i++ {
				if snap, ok := snaps[sig]; ok {
					sig = snap.sig
					*rocks -= cycleLength

					c.rockCount += cycleLength
					c.height += snap.heightInc
					c.floor += snap.floorInc
				} else {
					panic("Unknown sig")
				}
			}
			snap := snapshot{
				sig:       sig,
				heightInc: c.height - storedHeight,
				floorInc:  c.floor - storedFloor,
			}
			newSnaps[prevSig] = snap
		}
	}

	for *rocks >= cycleLength {
		if snap, ok := snaps[sig]; ok {
			sig = snap.sig
			*rocks -= cycleLength

			c.rockCount += cycleLength
			c.height += snap.heightInc
			c.floor += snap.floorInc
		} else {
			panic("Unknown sig")
		}
	}

	c.restore(sig)
}

//
// Rock stuff
//

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

//
// I/O & main
//

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

func main() {
	input := readInput()
	c := parseInput(input)

	const goal = 1000000000000
	c.dropRocks(goal)
	fmt.Println(c.height)
}
