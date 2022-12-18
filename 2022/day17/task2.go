package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const rockCycle = 5
const caveWidth = 7

type pnt struct {
	x, y int
}

//
// Cave stuff
//

type cave struct {
	rockCount, jetCycle, floor, height int
	currRock                           *rock
	space                              [][caveWidth]*rock
	pattern                            []int8
}

func NewCave() *cave {
	return &cave{
		rockCount: 0,
		jetCycle:  0,
		floor:     0,
		height:    0,
		currRock:  nil,
		space:     make([][caveWidth]*rock, 0),
		pattern:   make([]int8, 0),
	}
}

func (c *cave) setSpace(p pnt, val *rock) {
	y := p.y - c.floor
	c.space[y][p.x] = val
}

func (c *cave) getSpace(p pnt) *rock {
	y := p.y - c.floor
	return c.space[y][p.x]
}

func (c *cave) getSpaceRow(y int) [caveWidth]*rock {
	y -= c.floor
	return c.space[y]
}

func (c *cave) addRock() {
	c.currRock = NewRock(c.rockCount%5, c.height)
	c.height = c.currRock.pos.y + 1
	c.rockCount++

	spaceHeight := c.height - c.floor
	for len(c.space) <= spaceHeight {
		c.space = append(c.space, [caveWidth]*rock{nil, nil, nil, nil, nil, nil, nil})
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

	c.currRock.pos.x += int(direction)
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

// Top-down, row-by-row rainfilling algorithm to find a safe cutoff
func (c *cave) findCutoff() int {
	topRowNum := c.height - c.floor

	var prevFill [caveWidth]bool
	for i, v := range c.space[topRowNum] {
		prevFill[i] = (v == nil) // fill top row empty spaces
	}

	for rowNum := topRowNum - 1; rowNum >= 0; rowNum-- {
		row := c.space[rowNum]

		var newFill [caveWidth]bool
		cellFilled := false
		for x, r := range row {
			// Skip filed cells
			if newFill[x] {
				continue
			}

			// Fill if cell is empty and direct above is filled, otherwise skip
			if r == nil && prevFill[x] {
				newFill[x] = true
				cellFilled = true

				// Fill neighbors leftward
				for l := x - 1; l >= 0 && row[l] == nil && !newFill[l]; l-- {
					newFill[l] = true
				}
				// Fill neighbors rightward
				for r := x + 1; r < caveWidth && row[r] == nil; r++ {
					newFill[r] = true
				}
			}
		}

		if !cellFilled { // No cell was filled on this row
			return rowNum
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

type snapshots map[string]snapshot

func (snaps snapshots) String() string {
	lines := make([]string, 0, len(snaps))
	for _, s := range snaps {
		ps := []byte(strings.SplitN(s.sig, "-", 2)[1])
		l := fmt.Sprintf("%d %d %v\n", s.heightInc, s.floorInc, ps)
		lines = append(lines, l)
	}
	sort.Strings(lines)
	return strings.Join(lines, "")
}

func (c *cave) dropRocks(rocks int) {
	cycleLength := rockCycle * rockCycle

	if rocks >= cycleLength*rockCycle {
		snaps := make(snapshots)
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

	// Handle remainder
	for ; rocks > 0; rocks-- {
		c.dropOneRock()
	}
}

func (c *cave) cachedDropRocks(rocks *int, cycleLength int, snaps snapshots) {
	sig := c.signature()

	// Super cycle: exponentially cache larger jump sizes
	if *rocks >= cycleLength*rockCycle {
		newSnaps := make(snapshots)
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

	// Base cycle: handle the remainder
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
	pos      pnt
	rockType int
}

func NewRock(rockType, caveHeight int) *rock {
	r := new(rock)
	r.rockType = rockType
	r.setPos(caveHeight)
	return r
}

func (r *rock) setPos(caveHeight int) { r.pos.x, r.pos.y = 2, caveHeight+r.height()+2 }

// In theory, can eliminate all the `switch r.rockType` with an interface
// In practice, way too many duplicated method signatures
func (r *rock) height() int {
	switch r.rockType {
	case 0: // long
		return 1
	case 1: // cross
		return 3
	case 2: // reverse L
		return 3
	case 3: // tall
		return 4
	case 4: // square
		return 2
	default:
		panic("Invalid rock type")
	}
	return -1
}
func (r *rock) pnts() []pnt {
	x, y := r.pos.x, r.pos.y
	switch r.rockType {
	case 0: // long
		return []pnt{pnt{x, y}, pnt{x + 1, y}, pnt{x + 2, y}, pnt{x + 3, y}}
	case 1: // cross
		return []pnt{
			pnt{x + 1, y},
			pnt{x, y - 1}, pnt{x + 1, y - 1}, pnt{x + 2, y - 1},
			pnt{x + 1, y - 2},
		}
	case 2: // reverse L
		return []pnt{
			pnt{x + 2, y}, pnt{x + 2, y - 1}, pnt{x, y - 2}, pnt{x + 1, y - 2}, pnt{x + 2, y - 2},
		}
	case 3: // tall
		return []pnt{pnt{x, y}, pnt{x, y - 1}, pnt{x, y - 2}, pnt{x, y - 3}}
	case 4: // square
		return []pnt{pnt{x, y}, pnt{x + 1, y}, pnt{x, y - 1}, pnt{x + 1, y - 1}}
	default:
		panic("Invalid rock type")
	}
	return nil
}
func (r *rock) leftDelta() ([]pnt, []pnt) {
	var add, rem []pnt
	x, y := r.pos.x, r.pos.y
	switch r.rockType {
	case 0: // long
		add = []pnt{pnt{x - 1, y}}
		rem = []pnt{pnt{x + 3, y}}
	case 1: // cross
		add = []pnt{pnt{x, y}, pnt{x - 1, y - 1}, pnt{x, y - 2}}
		rem = []pnt{pnt{x + 1, y}, pnt{x + 2, y - 1}, pnt{x + 1, y - 2}}
	case 2: // reverse L
		add = []pnt{pnt{x + 1, y}, pnt{x + 1, y - 1}, pnt{x - 1, y - 2}}
		rem = []pnt{pnt{x + 2, y}, pnt{x + 2, y - 1}, pnt{x + 2, y - 2}}
	case 3: // tall
		nx := x - 1
		add = []pnt{pnt{nx, y}, pnt{nx, y - 1}, pnt{nx, y - 2}, pnt{nx, y - 3}}
		rem = []pnt{pnt{x, y}, pnt{x, y - 1}, pnt{x, y - 2}, pnt{x, y - 3}}
	case 4: // square
		add = []pnt{pnt{x - 1, y}, pnt{x - 1, y - 1}}
		rem = []pnt{pnt{x + 1, y}, pnt{x + 1, y - 1}}
	}
	return add, rem
}
func (r *rock) rightDelta() ([]pnt, []pnt) {
	var add, rem []pnt
	x, y := r.pos.x, r.pos.y
	switch r.rockType {
	case 0: // long
		add = []pnt{pnt{x + 4, y}}
		rem = []pnt{pnt{x, y}}
	case 1: // cross
		add = []pnt{pnt{x + 2, y}, pnt{x + 3, y - 1}, pnt{x + 2, y - 2}}
		rem = []pnt{pnt{x + 1, y}, pnt{x, y - 1}, pnt{x + 1, y - 2}}
	case 2: // reverse L
		add = []pnt{pnt{x + 3, y}, pnt{x + 3, y - 1}, pnt{x + 3, y - 2}}
		rem = []pnt{pnt{x + 2, y}, pnt{x + 2, y - 1}, pnt{x, y - 2}}
	case 3: // tall
		nx := x + 1
		add = []pnt{pnt{nx, y}, pnt{nx, y - 1}, pnt{nx, y - 2}, pnt{nx, y - 3}}
		rem = []pnt{pnt{x, y}, pnt{x, y - 1}, pnt{x, y - 2}, pnt{x, y - 3}}
	case 4: // square
		add = []pnt{pnt{x + 2, y}, pnt{x + 2, y - 1}}
		rem = []pnt{pnt{x, y}, pnt{x, y - 1}}
	}
	return add, rem
}
func (r *rock) downDelta() ([]pnt, []pnt) {
	var add, rem []pnt
	x, y := r.pos.x, r.pos.y
	switch r.rockType {
	case 0: // long
		ny := y - 1
		add = []pnt{pnt{x, ny}, pnt{x + 1, ny}, pnt{x + 2, ny}, pnt{x + 3, ny}}
		rem = []pnt{pnt{x, y}, pnt{x + 1, y}, pnt{x + 2, y}, pnt{x + 3, y}}
	case 1: // cross
		add = []pnt{pnt{x, y - 2}, pnt{x + 1, y - 3}, pnt{x + 2, y - 2}}
		rem = []pnt{pnt{x, y - 1}, pnt{x + 1, y}, pnt{x + 2, y - 1}}
	case 2: // reverse L
		add = []pnt{pnt{x, y - 3}, pnt{x + 1, y - 3}, pnt{x + 2, y - 3}}
		rem = []pnt{pnt{x + 2, y}, pnt{x, y - 2}, pnt{x + 1, y - 2}}
	case 3: // tall
		add = []pnt{pnt{x, y - 4}}
		rem = []pnt{pnt{x, y}}
	case 4: // square
		add = []pnt{pnt{x, y - 2}, pnt{x + 1, y - 2}}
		rem = []pnt{pnt{x, y}, pnt{x + 1, y}}
	}
	return add, rem
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
	pattern := make([]int8, len(input[0]))
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
