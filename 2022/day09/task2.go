package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func intMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

type move struct {
	dx, dy int
}

type knot struct {
	x, y int
}

func (k *knot) move(m move) {
	k.x += m.dx
	k.y += m.dy
}

func (k1 *knot) inRange(k2 *knot) bool {
	x1, x2 := intMinMax(k1.x, k2.x)
	y1, y2 := intMinMax(k1.y, k2.y)
	_, maxDiff := intMinMax(x2-x1, y2-y1)
	return maxDiff < 2
}

func (chaser *knot) chase(target *knot) {
	if chaser.inRange(target) {
		return
	}

	if chaser.x > target.x {
		chaser.x--
	} else if chaser.x < target.x {
		chaser.x++
	}

	if chaser.y > target.y {
		chaser.y--
	} else if chaser.y < target.y {
		chaser.y++
	}
}

func moveRope(moves []move, ropeLength int) int {
	rope := make([]*knot, ropeLength)
	for i := range rope {
		rope[i] = &knot{0, 0}
	}

	tailPositions := make(map[knot]byte)
	for _, m := range moves {
		rope[0].move(m)
		for i, knot := range rope[1:] {
			knot.chase(rope[i])
		}

		tail := *rope[ropeLength-1]
		tailPositions[tail] = 0
	}
	return len(tailPositions)
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

func parseInput(input []string) []move {
	moves := make([]move, 0, len(input))
	for _, line := range input {
		tokens := strings.Split(line, " ")

		steps, _ := strconv.Atoi(tokens[1])
		var m move
		switch tokens[0] {
		case "U":
			m.dx, m.dy = 0, -1
		case "D":
			m.dx, m.dy = 0, 1
		case "L":
			m.dx, m.dy = -1, 0
		case "R":
			m.dx, m.dy = 1, 0
		}

		for i := steps; i > 0; i-- {
			moves = append(moves, m)
		}
	}

	return moves
}

func main() {
	input := readInput()
	moves := parseInput(input)
	fmt.Println(moveRope(moves, 10))
}
