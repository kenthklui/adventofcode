package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

const jokerValue int = -1

type hand struct {
	cards      [5]int
	bid, combo int
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
	if h1.combo < h2.combo {
		return true
	} else if h1.combo > h2.combo {
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

func makeHand(cards [5]int, bid int) hand {
	h := hand{cards, bid, 0}

	jokers := make([]int, 0)
	for i, c := range h.cards {
		if c == jokerValue {
			jokers = append(jokers, i)
		}
	}

	for value := 1; value <= 12; value++ {
		for _, index := range jokers {
			h.cards[index] = value
		}
		combo := h.cardCombos()
		if combo > h.combo {
			h.combo = combo
		}
	}

	for _, index := range jokers {
		h.cards[index] = jokerValue
	}

	return h
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
	case 'J':
		return jokerValue
	case 'T':
		return 9
	case 'Q':
		return 10
	case 'K':
		return 11
	case 'A':
		return 12
	default: // r >= '2' && r <= '9'
		return int(b - '1')
	}
}

func parseHands(input []string) hands {
	hs := make(hands, len(input))
	for i, line := range input {
		before, after, _ := strings.Cut(line, " ")

		var cards [5]int
		for j := range cards {
			cards[j] = parseCard(before[j])
		}

		bid, err := strconv.Atoi(after)
		if err != nil {
			panic(err)
		}

		hs[i] = makeHand(cards, bid)
	}

	return hs
}

func main() {
	input := util.StdinReadlines()
	hs := parseHands(input)
	fmt.Println(hs.winnings())
}
