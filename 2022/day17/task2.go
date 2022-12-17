package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cave struct {
	rockCount, jetCount, floor, height int
	currRock                           *rock
	space                              [][7]*rock
	pattern                            []int8
}

func NewCave() *cave {
	return &cave{
		rockCount: 0,
		jetCount:  0,
		floor:     0,
		height:    0,
		currRock:  nil,
		space:     make([][7]*rock, 0),
		pattern:   make([]int8, 0),
	}
}

func (c *cave) tetrisMove(times int) {
	for i := 0; i < times; i++ {
		c.addRock()
		cont := true
		for cont {
			c.applyJet()
			cont = c.dropRock()
		}
	}
	c.truncateOld()
}

func (c *cave) setSpace(p pnt, val *rock) {
	y := p.y - c.floor
	c.space[y][p.x] = val
}

func (c *cave) getSpace(p pnt) *rock {
	y := p.y - c.floor
	return c.space[y][p.x]
}

func (c *cave) getSpaceRow(y int) [7]*rock {
	y -= c.floor
	return c.space[y]
}

func (c *cave) addRock() {
	c.currRock = NewRock(c.rockCount%5, c.height)
	c.height = c.currRock.pos.y + 1
	c.rockCount++

	spaceHeight := c.height - c.floor
	for len(c.space) <= spaceHeight {
		c.space = append(c.space, [7]*rock{nil, nil, nil, nil, nil, nil, nil})
	}
	for _, p := range c.currRock.pnts() {
		c.setSpace(p, c.currRock)
	}
}

func (c *cave) applyJet() {
	direction := c.pattern[c.jetCount%len(c.pattern)]
	c.jetCount++

	var add, rem []pnt
	switch direction {
	case -1:
		add, rem = c.currRock.leftDelta()
	case 1:
		add, rem = c.currRock.rightDelta()
	}
	// Collision check
	for _, p := range add {
		if p.x < 0 || p.x >= 7 {
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

func (c *cave) dropRock() bool {
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

func (c *cave) findCutoff() int {
	// Top-down, row-by-row rainfilling algorithm to find safe cutoff
	topRowNum := c.height - c.floor

	var prevRow [7]bool
	for i, v := range c.space[topRowNum] {
		prevRow[i] = (v == nil) // fill top row by default
	}

	for rowNum := topRowNum - 1; rowNum >= 0; rowNum-- {
		var row [7]bool
		emptyRow := true
		for j, v := range row {
			if v || c.space[rowNum][j] != nil || !prevRow[j] { // skip if filled
				continue
			}
			row[j] = true // fill cell
			emptyRow = false
			for k := j - 1; k >= 0; k-- { // fill leftward
				if row[k] || c.space[rowNum][k] != nil {
					break
				}
				row[k] = true
			}
			for k := j + 1; k < 7; k++ { // fill rightward
				if c.space[rowNum][k] != nil {
					break
				}
				row[k] = true
			}
		}

		if emptyRow { // No cell was filled on this row
			return rowNum
		}
		prevRow = row
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
	jetStr := fmt.Sprintf("%06d", c.jetCount%len(c.pattern))
	b.WriteString(jetStr)
	return b.String()
}

// This isn't needed if the iteration count is divisible
func (c *cave) restore(sig string) {
	c.space = make([][7]*rock, len(sig))
	fakeRock := NewRock(0, 0)
	for i, r := range sig {
		for j := 0; j < 7; j++ {
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

type pnt struct {
	x, y int
}

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

type snapshot struct {
	sig                 string
	heightInc, floorInc int
}

func (c *cave) cachedDropRocks(goal int) int {
	cycleLength := 10000 // Assume goal is divisble for now
	snaps := make(map[string]snapshot)

	var sig string
	for {
		if _, ok := snaps[sig]; ok {
			break
		}

		prevHeight := c.height
		prevFloor := c.floor
		c.tetrisMove(cycleLength)

		newSig := c.signature()
		snaps[sig] = snapshot{
			sig:       newSig,
			heightInc: c.height - prevHeight,
			floorInc:  c.floor - prevFloor,
		}

		sig = newSig
	}

	rockCount := c.rockCount
	height := c.height
	floor := c.floor
	for rockCount < goal {
		// TODO: Handle non divisor cases
		snap, _ := snaps[sig]
		rockCount += cycleLength
		sig = snap.sig
		height += snap.heightInc
		floor += snap.floorInc
	}
	return height
}

func main() {
	input := readInput()
	c := parseInput(input)

	goal := 1000000000000
	height := c.cachedDropRocks(goal)
	fmt.Println(height)
}
