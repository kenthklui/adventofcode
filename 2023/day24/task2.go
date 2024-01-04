package main

import (
	"fmt"
	"math/big"

	"github.com/kenthklui/adventofcode/util"
)

var zero = big.NewRat(0, 1)

type hailstone struct {
	x, y, z, dx, dy, dz int
}

func (h hailstone) coordinateSum() int { return h.x + h.y + h.z }

type equations struct {
	lines [][]*big.Rat
}

func newEquations(ints [][]int) *equations {
	lines := make([][]*big.Rat, len(ints))
	for i, row := range ints {
		lines[i] = make([]*big.Rat, len(row))
		for j, integer := range row {
			lines[i][j] = big.NewRat(int64(integer), 1)
		}
	}
	return &equations{lines}
}

func (e *equations) swapRow(i, j int) { e.lines[i], e.lines[j] = e.lines[j], e.lines[i] }
func (e *equations) scaleRow(i int, r *big.Rat) {
	for j := range e.lines[i] {
		e.lines[i][j].Mul(e.lines[i][j], r)
	}
}
func (e *equations) addScaledRow(i, j int, r *big.Rat) {
	scalar := big.NewRat(0, 1)
	for h := range e.lines[i] {
		scalar.Mul(r, e.lines[j][h])
		e.lines[i][h].Add(e.lines[i][h], scalar)
	}
}

// Solve a system of linear equations where each lines is:
// lines[0] * x_0 + lines[1] * x_1 + ... + lines[k] * x_k = lines[k+1]
// (ie. the final column is the constant on the other side)
// If equation is overdetermined or underdetermined, return nil
func (e *equations) gaussianElim() *solution {
	if len(e.lines) >= len(e.lines[0]) {
		return nil // overdetermined
	}

	if !e.toRowEchelon() {
		return nil
	}

	if !e.backwardSubstitution() {
		return nil
	}

	constantCol := len(e.lines)
	values := make([]*big.Rat, len(e.lines))
	for i := range values {
		values[i] = e.lines[i][constantCol]
	}

	return &solution{values}
}

func (e *equations) toRowEchelon() bool {
	scalar := big.NewRat(0, 1)
	for i, line1 := range e.lines {
		if line1[i].Cmp(zero) == 0 {
			for j, line2 := range e.lines[i+1:] {
				if line2[i].Cmp(zero) != 0 {
					e.swapRow(i, j+i+1)
					break
				}
			}
			if line1[i].Cmp(zero) == 0 {
				return false // Cannot find a non-zero coefficient for variable
			}
		}

		for j, line2 := range e.lines[i+1:] {
			if line2[i].Cmp(zero) == 0 {
				continue
			}

			j += i + 1
			scalar.Inv(line1[i])
			scalar.Mul(scalar, line2[i])
			scalar.Neg(scalar)
			e.addScaledRow(j, i, scalar)
		}
	}
	return true
}

func (e *equations) backwardSubstitution() bool {
	scalar := big.NewRat(0, 1)
	for i := len(e.lines) - 1; i >= 0; i-- {
		if e.lines[i][i].Cmp(zero) == 0 {
			return false
		}
		e.scaleRow(i, scalar.Inv(e.lines[i][i]))

		for j := 0; j < i; j++ {
			e.addScaledRow(j, i, scalar.Neg(e.lines[j][i]))
		}
	}
	return true
}

type solution struct {
	values []*big.Rat
}

func (s solution) allInt() bool {
	for _, v := range s.values {
		if !v.IsInt() {
			return false
		}
	}
	return true
}

func (s solution) hailstone() hailstone {
	return hailstone{
		x:  int(s.values[0].Num().Int64()),
		y:  int(s.values[1].Num().Int64()),
		z:  int(s.values[2].Num().Int64()),
		dx: int(s.values[3].Num().Int64()),
		dy: int(s.values[4].Num().Int64()),
		dz: int(s.values[5].Num().Int64()),
	}
}

func solve(hailstones []hailstone) hailstone {
	// Some combos of hailstones don't offer a solution. Try with every consecutive 4!
	for i, h1 := range hailstones {
		eqInput := make([][]int, 0, 6)
		for _, h2 := range hailstones[i+1 : i+4] {
			// Obtained by first solving for collision time t:
			// -t = (x-x1)/(dx-dx1) = (y-y1)/(dy-dy1) = (z-z1)/(dz-dz1)
			// Then, use 2 stones to eliminate cross terms and flatten to linear equations
			eqInput = append(eqInput, []int{
				h1.dy - h2.dy,
				h2.dx - h1.dx,
				0,
				h2.y - h1.y,
				h1.x - h2.x,
				0,
				(h2.y*h2.dx - h1.y*h1.dx) - (h2.x*h2.dy - h1.x*h1.dy),
			})
			eqInput = append(eqInput, []int{
				h1.dz - h2.dz,
				0,
				h2.dx - h1.dx,
				h2.z - h1.z,
				0,
				h1.x - h2.x,
				(h2.z*h2.dx - h1.z*h1.dx) - (h2.x*h2.dz - h1.x*h1.dz),
			})
		}
		eq := newEquations(eqInput)

		if sol := eq.gaussianElim(); sol != nil && sol.allInt() {
			return sol.hailstone()
		}
	}
	panic("Solution not found")
}

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

func main() {
	input := util.StdinReadlines()
	hailstones := parse(input)
	stone := solve(hailstones)
	fmt.Println(stone.coordinateSum())
}
