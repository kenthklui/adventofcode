package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct{ x, y int }

var dirs = [4]vec2{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func (v vec2) add(v2 vec2) vec2         { return vec2{v.x + v2.x, v.y + v2.y} }
func (v vec2) in(sizeX, sizeY int) bool { return v.x >= 0 && v.x < sizeX && v.y >= 0 && v.y < sizeY }
func (v vec2) edge(dir int) [2]vec2 {
	switch dir {
	case 0:
		return [2]vec2{v, {v.x + 1, v.y}}
	case 1:
		return [2]vec2{{v.x + 1, v.y}, {v.x + 1, v.y + 1}}
	case 2:
		return [2]vec2{{v.x, v.y + 1}, {v.x + 1, v.y + 1}}
	case 3:
		return [2]vec2{{v.x, v.y + 1}, v}
	default:
		panic("Invalid direction")
	}
}

type plot struct {
	typ byte
	reg *region
}

type region struct {
	typ                    byte
	plots                  []*plot
	edgePointH, edgePointV map[vec2]int
}

// Instead of counting sides, count corners
// NOTE: account for the edge case of a polygon having touching corners, ie.
func (r *region) sides() int {
	intersections := 0
	for point, hval := range r.edgePointH {
		if vval, ok := r.edgePointV[point]; ok {
			intersections += (hval + vval) / 2
		}
	}
	return intersections
}

func (r *region) fence() int {
	return len(r.plots) * r.sides()
}

type garden struct {
	plots   [][]*plot
	regions []*region
}

func newGarden(input []string) *garden {
	plots := make([][]*plot, len(input))
	for y, line := range input {
		plots[y] = make([]*plot, len(line))
		for x, c := range line {
			plots[y][x] = &plot{typ: byte(c)}
		}
	}

	g := garden{plots: plots}
	g.setRegions()

	return &g
}

func (g garden) at(v vec2) *plot {
	if v.in(len(g.plots[0]), len(g.plots)) {
		return g.plots[v.y][v.x]
	}
	return nil
}

func (g *garden) setRegions() {
	for y, row := range g.plots {
		for x, p := range row {
			if p.reg == nil {
				r := region{typ: p.typ, edgePointH: make(map[vec2]int), edgePointV: make(map[vec2]int)}
				g.expandRegion(vec2{x, y}, &r)
				g.regions = append(g.regions, &r)
			}
		}
	}
}

func (g *garden) expandRegion(v vec2, r *region) {
	p := g.at(v)
	p.reg = r
	r.plots = append(r.plots, p)

	for i, dxy := range dirs {
		neighbor := v.add(dxy)
		neighborPlot := g.at(neighbor)
		if neighborPlot == nil || neighborPlot.typ != r.typ {
			if i%2 == 0 { // If direction is vertical, then edge between is horizontal
				r.edgePointH[v.edge(i)[0]]++
				r.edgePointH[v.edge(i)[1]]++
			} else {
				r.edgePointV[v.edge(i)[0]]++
				r.edgePointV[v.edge(i)[1]]++
			}
		} else if neighborPlot.reg == nil {
			g.expandRegion(neighbor, r)
		}
	}
}

func (g garden) fences() int {
	fences := 0
	for _, r := range g.regions {
		fences += r.fence()
	}
	return fences
}

func solve(input []string) (output string) {
	g := newGarden(input)
	return strconv.Itoa(g.fences())
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
