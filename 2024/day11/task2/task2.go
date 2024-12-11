package main

import (
	"fmt"
	"strconv"

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

type tuple struct {
	left, right int
}

func step(val int) tuple {
	if val == 0 {
		val = 1
		return tuple{1, -1}
	} else if d := digits(val); d%2 == 0 {
		split := pow(10, d/2)
		return tuple{val / split, val % split}
	} else {
		return tuple{val * 2024, -1}
	}
}

type cache map[int]int

type stoneChain struct {
	stones []int
	caches []cache
}

func (sc *stoneChain) countAfter(steps int) int {
	sc.caches = make([]cache, steps)
	for i := range sc.caches {
		sc.caches[i] = make(cache)
	}
	count := 0
	for _, stone := range sc.stones {
		count += sc.countAfterSteps(stone, steps)
	}
	return count
}

func (sc *stoneChain) countAfterSteps(val, stepsRemain int) int {
	if val < 0 {
		return 0
	}
	if stepsRemain == 0 {
		return 1
	}

	count, cached := sc.caches[stepsRemain-1][val]
	if !cached {
		next := step(val)
		count += sc.countAfterSteps(next.left, stepsRemain-1)
		count += sc.countAfterSteps(next.right, stepsRemain-1)
		sc.caches[stepsRemain-1][val] = count
	}
	return count
}

func parse(line string) *stoneChain {
	return &stoneChain{stones: util.ParseLineInts(line)}
}

var STEPS = 75

func solve(input []string) (output string) {
	sc := parse(input[0])
	return strconv.Itoa(sc.countAfter(STEPS))
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
