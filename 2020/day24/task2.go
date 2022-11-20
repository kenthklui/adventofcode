package main

import (
	"bufio"
	"fmt"
	"os"
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

func parseDirection(dir string) (int, int) {
	switch dir {
	case "e":
		return 0, 1
	case "w":
		return 0, -1
	case "se":
		return 1, 0
	case "sw":
		return 1, -1
	case "ne":
		return -1, 1
	case "nw":
		return -1, 0
	}

	return 0, 0
}

// false -> white, true -> black
func setTiles(input []string) map[int]bool {
	tiles := make(map[int]bool)

	var dir string
	var x, y, dx, dy int
	for _, line := range input {
		x, y = 0, 0
		for _, b := range line {
			dir += string(b)

			switch b {
			case 'e':
				dx, dy = parseDirection(dir)
				x += dx
				y += dy
				dir = ""
			case 'w':
				dx, dy = parseDirection(dir)
				x += dx
				y += dy
				dir = ""
			}
		}

		key := x*32768 + y
		if _, ok := tiles[key]; ok {
			tiles[key] = !tiles[key]
		} else {
			tiles[key] = true
		}
	}

	return tiles
}

func flipTiles(tiles map[int]bool) map[int]bool {
	newTiles := make(map[int]bool)

	for key := range tiles {
		neighborhood := []int{
			key - 32768,
			key - 32767,
			key - 1,
			key,
			key + 1,
			key + 32767,
			key + 32768,
		}

		for _, newKey := range neighborhood {
			if _, ok := newTiles[newKey]; ok {
				continue
			}

			blackNeighbors := countBlackNeighbors(newKey, tiles)

			black, ok := tiles[newKey]
			if ok && black { // black
				if blackNeighbors == 0 || blackNeighbors > 2 {
					newTiles[newKey] = false
				} else {
					newTiles[newKey] = true
				}
			} else { // white
				if blackNeighbors == 2 {
					newTiles[newKey] = true
				} else {
					newTiles[newKey] = false
				}
			}
		}
	}

	return newTiles
}

func countBlackNeighbors(key int, tiles map[int]bool) int {
	blackNeighbors := 0
	neighbors := []int{
		key - 32768,
		key - 32767,
		key - 1,
		key + 1,
		key + 32767,
		key + 32768,
	}

	for _, neighbor := range neighbors {
		if tile, ok := tiles[neighbor]; ok && tile {
			blackNeighbors++
		}
	}

	return blackNeighbors
}

func countBlack(tiles map[int]bool) int {
	black := 0
	for _, t := range tiles {
		if t {
			black++
		}
	}

	return black
}

func main() {
	input := readInput()
	tiles := setTiles(input)

	for i := 1; i <= 100; i++ {
		tiles = flipTiles(tiles)
	}
	fmt.Printf("%d\n", countBlack(tiles))
}
