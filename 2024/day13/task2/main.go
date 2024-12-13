package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

var COSTS = [2]int{3, 1}
var OFFSET = 10000000000000

// TODO: Move all system of linear equations related code to a separate package
var zero = big.NewRat(0, 1)

type equations struct {
	lines [][]*big.Rat
}

type solution struct {
	values []*big.Rat
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

func parse(input []string) *equations {
	ints := util.ParseInts(input)
	eqInts := [][]int{
		{ints[0][0], ints[1][0], ints[2][0] + OFFSET},
		{ints[0][1], ints[1][1], ints[2][1] + OFFSET},
	}
	return newEquations(eqInts)
}

func solve(input []string) (output string) {
	tokens := 0
	for i := 0; i < len(input); i += 4 {
		eq := parse(input[i : i+3])
		sol := eq.gaussianElim()
		if sol.values[0].IsInt() && sol.values[1].IsInt() {
			tokens += int(sol.values[0].Num().Int64())*COSTS[0] + int(sol.values[1].Num().Int64())*COSTS[1]
		}
	}
	return strconv.Itoa(tokens)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
