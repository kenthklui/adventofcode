package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type set struct {
	red, green, blue int
}

type game struct {
	Id   int
	bag  set
	Sets []set
}

func (g *game) MakePossible() {
	for _, s := range g.Sets {
		if s.red > g.bag.red {
			g.bag.red = s.red
		}
		if s.green > g.bag.green {
			g.bag.green = s.green
		}
		if s.blue > g.bag.blue {
			g.bag.blue = s.blue
		}
	}
}

func (g *game) Power() int {
	return g.bag.red * g.bag.green * g.bag.blue
}

func parseGame(line string) *game {
	before, after, _ := strings.Cut(line, ": ")
	setStrs := strings.Split(after, "; ")

	_, a, _ := strings.Cut(before, " ")
	id, _ := strconv.Atoi(a)

	g := game{id, set{}, make([]set, len(setStrs))}
	for i, setStr := range setStrs {
		for _, s := range strings.Split(setStr, ", ") {
			before, after, _ := strings.Cut(s, " ")
			switch after {
			case "red":
				g.Sets[i].red, _ = strconv.Atoi(before)
			case "green":
				g.Sets[i].green, _ = strconv.Atoi(before)
			case "blue":
				g.Sets[i].blue, _ = strconv.Atoi(before)
			default:
				panic("Not a color")
			}
		}
	}
	return &g
}

func parseGames(input []string) []*game {
	games := make([]*game, len(input))
	for i, line := range input {
		games[i] = parseGame(line)
	}
	return games
}

func main() {
	input := util.StdinReadlines()
	games := parseGames(input)
	sum := 0
	for _, g := range games {
		g.MakePossible()
		sum += g.Power()
	}
	fmt.Println(sum)
}
