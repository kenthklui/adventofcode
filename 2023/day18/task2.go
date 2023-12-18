package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type trench struct {
	min, max, value int
}

type intersect struct {
	value, typ int
}

type lagoon struct {
	hori, vert []trench
}

func (l lagoon) yBreakpoints() []int {
	breakpoints := make([]int, 0, len(l.hori))
	for _, h := range l.hori {
		breakpoints = append(breakpoints, h.value)
	}
	sort.Ints(breakpoints)
	return slices.Compact(breakpoints)
}

func (l lagoon) rowVolume(y int) int {
	intersects := make([]intersect, 0, len(l.vert))
	for _, t := range l.vert {
		if y == t.min {
			intersects = append(intersects, intersect{t.value, 0})
		} else if y == t.max {
			intersects = append(intersects, intersect{t.value, 2})
		} else if y > t.min && y < t.max {
			intersects = append(intersects, intersect{t.value, 1})
		}
	}

	rowVol := 0
	inside := -1
	for i := 0; i < len(intersects); i++ {
		inter := intersects[i]
		if inter.typ == 1 {
			if inside != -1 {
				rowVol += inter.value - inside + 1
				inside = -1
			} else {
				inside = inter.value
			}
		} else {
			next := intersects[i+1]
			i++
			if inside != -1 {
				if inter.typ != next.typ {
					rowVol += next.value - inside + 1
					inside = -1
				}
			} else {
				if inter.typ != next.typ {
					inside = inter.value
				} else {
					rowVol += next.value - inter.value + 1
				}
			}
		}
	}
	return rowVol
}

func (l lagoon) volume() int {
	bps := l.yBreakpoints()

	vol := 0
	// Add volume at breakpoints
	for _, y := range bps {
		vol += l.rowVolume(y)
	}
	// Add volume in-between breakpoints
	for i := range bps[1:] {
		rowCount := bps[i+1] - bps[i] - 1
		vol += rowCount * l.rowVolume(bps[i]+1)
	}

	return vol
}

func parseLagoon(input []string) lagoon {
	horizontals, verticals := make([]trench, 0), make([]trench, 0)
	x, y := 0, 0
	for _, line := range input {
		tokens := strings.Split(line, " ")
		if distance, err := strconv.ParseInt(tokens[2][2:7], 16, 64); err == nil {
			switch int(tokens[2][7] - '0') {
			case 0: // right
				newX := x + int(distance)
				horizontals = append(horizontals, trench{x, newX, y})
				x = newX
			case 1: // down
				newY := y + int(distance)
				verticals = append(verticals, trench{y, newY, x})
				y = newY
			case 2: // left
				newX := x - int(distance)
				horizontals = append(horizontals, trench{newX, x, y})
				x = newX
			case 3: // up
				newY := y - int(distance)
				verticals = append(verticals, trench{newY, y, x})
				y = newY
			default:
				panic("Invalid direction")
			}
		} else {
			panic(err)
		}
	}
	slices.SortFunc(horizontals, func(t1, t2 trench) int { return cmp.Compare(t1.value, t2.value) })
	slices.SortFunc(verticals, func(t1, t2 trench) int { return cmp.Compare(t1.value, t2.value) })

	return lagoon{horizontals, verticals}
}

func main() {
	input := util.StdinReadlines()
	l := parseLagoon(input)
	fmt.Println(l.volume())
}
