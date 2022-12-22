package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dir int

const right dir = 0
const down dir = 1
const left dir = 2
const up dir = 3
const invalidDir dir = -1

/*
	Handcrafted from looking at a paper cube

oldFaceNum direction -> newFaceNum newDirection
1 up -> 6 right
1 down -> 3 down
1 right -> 2 right
1 left -> 4 right
2 up -> 6 up
2 down -> 3 left
2 right -> 5 left
2 left -> 1 left
3 up -> 1 up
3 down -> 5 down
3 right -> 2 up
3 left -> 4 down
4 up -> 3 right
4 down -> 6 down
4 right -> 5 right
4 left -> 1 right
5 up -> 3 up
5 down -> 6 left
5 right -> 2 left
5 left -> 4 left
6 up -> 4 up
6 down -> 2 down
6 right -> 5 up
6 left -> 1 down
*/
var mapChange = [6][6]dir{
	{invalidDir, right, down, right, invalidDir, right},
	{left, invalidDir, left, invalidDir, left, up},
	{up, up, invalidDir, down, down, invalidDir},
	{right, invalidDir, right, invalidDir, right, down},
	{invalidDir, left, up, left, invalidDir, left},
	{down, down, invalidDir, up, up, invalidDir},
}

type vec2d [2]int

type player struct {
	cell   *cell
	facing dir
}

func (p *player) move() bool {
	var next *cell
	switch p.facing {
	case right:
		next = p.cell.right
	case down:
		next = p.cell.down
	case left:
		next = p.cell.left
	case up:
		next = p.cell.up
	}
	if next.val == 2 { // wall
		return false
	}

	if next.face != p.cell.face {
		p.facing = mapChange[p.cell.face][next.face]
	}
	p.cell = next
	return true
}

func (p *player) password() int {
	return 1000*(p.cell.pos[0]+1) + 4*(p.cell.pos[1]+1) + int(p.facing)
}

type cell struct {
	face                  int
	pos                   vec2d
	val                   int
	right, down, left, up *cell
}

func (c *cell) setDown(downCell *cell) {
	if c.down != nil {
		panic("Already set")
	}
	c.down = downCell
	if downCell.up != nil {
		panic("Already set")
	}
	downCell.up = c
}
func (c *cell) setRight(rightCell *cell) {
	if c.right != nil {
		panic("Already set")
	}
	c.right = rightCell
	if rightCell.left != nil {
		panic("Already set")
	}
	rightCell.left = c
}

type board struct {
	cells map[vec2d]*cell
	p     *player
}

func turn(turnDir byte, facing dir) dir {
	switch turnDir {
	case 'R':
		switch facing {
		case up:
			return right
		case down:
			return left
		case right:
			return down
		case left:
			return up
		}
	case 'L':
		switch facing {
		case up:
			return left
		case down:
			return right
		case right:
			return up
		case left:
			return down
		}
	}
	panic("Invalid turn")
}

func (b *board) run(moves []int, turns []byte) {
	for i, m := range moves {
		fmt.Println("Move", m)
		for j := 0; j < m; j++ {
			if !b.p.move() {
				break
			}
		}
		if i < len(turns) {
			fmt.Printf("Turn %c\n", turns[i])
			b.p.facing = turn(turns[i], b.p.facing)
		} else {
			break
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

func parseCells(input []string, boardWidth int) map[vec2d]*cell {
	// This whole section is specific to my input
	faceSize := boardWidth / 3
	faceCornerRows := []int{0, 0, faceSize, faceSize * 2, faceSize * 2, faceSize * 3}
	faceCornerCols := []int{faceSize, faceSize * 2, faceSize, 0, faceSize, 0}

	var faceCells [6]map[vec2d]*cell
	cells := make(map[vec2d]*cell, boardWidth*boardWidth)
	var row, col, val int

	for faceNum := range faceCornerRows {
		faceCells[faceNum] = make(map[vec2d]*cell)
		cornerRow, cornerCol := faceCornerRows[faceNum], faceCornerCols[faceNum]
		for row := 0; row < faceSize; row++ {
			boardRow := cornerRow + row
			for col := 0; col < faceSize; col++ {
				boardCol := cornerCol + col

				switch input[boardRow][boardCol] {
				case ' ':
					val = 0
				case '.':
					val = 1
				case '#':
					val = 2
				}

				boardPos := [2]int{boardRow, boardCol}
				newCell := cell{
					face: faceNum,
					val:  val,
					pos:  boardPos,
				}
				cells[boardPos] = &newCell

				facePos := [2]int{row, col}
				faceCells[faceNum][facePos] = &newCell
			}
		}
	}

	// Set regular neighbors
	for faceNum := range faceCornerRows {
		cornerRow, cornerCol := faceCornerRows[faceNum], faceCornerCols[faceNum]
		for r := 0; r < faceSize; r++ {
			row = cornerRow + r
			for c := 0; c < faceSize; c++ {
				col = cornerCol + c
				pos := [2]int{row, col}

				sourceCell := cells[pos]
				rightPos, downPos := [2]int{row, col + 1}, [2]int{row + 1, col}
				if c+1 < faceSize {
					sourceCell.setRight(cells[rightPos])
				}
				if r+1 < faceSize {
					sourceCell.setDown(cells[downPos])
				}
			}
		}
	}

	n := faceSize - 1
	for i := 0; i < faceSize; i++ {
		// 0 up -> 5
		faceCells[0][vec2d{0, i}].up = faceCells[5][vec2d{i, 0}]
		faceCells[5][vec2d{i, 0}].left = faceCells[0][vec2d{0, i}]
		// 0 down -> 2
		faceCells[0][vec2d{n, i}].setDown(faceCells[2][vec2d{0, i}])
		// 0 right -> 1
		faceCells[0][vec2d{i, n}].setRight(faceCells[1][vec2d{i, 0}])
		// 0 left -> 3
		faceCells[0][vec2d{i, 0}].left = faceCells[3][vec2d{n - i, 0}]
		faceCells[3][vec2d{n - i, 0}].left = faceCells[0][vec2d{i, 0}]

		// 1 up -> 5
		faceCells[5][vec2d{n, i}].setDown(faceCells[1][vec2d{0, i}])
		// 1 down -> 2
		faceCells[1][vec2d{n, i}].down = faceCells[2][vec2d{i, n}]
		faceCells[2][vec2d{i, n}].right = faceCells[1][vec2d{n, i}]
		// 1 left -> 0: DONE
		// 1 right -> 4
		faceCells[1][vec2d{i, n}].right = faceCells[4][vec2d{n - i, n}]
		faceCells[4][vec2d{n - i, n}].right = faceCells[1][vec2d{i, n}]

		// 2 up -> 0: DONE
		// 2 down -> 4
		faceCells[2][vec2d{n, i}].setDown(faceCells[4][vec2d{0, i}])
		// 2 right -> 1: DONE
		// 2 left -> 3
		faceCells[2][vec2d{i, 0}].left = faceCells[3][vec2d{0, i}]
		faceCells[3][vec2d{0, i}].up = faceCells[2][vec2d{i, 0}]

		// 3 up -> 2: DONE
		// 3 down -> 5
		faceCells[3][vec2d{n, i}].setDown(faceCells[5][vec2d{0, i}])
		// 3 right -> 4
		faceCells[3][vec2d{i, n}].setRight(faceCells[4][vec2d{i, 0}])
		// 3 left -> 0: DONE

		// 4 up -> 2: DONE
		// 4 down -> 5
		faceCells[4][vec2d{n, i}].down = faceCells[5][vec2d{i, n}]
		faceCells[5][vec2d{i, n}].right = faceCells[4][vec2d{n, i}]
		// 4 right -> 1: DONE
		// 4 left -> 3: DONE

		// 6: ALL DONE
		// 12 pairs in total handled
	}

	return cells
}

func parseInput(input []string) (*board, []int, []byte) {
	var boardWidth int
	var lineNum int
	var line string
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
	b.cells = parseCells(input, boardWidth)

	for colNum := 0; colNum < boardWidth; colNum++ {
		if c, ok := b.cells[vec2d{0, colNum}]; ok && c.val == 1 {
			b.p = &player{cell: c, facing: right}
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
	cube, moves, turns := parseInput(input)

	for _, c := range cube.cells {
		if c.up == nil || c.down == nil || c.right == nil || c.left == nil {
			fmt.Println(*c)
		}
	}
	cube.run(moves, turns)
	fmt.Println(cube.p.password())
}
