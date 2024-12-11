package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

func digits(i int) int {
	j, mult := 1, 10
	for mult <= i {
		j++
		mult *= 10
	}
	return j
}

func pow(base, exp int) int {
	if exp == 0 {
		return 1
	} else if exp == 1 {
		return base
	} else if exp%2 == 0 {
		return pow(base*base, exp/2)
	} else {
		return base * pow(base*base, (exp-1)/2)
	}
}

type stone struct {
	val   int
	right *stone
}

func (s *stone) step() bool {
	if s.val == 0 {
		s.val = 1
		return false
	} else {
		d := digits(s.val)
		if d%2 == 0 {
			split := pow(10, d/2)
			leftVal, rightVal := s.val/split, s.val%split

			s.val = leftVal
			rightStone := &stone{val: rightVal}
			if s.right != nil {
				rightStone.right = s.right
			}
			s.right = rightStone
			return true
		} else {
			s.val *= 2024
			return false
		}
	}
}

type stoneChain struct {
	stones *stone
	count  int
}

func (sc *stoneChain) step() {
	stones := make([]*stone, 0, sc.count)
	for s := sc.stones; s != nil; s = s.right {
		stones = append(stones, s)
	}
	for _, s := range stones {
		if s.step() {
			sc.count++
		}
	}
}

func (sc stoneChain) print() {
	var b strings.Builder
	for s := sc.stones; s != nil; s = s.right {
		b.WriteString(fmt.Sprintf("%d ", s.val))
	}
	fmt.Println(b.String())
}

func parse(line string) *stoneChain {
	ints := util.ParseLineInts(line)
	stones := make([]*stone, len(ints))
	for i, val := range ints {
		stones[i] = &stone{val: val}
	}
	for i := range stones[1:] {
		stones[i].right = stones[i+1]
	}
	return &stoneChain{stones[0], len(stones)}
}

var STEPS = 25

func solve(input []string) (output string) {
	sc := parse(input[0])
	for i := 0; i < STEPS; i++ {
		sc.step()
	}
	return strconv.Itoa(sc.count)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
