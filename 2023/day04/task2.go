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

func (c card) matches() int {
	matches := 0
	for num := range c.Numbers {
		if _, ok := c.Winners[num]; ok {
			matches++
		}
	}
	return matches
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

func parseCards(input []string) []card {
	cards := make([]card, len(input))
	for i, line := range input {
		cards[i] = parseCard(line)
	}
	return cards
}

func countCards(cards []card) int {
	cardCounts := make([]int, len(cards))
	for i := range cardCounts {
		cardCounts[i] = 1
	}

	sum := 0
	for i, c := range cards {
		sum += cardCounts[i]
		m := c.matches()
		for j := 1; j <= m; j++ {
			cardCounts[i+j] += cardCounts[i]
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	cards := parseCards(input)
	fmt.Println(countCards(cards))
}
