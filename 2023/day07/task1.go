package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type hand struct {
	cards [5]int
	bid   int
}

func (h hand) cardCombos() int {
	groups := make(map[int]int)
	for _, c := range h.cards {
		if _, found := groups[c]; found {
			groups[c]++
		} else {
			groups[c] = 1
		}
	}

	trio, pairs := 0, 0
	for _, count := range groups {
		if count == 5 {
			return 6
		} else if count == 4 {
			return 5
		} else if count == 3 {
			trio++
		} else if count == 2 {
			pairs++
		}
	}

	if trio > 0 {
		if pairs > 0 {
			return 4
		} else {
			return 3
		}
	}
	return pairs
}

func cmp(h1, h2 hand) bool {
	if h1.cardCombos() < h2.cardCombos() {
		return true
	} else if h1.cardCombos() > h2.cardCombos() {
		return false
	}

	for i := range h1.cards {
		if h1.cards[i] < h2.cards[i] {
			return true
		} else if h1.cards[i] > h2.cards[i] {
			return false
		}
	}

	// Identical hands...?
	return false
}

type hands []hand

func (hs hands) Len() int           { return len(hs) }
func (hs hands) Less(i, j int) bool { return cmp(hs[i], hs[j]) }
func (hs hands) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }

func (hs hands) winnings() int {
	sort.Sort(hs)
	sum := 0
	for i, h := range hs {
		sum += (i + 1) * h.bid
	}
	return sum
}

func parseCard(b byte) int {
	switch b {
	case 'T':
		return 9
	case 'J':
		return 10
	case 'Q':
		return 11
	case 'K':
		return 12
	case 'A':
		return 13
	default: // r >= '2' && r <= '9'
		return int(b - '1')
	}
}

func parseHands(input []string) hands {
	var err error

	hs := make(hands, len(input))
	for i, line := range input {
		before, after, _ := strings.Cut(line, " ")

		for j := range hs[i].cards {
			hs[i].cards[j] = parseCard(before[j])
		}
		if hs[i].bid, err = strconv.Atoi(after); err != nil {
			panic(err)
		}
	}

	return hs
}

func main() {
	input := util.StdinReadlines()
	hs := parseHands(input)
	fmt.Println(hs.winnings())
}
