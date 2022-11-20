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

func playDecks(ds []deck) ([]deck, int) {
	for ds[0].Cards.Len() != 0 && ds[1].Cards.Len() != 0 {
		card0 := ds[0].Cards.Remove(ds[0].Cards.Front()).(int)
		card1 := ds[1].Cards.Remove(ds[1].Cards.Front()).(int)

		if card0 > card1 {
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
