package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

type card struct {
	Winners, Numbers map[int]void
}

func (c card) points() int {
	matches := 0
	for num := range c.Numbers {
		if _, ok := c.Winners[num]; ok {
			matches++
		}
	}
	if matches == 0 {
		return 0
	}
	score := 1
	for i := matches; i > 1; i-- {
		score *= 2
	}
	return score
}

func parseCard(line string) card {
	c := card{make(map[int]void), make(map[int]void)}

	_, str, _ := strings.Cut(line, ": ")
	winStr, numStr, _ := strings.Cut(str, " | ")
	for _, s := range strings.Split(winStr, " ") {
		if s == "" {
			continue
		}
		if n, err := strconv.Atoi(s); err == nil {
			c.Winners[n] = nul
		} else {
			panic(err)
		}
	}
	for _, s := range strings.Split(numStr, " ") {
		if s == "" {
			continue
		}
		if n, err := strconv.Atoi(s); err == nil {
			c.Numbers[n] = nul
		} else {
			panic(err)
		}
	}

	return c
}

func main() {
	input := util.StdinReadlines()
	pointSum := 0
	for _, line := range input {
		pointSum += parseCard(line).points()
	}
	fmt.Println(pointSum)
}
