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
func flipTiles(input []string) map[string]bool {
	tiles := make(map[string]bool)

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

		key := fmt.Sprintf("%d,%d", x, y)
		if _, ok := tiles[key]; ok {
			tiles[key] = !tiles[key]
		} else {
			tiles[key] = true
		}
	}

	return tiles
}

func countBlack(tiles map[string]bool) int {
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
	tiles := flipTiles(input)
	fmt.Println(countBlack(tiles))
}
