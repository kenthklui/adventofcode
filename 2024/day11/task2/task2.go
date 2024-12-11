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
	if stepsRemain == 0 {
		return 1
	}

	stepsRemain--
	count, cached := sc.caches[stepsRemain][val]
	if !cached {
		if val == 0 {
			count = sc.countAfterSteps(1, stepsRemain)
		} else if d := digits(val); d%2 == 0 {
			split := pow(10, d/2)
			count = sc.countAfterSteps(val/split, stepsRemain) + sc.countAfterSteps(val%split, stepsRemain)
		} else {
			count = sc.countAfterSteps(val*2024, stepsRemain)
		}
		sc.caches[stepsRemain][val] = count
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
