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

type game []set

var defaultBag = set{12, 13, 14}

func (g game) Possible() bool {
	for _, s := range g {
		if s.red > defaultBag.red {
			return false
		}
		if s.green > defaultBag.green {
			return false
		}
		if s.blue > defaultBag.blue {
			return false
		}
	}
	return true
}

func parseGame(line string) game {
	_, after, _ := strings.Cut(line, ": ")
	setStrs := strings.Split(after, "; ")
	g := make([]set, len(setStrs))

	for i, setStr := range setStrs {
		for _, s := range strings.Split(setStr, ", ") {
			before, after, _ := strings.Cut(s, " ")
			switch after {
			case "red":
				g[i].red, _ = strconv.Atoi(before)
			case "green":
				g[i].green, _ = strconv.Atoi(before)
			case "blue":
				g[i].blue, _ = strconv.Atoi(before)
			default:
				panic("Not a color")
			}
		}
	}
	return g
}

func parseGames(input []string) []game {
	games := make([]game, len(input))
	for i, line := range input {
		games[i] = parseGame(line)
	}
	return games
}

func main() {
	input := util.StdinReadlines()
	games := parseGames(input)
	sum := 0

	for i, g := range games {
		if g.Possible() {
			id := i + 1
			sum += id
		}
	}
	fmt.Println(sum)
}
