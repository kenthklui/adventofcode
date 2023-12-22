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
	axis, length, minX, minY, minZ, maxX, maxY, maxZ int
	blockedBy                                        map[*brick]void
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
	b := brick{blockedBy: make(map[*brick]void)}

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
	bricks   map[*brick]void
	maxZ     int
}

func makeArea(bricks []*brick) *area {
	slices.SortFunc(bricks, func(b1, b2 *brick) int {
		return cmp.Compare(b1.minZ, b2.minZ)
	})
	brickMap := make(map[vec]*brick)
	maxZ := 0
	a := area{brickMap, make(map[*brick]void), maxZ}
	for _, b := range bricks {
		a.dropIn(b)
	}
	return &a
}

func (a *area) droppable(b *brick) (bool, map[*brick]void) {
	a.bricks[b] = nul

	blockedBy := make(map[*brick]void)
	if b.minZ == floorZ {
		return false, blockedBy
	}

	droppable := true
	for _, v := range b.dropVecs(1) {
		if blocker, ok := a.brickMap[v]; ok {
			droppable = false
			blockedBy[blocker] = nul
		}
	}
	return droppable, blockedBy
}

func (a *area) dropIn(b *brick) {
	for {
		if droppable, blockedBy := a.droppable(b); droppable {
			b.minZ--
			b.maxZ--
		} else {
			b.blockedBy = blockedBy
			for _, v := range b.vecs() {
				a.brickMap[v] = b
			}
			break
		}
	}

}

func (a *area) deleteBrick(brickToDelete *brick) bool {
	for b := range a.bricks {
		if _, ok := b.blockedBy[brickToDelete]; ok && len(b.blockedBy) == 1 {
			return false
		}
	}

	for b := range a.bricks {
		if _, ok := b.blockedBy[brickToDelete]; ok && len(b.blockedBy) == 1 {
			delete(b.blockedBy, brickToDelete)
		}
	}
	for _, v := range brickToDelete.vecs() {
		delete(a.brickMap, v)
	}
	return true
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
		if a.deleteBrick(b) {
			count++
		}
	}
	fmt.Println(count)
}
