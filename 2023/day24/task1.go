package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct {
	x, y float64
}

func (v vec2) negate() vec2          { return vec2{-v.x, -v.y} }
func (v vec2) scalar(f float64) vec2 { return vec2{v.x * f, v.y * f} }

func sum2d(v1, v2 vec2) vec2      { return vec2{v1.x + v2.x, v1.y + v2.y} }
func subtract2d(v1, v2 vec2) vec2 { return sum2d(v1, v2.negate()) }
func dot2d(v1, v2 vec2) float64   { return v1.x*v2.x + v1.y*v2.y }
func cross2d(v1, v2 vec2) float64 { return v1.x*v2.y - v1.y*v2.x }
func parallel2d(v1, v2 vec2) bool { return cross2d(v1, v2) == 0 }

type intersection struct {
	t1, t2 float64
	point  vec2
}

type area struct {
	minX, minY, maxX, maxY float64
}

func (i intersection) within(a area) bool {
	if i.t1 < 0.0 {
		// fmt.Println("Crossed in the past for hailstone A.")
		return false
	}
	if i.t2 < 0.0 {
		// fmt.Println("Crossed in the past for hailstone A.")
		return false
	}

	return a.minX <= i.point.x && i.point.x <= a.maxX &&
		a.minY <= i.point.y && i.point.y <= a.maxY
}

// solve v1 + k1 * dv1 == v2 + k2 * dv2
func solve(v1, dv1, v2, dv2 vec2) intersection {
	numerT1 := cross2d(subtract2d(v2, v1), dv2)
	numerT2 := cross2d(subtract2d(v2, v1), dv1)
	denom := cross2d(dv1, dv2)

	t1, t2 := numerT1/denom, numerT2/denom
	point := sum2d(v1, dv1.scalar(t1))

	return intersection{t1, t2, point}
}

type hailstone struct {
	x, y, z, dx, dy, dz int
}

func (h hailstone) pos2d() vec2 { return vec2{float64(h.x), float64(h.y)} }
func (h hailstone) dir2d() vec2 { return vec2{float64(h.dx), float64(h.dy)} }

func parse(input []string) []hailstone {
	hailstones := make([]hailstone, 0, len(input))
	for _, line := range input {
		var h hailstone
		if n, err := fmt.Sscanf(line, "%d, %d, %d @ %d, %d, %d",
			&h.x, &h.y, &h.z, &h.dx, &h.dy, &h.dz); err == nil {
		} else if err != nil {
			panic(err)
		} else if n != 6 {
			panic("Failed to parse coordinates")
		}
		hailstones = append(hailstones, h)
	}
	return hailstones
}

// var testArea = area{7, 7, 27, 27}

var testArea = area{200000000000000, 200000000000000, 400000000000000, 400000000000000}

func main() {
	input := util.StdinReadlines()
	hailstones := parse(input)

	inRange := 0
	for i, h1 := range hailstones {
		for _, h2 := range hailstones[i+1:] {
			v1, dv1, v2, dv2 := h1.pos2d(), h1.dir2d(), h2.pos2d(), h2.dir2d()
			if parallel2d(dv1, dv2) {
				// fmt.Println("Parallel!")
			} else {
				inter := solve(v1, dv1, v2, dv2)
				if inter.within(testArea) {
					inRange++
				}
			}
		}
	}
	fmt.Println(inRange)
}
