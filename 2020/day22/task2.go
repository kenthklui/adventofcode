package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

type deck []int

func (d deck) Score() int {
	score := 0
	for i, card := range d {
		score += (len(d) - i) * card
	}
	return score
}

func checkPlayed(roundHistory map[string]void, decks [2]deck) bool {
	var b strings.Builder
	for _, c := range decks[0] {
		b.WriteRune(rune(c))
	}
	// A game is uniquely identified by the contents in one of the two decks
	/*
		b.WriteRune(rune(255))
		for _, c := range decks[1] {
			b.WriteRune(rune(c))
		}
	*/
	key := b.String()

	_, played := roundHistory[key]
	if !played {
		roundHistory[key] = nul
	}
	return played
}

func playDecks(decks [2]deck) ([2]deck, int) {
	roundHistory := make(map[string]void)
	var roundWinner, card0, card1 int
	for len(decks[0]) > 0 && len(decks[1]) > 0 {
		if checkPlayed(roundHistory, decks) {
			return decks, 0
		}

		card0, decks[0] = decks[0][0], decks[0][1:]
		card1, decks[1] = decks[1][0], decks[1][1:]

		if card0 <= len(decks[0]) && card1 <= len(decks[1]) {
			newDecks := [2]deck{slices.Clone(decks[0][:card0]), slices.Clone(decks[1][:card1])}
			_, roundWinner = playDecks(newDecks)
		} else if card0 > card1 {
			roundWinner = 0
		} else {
			roundWinner = 1
		}

		if roundWinner == 1 {
			card0, card1 = card1, card0
		}
		decks[roundWinner] = append(decks[roundWinner], card0, card1)
	}

	if len(decks[0]) == 0 {
		return decks, 1
	} else {
		return decks, 0
	}
}

func readDecks(input []string) [2]deck {
	var decks [2]deck
	deckNum := 0
	for _, line := range input {
		if strings.Contains(line, "Player") {
			decks[deckNum] = make(deck, 0)
		} else if line == "" {
			deckNum++
		} else {
			if card, err := strconv.Atoi(line); err == nil {
				decks[deckNum] = append(decks[deckNum], card)
			} else {
				panic(err)
			}
		}
	}
	return decks
}

func main() {
	input := util.StdinReadlines()
	decks := readDecks(input)
	endDecks, winner := playDecks(decks)
	fmt.Println(endDecks[winner].Score())
}
