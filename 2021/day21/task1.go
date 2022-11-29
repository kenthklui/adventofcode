package main

import (
	"bufio"
	"fmt"
	"os"
)

type game struct {
	p1, p2                                       int
	score1, score2                               int
	nextPlayer, lastDice, dicePerRoll, rollCount int
	diceCap, boardSize, winningScore             int
}

func NewGame(p1, p2 int) *game {
	return &game{
		p1:           p1,
		p2:           p2,
		score1:       0,
		score2:       0,
		nextPlayer:   1,
		lastDice:     0,
		dicePerRoll:  3,
		rollCount:    0,
		diceCap:      100,
		boardSize:    10,
		winningScore: 1000,
	}
}

func (g *game) rollValue() int {
	value := 0
	for i := 0; i < g.dicePerRoll; i++ {
		g.lastDice++
		if g.lastDice > g.diceCap {
			g.lastDice -= g.diceCap
		}

		value += g.lastDice
		value %= g.boardSize // Keep it small for now
	}

	return value
}

func (g *game) roll() (int, int, int) {
	rollvalue := g.rollValue()
	g.rollCount += g.dicePerRoll

	switch g.nextPlayer {
	case 1:
		g.p1 += rollvalue
		if g.p1 > 10 {
			g.p1 -= 10
		}

		g.score1 += g.p1
		if g.score1 >= g.winningScore {
			return 1, g.score2, g.rollCount
		}

		g.nextPlayer = 2
	case 2:
		g.p2 += rollvalue
		if g.p2 > 10 {
			g.p2 -= 10
		}

		g.score2 += g.p2
		if g.score2 >= g.winningScore {
			return 2, g.score1, g.rollCount
		}

		g.nextPlayer = 1
	}

	return 0, 0, g.rollCount
}

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

func parseInput(input []string) (int, int) {
	var p1, p2 int

	if n, err := fmt.Sscanf(input[0], "Player 1 starting position: %d", &p1); err != nil {
		panic(err)
	} else if n != 1 {
		panic("Failed to parse starting position for player 1")
	}

	if n, err := fmt.Sscanf(input[1], "Player 2 starting position: %d", &p2); err != nil {
		panic(err)
	} else if n != 1 {
		panic("Failed to parse starting position for player 2")
	}

	return p1, p2
}

func main() {
	input := readInput()
	p1, p2 := parseInput(input)

	g := NewGame(p1, p2)
	var winningPlayer, losingScore, diceRolls int
	for winningPlayer == 0 {
		winningPlayer, losingScore, diceRolls = g.roll()
	}

	fmt.Println(losingScore * diceRolls)
}
