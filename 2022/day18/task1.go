package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type void struct{}
type grid map[int]plane
type plane map[int]line
type line map[int]void

var empty void

func parseInput(input []string) grid {
	g := make(grid)

	for _, text := range input {
		splits := strings.Split(text, ",")
		x, _ := strconv.Atoi(splits[0])
		y, _ := strconv.Atoi(splits[1])
		z, _ := strconv.Atoi(splits[2])
		if _, ok := g[x]; !ok {
			g[x] = make(plane)
		}
		if _, ok := g[x][y]; !ok {
			g[x][y] = make(line)
		}
		g[x][y][z] = empty
	}

	return g
}

func countSides(g grid) int {
	sides := 0
	for x, p := range g {
		for y, l := range p {
			for z := range l {
				sides += 6
				if _, ok := g[x][y][z-1]; ok {
					sides--
				}
				if _, ok := g[x][y][z+1]; ok {
					sides--
				}
				if _, ok := g[x][y-1][z]; ok {
					sides--
				}
				if _, ok := g[x][y+1][z]; ok {
					sides--
				}
				if _, ok := g[x-1][y][z]; ok {
					sides--
				}
				if _, ok := g[x+1][y][z]; ok {
					sides--
				}
			}
		}
	}

	return sides
}

func main() {
	input := readInput()
	g := parseInput(input)
	fmt.Println(countSides(g))
}
