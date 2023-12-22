package main

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/kenthklui/adventofcode/util"
)

const floorZ = 1

type void struct{}

var nul void

type vec struct{ x, y, z int }

func (v vec) drop(dz int) vec { return vec{v.x, v.y, v.z - dz} }

type brick struct {
	index, axis, length, minX, minY, minZ, maxX, maxY, maxZ int
	below, above                                            []*brick
}

func (b brick) vecs() []vec {
	vecs := make([]vec, 0, b.length)
	switch b.axis {
	case 0:
		vecs = make([]vec, 0, b.maxX-b.minX+1)
		for x := b.minX; x <= b.maxX; x++ {
			vecs = append(vecs, vec{x, b.minY, b.minZ})
		}
	case 1:
		vecs = make([]vec, 0, b.maxY-b.minY+1)
		for y := b.minY; y <= b.maxY; y++ {
			vecs = append(vecs, vec{b.minX, y, b.minZ})
		}
	case 2:
		vecs = make([]vec, 0, b.maxZ-b.minZ+1)
		for z := b.minZ; z <= b.maxZ; z++ {
			vecs = append(vecs, vec{b.minX, b.minY, z})
		}
	}
	return vecs
}

func (b brick) dropVecs(dz int) []vec {
	vs := b.vecs()
	for i, v := range vs {
		vs[i] = v.drop(dz)
	}
	return vs
}

const lineFormat = "%d,%d,%d~%d,%d,%d"

func parseBrick(line string) *brick {
	b := brick{below: make([]*brick, 0), above: make([]*brick, 0)}
	if n, err := fmt.Sscanf(line, lineFormat,
		&b.minX, &b.minY, &b.minZ, &b.maxX, &b.maxY, &b.maxZ); err != nil {
		panic(err)
	} else if n != 6 {
		panic("Failed to parse 6 coordinates")
	}
	if b.minX != b.maxX {
		b.axis = 0
		b.length = b.maxX - b.minX + 1
	} else if b.minY != b.maxY {
		b.axis = 1
		b.length = b.maxY - b.minY + 1
	} else if b.minZ != b.maxZ {
		b.axis = 2
		b.length = b.maxZ - b.minZ + 1
	}

	return &b
}

type area struct {
	brickMap map[vec]*brick
	bricks   []*brick
}

func makeArea(bricks []*brick) *area {
	slices.SortFunc(bricks, func(b1, b2 *brick) int {
		return cmp.Compare(b1.minZ, b2.minZ)
	})
	a := &area{make(map[vec]*brick), bricks}
	for i, b := range bricks {
		b.index = i
		a.dropIn(b)
	}
	return a
}

func (a *area) droppable(b *brick) (bool, []*brick) {
	bricksBelow := make([]*brick, 0, b.length)

	if b.minZ == floorZ {
		return false, bricksBelow
	}

	droppable := true
	blockers := make(map[*brick]void)
	for _, v := range b.dropVecs(1) {
		if b, ok := a.brickMap[v]; ok {
			droppable = false
			blockers[b] = nul
		}
	}
	for blocker := range blockers {
		bricksBelow = append(bricksBelow, blocker)
	}

	return droppable, bricksBelow
}

func (a *area) dropIn(b *brick) {
	for {
		if droppable, below := a.droppable(b); droppable {
			b.minZ--
			b.maxZ--
		} else {
			b.below = below
			for _, blockBelow := range b.below {
				blockBelow.above = append(blockBelow.above, b)
			}

			for _, v := range b.vecs() {
				a.brickMap[v] = b
			}
			break
		}
	}
}

func (a *area) ifDeleteBrick(brickToDelete *brick) int {
	fallingBricks := 0

	bricksMoved := make([]bool, len(a.bricks))
	bricksMoved[brickToDelete.index] = true

	bricksToCheck := make([]*brick, 0)
	for _, blockAbove := range brickToDelete.above {
		bricksToCheck = append(bricksToCheck, blockAbove)
	}

	bricksChecked := 0
	for len(bricksToCheck) > 0 {
		bricksChecked++
		b := bricksToCheck[0]
		bricksToCheck = bricksToCheck[1:]

		if bricksMoved[b.index] {
			continue
		}

		keep := false
		for _, blockBelow := range b.below {
			if !bricksMoved[blockBelow.index] {
				keep = true
				break
			}
		}
		if keep {
			continue
		}

		bricksMoved[b.index] = true
		fallingBricks++
		for _, blockAbove := range b.above {
			bricksToCheck = append(bricksToCheck, blockAbove)
		}
	}
	return fallingBricks
}

func parseBricks(input []string) []*brick {
	bricks := make([]*brick, 0, len(input))
	for _, line := range input {
		bricks = append(bricks, parseBrick(line))
	}
	slices.SortFunc(bricks, func(b1, b2 *brick) int {
		return cmp.Compare(b1.minZ, b2.minZ)
	})
	return bricks
}

func main() {
	input := util.StdinReadlines()
	bricks := parseBricks(input)

	a := makeArea(bricks)
	count := 0
	for _, b := range bricks {
		count += a.ifDeleteBrick(b)
	}
	fmt.Println(count)
}
