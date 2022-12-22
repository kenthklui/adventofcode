package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type player struct {
	row, col int
	facing   [2]int
}

var right = [2]int{0, 1}
var left = [2]int{0, -1}
var up = [2]int{-1, 0}
var down = [2]int{1, 0}

func (p *player) turn(dir byte) {
	switch dir {
	case 'R':
		switch p.facing {
		case up:
			p.facing = right
		case down:
			p.facing = left
		case right:
			p.facing = down
		case left:
			p.facing = up
		}
	case 'L':
		switch p.facing {
		case up:
			p.facing = left
		case down:
			p.facing = right
		case right:
			p.facing = up
		case left:
			p.facing = down
		}
	}
}

func (p *player) turnScore() int {
	switch p.facing {
	case up:
		return 3
	case down:
		return 1
	case right:
		return 0
	case left:
		return 2
	default:
		panic("Invalid direction")
	}
}

func (p *player) password() int {
	return 1000*(p.row+1) + 4*(p.col+1) + p.turnScore()
}

type board struct {
	height, width int
	Map           [][]int
	p             *player
}

func (b *board) run(moves []int, turns []byte) {
	for i, m := range moves {
		b.movePlayer(m)
		if i < len(turns) {
			b.p.turn(turns[i])
		} else {
			break
		}
	}
}

func (b *board) movePlayer(square int) {
	for i := 0; i < square; i++ {
		newRow := b.p.row + b.p.facing[0]
		for {
			if newRow < 0 {
				newRow += b.height
			} else if newRow >= b.height {
				newRow -= b.height
			}
			if b.Map[newRow][b.p.col] != 0 {
				break
			}
			newRow += b.p.facing[0]
		}
		if b.Map[newRow][b.p.col] == 2 { // hit a wall
			break
		} else {
			b.p.row = newRow
		}

		newCol := b.p.col + b.p.facing[1]
		for {
			if newCol < 0 {
				newCol += b.width
			} else if newCol >= b.width {
				newCol -= b.width
			}
			if b.Map[b.p.row][newCol] != 0 {
				break
			}
			newCol += b.p.facing[1]
		}
		if b.Map[b.p.row][newCol] == 2 { // hit a wall
			break
		} else {
			b.p.col = newCol
		}
	}
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

func parseInput(input []string) (*board, []int, []byte) {
	var line string
	var lineNum, boardWidth int
	for lineNum, line = range input {
		if len(line) > boardWidth {
			boardWidth = len(line)
		}
		if line == "" {
			break
		}
	}
	boardHeight := lineNum

	b := new(board)
	b.height = boardHeight
	b.width = boardWidth
	b.Map = make([][]int, 0, boardHeight)
	for lineNum, line = range input[:boardHeight] {
		b.Map = append(b.Map, make([]int, boardWidth))
		for colNum, r := range line {
			switch r {
			// case ' ': // Do nothing, leave as 0
			case '.':
				b.Map[lineNum][colNum] = 1
			case '#':
				b.Map[lineNum][colNum] = 2
			}
		}
	}
	for colNum, r := range b.Map[0] {
		if r == 1 {
			b.p = &player{row: 0, col: colNum, facing: right}
			break
		}
	}

	instruction := input[boardHeight+1]
	moves, turns := make([]int, 0), make([]byte, 0)
	for len(instruction) > 0 {
		nextTurn := strings.IndexAny(instruction, "LR")
		if nextTurn == -1 {
			if tiles, err := strconv.Atoi(instruction); err == nil {
				moves = append(moves, tiles)
				break
			} else {
				panic(err)
			}
		} else {
			if tiles, err := strconv.Atoi(instruction[:nextTurn]); err == nil {
				moves = append(moves, tiles)
			} else {
				panic(err)
			}
			turns = append(turns, instruction[nextTurn])
			instruction = instruction[nextTurn+1:]
		}
	}

	return b, moves, turns
}

func main() {
	input := readInput()
	b, moves, turns := parseInput(input)
	b.run(moves, turns)
	fmt.Println(b.p.password())
}
