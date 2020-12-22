package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"container/list"
)

func readInput() []string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

type deck struct {
	Cards *list.List
}

func (d deck) Score() int {
	multiplier := d.Cards.Len()
	score := 0

	for e := d.Cards.Front(); e != nil; e = e.Next() {
		score += multiplier * e.Value.(int)
		multiplier--
	}

	return score
}

func checkPlayed(roundHistory map[string]bool, ds []deck) bool {
	s1 := make([]rune, ds[0].Cards.Len())
	for i, e := 0, ds[0].Cards.Front(); e != nil; e = e.Next() {
		s1[i] = rune(e.Value.(int))
		i++
	}

	s2 := make([]rune, ds[1].Cards.Len())
	for i, e := 0, ds[1].Cards.Front(); e != nil; e = e.Next() {
		s2[i] = rune(e.Value.(int))
		i++
	}

	key := fmt.Sprintf("%s_%s", string(s1), string(s2))
	if _, ok := roundHistory[key]; ok {
		return true
	}

	roundHistory[key] = true
	return false
}

func recursiveCombat(card0, card1 int, ds []deck) int {
	d1 := deck{list.New()}
	for i, e := 0, ds[0].Cards.Front(); i < card0; i++ {
		d1.Cards.PushBack(e.Value)
		e = e.Next()
	}

	d2 := deck{list.New()}
	for i, e := 0, ds[1].Cards.Front(); i < card1; i++ {
		d2.Cards.PushBack(e.Value)
		e = e.Next()
	}

	_, roundWinner := playDecks([]deck{d1, d2})
	return roundWinner
}

func playDecks(ds []deck) ([]deck, int) {
	roundHistory := make(map[string]bool)
	for ds[0].Cards.Len() != 0 && ds[1].Cards.Len() != 0 {
		if checkPlayed(roundHistory, ds) {
			return nil, 0
		}

		card0 := ds[0].Cards.Remove(ds[0].Cards.Front()).(int)
		card1 := ds[1].Cards.Remove(ds[1].Cards.Front()).(int)

		var roundWinner int
		if card0 <= ds[0].Cards.Len() && card1 <= ds[1].Cards.Len() {
			roundWinner = recursiveCombat(card0, card1, ds)
		} else if card0 > card1 {
			roundWinner = 0
		} else {
			roundWinner = 1
		}

		if roundWinner == 0 {
			ds[0].Cards.PushBack(card0)
			ds[0].Cards.PushBack(card1)
		} else {
			ds[1].Cards.PushBack(card1)
			ds[1].Cards.PushBack(card0)
		}
	}

	if ds[0].Cards.Len() == 0 {
		return ds, 1
	} else {
		return ds, 0
	}
}

func readDecks(input []string) []deck {
	ds := make([]deck, 0)
	i := 0
	for _, line := range input {
		if strings.Contains(line, "Player") {
			ds = append(ds, deck{list.New()})
		} else if line == "" {
			i++
		} else {
			card, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			ds[i].Cards.PushBack(card)
		}
	}
	return ds
}

func main() {
	input := readInput()
	ds := readDecks(input)
	ds, winner := playDecks(ds)

	fmt.Printf("Winner is player %d\n", winner+1)
	fmt.Println(ds[winner].Score())
}
