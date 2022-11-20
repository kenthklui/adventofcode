package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type board struct {
	raw []int

	unmarkedSum    int
	colCompletions []int
	rowCompletions []int
}

func NewBoard(boardNums []int) *board {
	copyNums := make([]int, len(boardNums))
	sum := 0
	for i, n := range boardNums {
		copyNums[i] = n
		sum += n
	}

	return &board{
		raw:            copyNums,
		unmarkedSum:    sum,
		colCompletions: []int{0, 0, 0, 0, 0},
		rowCompletions: []int{0, 0, 0, 0, 0},
	}
}

func (b *board) Score(winningNumber int) int {
	return b.unmarkedSum * winningNumber
}

func (b *board) MarkboardPosition(position int) bool {
	num := b.raw[position]
	b.unmarkedSum -= num

	row := position / 5
	col := position % 5
	b.colCompletions[col]++
	b.rowCompletions[row]++

	return (b.colCompletions[col] == 5 || b.rowCompletions[row] == 5)
}

type boardPosition struct {
	BoardIndex int
	BoardPos   int
}

func parseInput(input []string) ([]int, map[int][]boardPosition, []*board) {
	drawnNums := make([]int, 0)
	positions := make(map[int][]boardPosition)
	for _, s := range strings.Split(input[0], ",") {
		if n, err := strconv.Atoi(s); err == nil {
			drawnNums = append(drawnNums, n)
			positions[n] = make([]boardPosition, 0)
		}
	}

	boardNums := make([]int, 0)
	boards := make([]*board, 0)

	for _, s := range input[2:] {
		if s == "" {
			newBoard := NewBoard(boardNums)
			boards = append(boards, newBoard)

			boardNums = nil
		} else {
			rowSplit := strings.Split(strings.Trim(s, " "), " ")
			for _, s := range rowSplit {
				n, err := strconv.Atoi(s)
				if err == nil {
					if _, ok := positions[n]; ok {
						newPosition := boardPosition{len(boards), len(boardNums)}
						positions[n] = append(positions[n], newPosition)
					}
					boardNums = append(boardNums, n)
				}
			}
		}
	}

	newBoard := NewBoard(boardNums)
	boards = append(boards, newBoard)

	return drawnNums, positions, boards
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

func main() {
	input := readInput()
	drawnNums, positions, boards := parseInput(input)

	for _, n := range drawnNums {
		if boardPos, ok := positions[n]; ok {
			for _, pos := range boardPos {
				b := boards[pos.BoardIndex]
				if b.MarkboardPosition(pos.BoardPos) {
					fmt.Println(b.Score(n))
					return
				}
			}
		}
	}
}
