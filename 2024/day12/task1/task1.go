package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type plot struct {
	typ   byte
	perim int
	reg   *region
}

type region struct {
	typ   byte
	plots []*plot
}

func (r *region) perim() int {
	perim := 0
	for _, p := range r.plots {
		perim += p.perim
	}
	return perim
}

func (r *region) fence() int {
	return len(r.plots) * r.perim()
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
	g.setPerim()
	g.setRegions()

	return &g
}

func (g garden) at(x, y int) *plot {
	if x < 0 || y < 0 || y >= len(g.plots) || x >= len(g.plots[y]) {
		return nil
	} else {
		return g.plots[y][x]
	}
}

func (g *garden) setPerim() {
	for y := range g.plots {
		g.plots[y][0].perim++
		g.plots[y][len(g.plots[y])-1].perim++

		for x := range g.plots[y][1:] {
			if g.plots[y][x].typ != g.plots[y][x+1].typ {
				g.plots[y][x].perim++
				g.plots[y][x+1].perim++
			}
		}
	}
	for x := range g.plots[0] {
		g.plots[0][x].perim++
		g.plots[len(g.plots)-1][x].perim++

		for y := range g.plots[1:] {
			if g.plots[y][x].typ != g.plots[y+1][x].typ {
				g.plots[y][x].perim++
				g.plots[y+1][x].perim++
			}
		}
	}
}

func (g *garden) setRegions() {
	for y, row := range g.plots {
		for x, p := range row {
			if p.reg == nil {
				r := region{typ: p.typ}
				g.expandRegion(x, y, &r)
				g.regions = append(g.regions, &r)
			}
		}
	}
}

func (g *garden) expandRegion(x, y int, r *region) {
	p := g.plots[y][x]
	p.reg, r.plots = r, append(r.plots, p)
	if neighbor := g.at(x-1, y); neighbor != nil && neighbor.reg == nil && neighbor.typ == r.typ {
		g.expandRegion(x-1, y, r)
	}
	if neighbor := g.at(x+1, y); neighbor != nil && neighbor.reg == nil && neighbor.typ == r.typ {
		g.expandRegion(x+1, y, r)
	}
	if neighbor := g.at(x, y-1); neighbor != nil && neighbor.reg == nil && neighbor.typ == r.typ {
		g.expandRegion(x, y-1, r)
	}
	if neighbor := g.at(x, y+1); neighbor != nil && neighbor.reg == nil && neighbor.typ == r.typ {
		g.expandRegion(x, y+1, r)
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
