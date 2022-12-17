package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cave struct {
	rockCount, jetCount, height int
	currRock                    *rock
	space                       [][7]*rock
	pattern                     []int8
}

func NewCave() *cave {
	return &cave{
		rockCount: 0,
		jetCount:  0,
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
}

func (c *cave) addRock() {
	c.currRock = NewRock(c.rockCount%5, c.height)
	c.height = c.currRock.pos.y + 1
	c.rockCount++

	for len(c.space) <= c.height {
		c.space = append(c.space, [7]*rock{nil, nil, nil, nil, nil, nil, nil})
	}
	for _, p := range c.currRock.pnts() {
		c.space[p.y][p.x] = c.currRock
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
	c.currRock.pos.x += int(direction)
}

func (c *cave) dropRock() bool {
	add, rem := c.currRock.downDelta()
	// Collision check
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

type pnt struct {
	x, y int
}

type rock struct {
	pos      pnt
	rockType int

	// The following return (added, removed []pnt)
}

func NewRock(rockType, caveHeight int) *rock {
	r := new(rock)
	r.rockType = rockType
	r.setPos(caveHeight)
	return r
}

func (r *rock) setPos(caveHeight int) { r.pos.x, r.pos.y = 2, caveHeight+r.height()+2 }

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

func main() {
	input := readInput()
	c := parseInput(input)

	c.tetrisMove(2022)
	fmt.Println(c.height)
}
