package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct {
	x, y int
}

func (p *point) String() string { return fmt.Sprintf("[%d, %d]", p.x, p.y) }

func (p *point) equals(p1 *point) bool { return p.x == p1.x && p.y == p1.y }

func (p *point) foldX(x int) {
	if p.x > x {
		p.x = x + x - p.x
	}
}

func (p *point) foldY(y int) {
	if p.y > y {
		p.y = y + y - p.y
	}
}

type points []*point

func (ps points) Len() int { return len(ps) }

func (ps points) Less(i, j int) bool {
	if ps[i].y == ps[j].y {
		return ps[i].x < ps[j].x
	} else {
		return ps[i].y < ps[j].y
	}
}

func (ps points) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

func (ps points) foldX(x int) {
	for i := range ps {
		ps[i].foldX(x)
	}
}
func (ps points) foldY(y int) {
	for i := range ps {
		ps[i].foldY(y)
	}
}
func (ps points) dedup() points {
	sort.Sort(ps)

	currPoint := ps[0]
	index := 1
	for i, p := range ps[1:] {
		if p.equals(currPoint) {
			continue
		} else {
			ps.Swap(index, i+1)

			index++
			currPoint = p
		}
	}

	return ps[:index]
}

func fold(ps points, foldStr string) points {
	var axis string
	var index int

	if _, err := fmt.Sscanf(foldStr, "fold along %1s=%d", &axis, &index); err == nil {
		switch axis {
		case "x":
			ps.foldX(index)
		case "y":
			ps.foldY(index)
		default:
			panic("Invalid axis")
		}
	} else {
		panic(err)
	}

	return ps.dedup()
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

func parseInput(input []string) (points, []string) {
	ps := make(points, 0)
	folds := make([]string, 0)

	var i, x, y int
	var line string

	for i, line = range input {
		if _, err := fmt.Sscanf(line, "%d,%d", &x, &y); err == nil {
			ps = append(ps, &point{x, y})
		} else {
			break
		}
	}

	folds = input[i+1:]

	return ps, folds
}

func main() {
	input := readInput()
	ps, folds := parseInput(input)

	ps = fold(ps, folds[0])

	fmt.Println(len(ps))
}
