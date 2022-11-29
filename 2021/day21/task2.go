package main

import (
	"bufio"
	"fmt"
	"os"
)

type game struct {
	p1, score1 int
	p2, score2 int
	nextPlayer int
	winner     int
}

const boardSize = 10
const winningScore = 21

func NewGame(p1, p2 int) game {
	return game{
		p1:         p1,
		score1:     0,
		p2:         p2,
		score2:     0,
		nextPlayer: 1,
		winner:     0,
	}
}

func NewWonGame(winner int) game {
	return game{
		p1:         0,
		score1:     0,
		p2:         0,
		score2:     0,
		nextPlayer: 0,
		winner:     winner,
	}
}

type reverseGraph map[game](map[game]uint64)

func (g game) buildGraph() reverseGraph {
	// Map to track which game states lead to which, sorta like a directed graph
	forwardGraph := make(map[game][]game)

	// Map to track game states to handle/play
	gameQueue := make([]game, 1)
	// Insert initial game state and begin generating game state graph
	gameQueue[0] = g
	for len(gameQueue) > 0 {
		// Pop tail
		end := len(gameQueue) - 1
		pop := gameQueue[end]
		gameQueue = gameQueue[:end]

		// Skip if already processed before and cached
		if _, ok := forwardGraph[pop]; ok {
			continue
		}

		nextGames := pop.roll()
		forwardGraph[pop] = nextGames
		gameQueue = append(gameQueue, nextGames...)
	}

	// { gameState : { sourceGameState : numberOfPaths, ... }
	rg := make(reverseGraph)
	for gameState, nextStates := range forwardGraph {
		for _, nextState := range nextStates {
			sourceStates, ok := rg[nextState]
			if !ok {
				sourceStates = make(map[game]uint64)
				rg[nextState] = sourceStates
			}

			_, ok = sourceStates[gameState]
			if ok {
				sourceStates[gameState]++
			} else {
				sourceStates[gameState] = 1
			}
		}
	}

	return rg
}

func (g game) countWins(rg reverseGraph, wantedState game, memoizer map[game]uint64) uint64 {
	if wantedState == g {
		return 1
	}

	if memo, ok := memoizer[wantedState]; ok {
		return memo
	}

	wins := uint64(0)
	for sourceState, paths := range rg[wantedState] {
		wins += paths * g.countWins(rg, sourceState, memoizer)
	}

	memoizer[wantedState] = wins
	return wins
}

func (g game) play() map[int]uint64 {
	results := make(map[int]uint64)
	memoizer := make(map[game]uint64)

	rg := g.buildGraph()
	p1win := NewWonGame(1)
	results[1] = g.countWins(rg, p1win, memoizer)
	p2win := NewWonGame(2)
	results[2] = g.countWins(rg, p2win, memoizer)

	return results
}

func (g game) roll() []game {
	nextGames := make([]game, 0)

	switch g.nextPlayer {
	case 0:
		// Game is over, we're done
	case 1:
		for roll1 := 1; roll1 <= 3; roll1++ {
			for roll2 := 1; roll2 <= 3; roll2++ {
				for roll3 := 1; roll3 <= 3; roll3++ {
					rollSum := roll1 + roll2 + roll3

					gx := g
					gx.p1 += rollSum
					if gx.p1 > 10 {
						gx.p1 -= 10
					}
					gx.score1 += gx.p1

					if gx.score1 >= winningScore {
						gx = NewWonGame(1)
					} else {
						gx.nextPlayer = 2
					}

					nextGames = append(nextGames, gx)
				}
			}
		}
	case 2:
		for roll1 := 1; roll1 <= 3; roll1++ {
			for roll2 := 1; roll2 <= 3; roll2++ {
				for roll3 := 1; roll3 <= 3; roll3++ {
					rollSum := roll1 + roll2 + roll3

					gx := g
					gx.p2 += rollSum
					if gx.p2 > 10 {
						gx.p2 -= 10
					}
					gx.score2 += gx.p2

					if gx.score2 >= winningScore {
						gx = NewWonGame(2)
					} else {
						gx.nextPlayer = 1
					}

					nextGames = append(nextGames, gx)
				}
			}
		}
	}

	return nextGames
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
	results := g.play()

	if results[1] > results[2] {
		fmt.Println(results[1])
	} else {
		fmt.Println(results[2])
	}
}
