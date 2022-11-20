package main

import (
	"bufio"
	"fmt"
	"os"
)

var gridWidth = 32768

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

func parseDirection(dir string) int {
	switch dir {
	case "e":
		return 1
	case "w":
		return -1
	case "se":
		return gridWidth
	case "sw":
		return gridWidth - 1
	case "ne":
		return -gridWidth + 1
	case "nw":
		return -gridWidth
	}

	return 0
}

// false -> white, true -> black
func setTiles(input []string) map[int]bool {
	tiles := make(map[int]bool)

	var dir string
	for _, line := range input {
		key := 0
		for _, b := range line {
			dir += string(b)

			switch b {
			case 'e':
				key += parseDirection(dir)
				dir = ""
			case 'w':
				key += parseDirection(dir)
				dir = ""
			}
		}

		if _, ok := tiles[key]; ok {
			tiles[key] = !tiles[key]
		} else {
			tiles[key] = true
		}
	}

	return tiles
}

func neighbors(key int) []int {
	return []int{
		key - 32768,
		key - 32767,
		key - 1,
		key + 1,
		key + 32767,
		key + 32768,
	}
}

func flipTiles(tiles map[int]bool) {
	flippedTiles := make(map[int]bool)

	for key := range tiles {
		neighborhood := append(neighbors(key), key)
		for _, newKey := range neighborhood {
			// Skip tiles already processed
			if _, ok := flippedTiles[newKey]; ok {
				continue
			}

			bn := blackNeighbors(newKey, tiles)

			if black, ok := tiles[newKey]; ok && black { // black
				if bn == 0 || bn > 2 {
					flippedTiles[newKey] = false
				}
			} else { // white
				if bn == 2 {
					flippedTiles[newKey] = true
				}
			}
		}
	}

	for key, flip := range flippedTiles {
		tiles[key] = flip
	}
}

func blackNeighbors(key int, tiles map[int]bool) int {
	count := 0
	for _, neighbor := range neighbors(key) {
		if tile, ok := tiles[neighbor]; ok && tile {
			count++
		}
	}

	return count
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
		flipTiles(tiles)
	}
	fmt.Printf("%d\n", countBlack(tiles))
}
