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

type tile struct {
	Id   int
	Size int
	data string
}

func NewTile(id int, data string) *tile {
	t := &tile{Id: id, Size: 10, data: data}

	return t
}

type tiles map[int]*tile

func readTiles(input []string) tiles {
	var id int
	var data string
	ts := make(tiles)

	for _, line := range input {
		if line == "" {
			ts[id] = NewTile(id, data)
			data = ""
		} else if line[:4] == "Tile" {
			if _, err := fmt.Sscanf(line, "Tile %d:", &id); err != nil {
				panic(err)
			}
			continue
		} else {
			data += line
		}
	}
	if data != "" {
		ts[id] = NewTile(id, data)
	}

	return ts
}

type tileEdges struct {
	Id    int
	Edges []int
}

func (te *tileEdges) String() string {
	return fmt.Sprintf("[%d: %s]", te.Id, te.Edges)
}

func (t *tile) Edges() *tileEdges {
	e := make([]int, 4)

	for i := 0; i < t.Size; i++ {
		significance := 1 << i

		topPos := i
		rightPos := (i+1)*t.Size - 1
		bottomPos := t.Size*(t.Size-1) + i
		leftPos := i * t.Size

		pos := []int{topPos, rightPos, bottomPos, leftPos}
		for j, p := range pos {
			if t.data[p] == '#' {
				e[j] += significance
			}
		}
	}

	return &tileEdges{t.Id, e}
}

func (te *tileEdges) Flip() *tileEdges {
	flipped := make([]int, 4)
	flipped[0] = te.Edges[3]
	flipped[1] = te.Edges[2]
	flipped[2] = te.Edges[1]
	flipped[3] = te.Edges[0]

	return &tileEdges{te.Id, flipped}
}

func (te *tileEdges) Reverse() *tileEdges {
	reversed := make([]int, 4)
	for i, e := range te.Edges {
		for j := 0; j < 10; j++ {
			reversed[i] <<= 1
			if (e & 1) == 1 {
				reversed[i] |= 1
			}
			e >>= 1
		}
	}

	return &tileEdges{te.Id, reversed}
}

func (te *tileEdges) Rotate(rotate int) *tileEdges {
	if rotate == 0 {
		return te
	}

	reversed := te.Reverse().Edges
	rotated := make([]int, 4)
	switch rotate {
	case 1:
		rotated[0] = reversed[3]
		rotated[1] = te.Edges[0]
		rotated[2] = reversed[1]
		rotated[3] = te.Edges[2]
	case 2:
		rotated[0] = reversed[2]
		rotated[1] = reversed[3]
		rotated[2] = reversed[0]
		rotated[3] = reversed[1]
	case 3:
		rotated[0] = te.Edges[1]
		rotated[1] = reversed[2]
		rotated[2] = te.Edges[3]
		rotated[3] = reversed[0]
	}

	return &tileEdges{te.Id, rotated}
}

func listTileEdges(ts tiles) []*tileEdges {
	tes := make([]*tileEdges, 0, len(ts)*8)
	for _, t := range ts {
		edges := t.Edges()
		tes = append(tes, edges)
		tes = append(tes, edges.Rotate(1))
		tes = append(tes, edges.Rotate(2))
		tes = append(tes, edges.Rotate(3))
		flipped := edges.Flip()
		tes = append(tes, flipped)
		tes = append(tes, flipped.Rotate(1))
		tes = append(tes, flipped.Rotate(2))
		tes = append(tes, flipped.Rotate(3))
	}

	return tes
}

type candidates map[*tileEdges]bool
type posMap map[int]candidates
type edgeMap map[int]posMap

func buildEdgeMap(tes []*tileEdges) edgeMap {
	em := make(edgeMap)
	for _, te := range tes {
		for side, edge := range te.Edges {
			if _, ok := em[edge]; !ok {
				em[edge] = make(posMap)
			}
			if _, ok := em[edge][side]; !ok {
				em[edge][side] = make(candidates)
			}
			em[edge][side][te] = true
		}
	}

	return em
}

func calculateGridsize(ts tiles) int {
	var gridSize int
	for gridSize = 1; gridSize <= 20; gridSize++ {
		if len(ts) == gridSize*gridSize {
			break
		}
	}

	return gridSize
}

type row []*tileEdges
type solution []row

func findSolution(ts tiles, gridSize int, tes []*tileEdges, em edgeMap) solution {

	s := make(solution, gridSize)
	for i := range s {
		s[i] = make(row, gridSize)
	}

	used := make(map[int]bool)
	for t := range ts {
		used[t] = false
	}
	// fmt.Println(used)

	for _, first := range tes {
		used[first.Id] = true
		s[0][0] = first
		// fmt.Printf("\nTrying first tile: %q\n", first)
		if success := fitTile(s, 0, 1, used, em); success {
			return s
		}
		used[first.Id] = false
	}

	return nil
}

func fitTile(s solution, row, column int, used map[int]bool, em edgeMap) bool {
	if row == len(s) {
		return true
	}

	// fmt.Printf("Fitting tile [%d, %d]\n", row, column)
	checkLeft := column != 0
	checkTop := row != 0

	var ok bool
	var leftMatches, topMatches candidates
	if checkLeft {
		left := s[row][column-1].Edges[1]
		if leftMatches, ok = em[left][3]; !ok {
			return false
		}
		// fmt.Printf("Fits left %d: %s\n", left, leftMatches)
	}
	if checkTop {
		top := s[row-1][column].Edges[2]
		if topMatches, ok = em[top][0]; !ok {
			return false
		}
		// fmt.Printf("Fits top %d: %s\n", top, topMatches)
	}

	filtered := make(candidates)
	if checkLeft && checkTop {
		// fmt.Printf("Left & top: %s, %s\n", leftMatches, topMatches)
		for lm := range leftMatches {
			if used[lm.Id] {
				continue
			}

			if _, ok = topMatches[lm]; ok {
				filtered[lm] = true
			}
		}
	} else if checkLeft {
		// fmt.Printf("Left: %s\n", leftMatches)
		for lm := range leftMatches {
			if !used[lm.Id] {
				filtered[lm] = true
			}
		}
	} else if checkTop {
		// fmt.Printf("Top: %s\n", topMatches)
		for tm := range topMatches {
			if !used[tm.Id] {
				filtered[tm] = true
			}
		}
	}
	// fmt.Printf("Filtered: %s\n", filtered)

	nextColumn := column + 1
	nextRow := row
	if nextColumn == len(s) {
		nextColumn = 0
		nextRow++
	}
	for te := range filtered {
		// fmt.Printf("Trying #%d for [%d, %d]\n", te.Id, row, column)
		used[te.Id] = true
		s[row][column] = te
		if success := fitTile(s, nextRow, nextColumn, used, em); success {
			return true
		}
		used[te.Id] = false
	}

	return false
}

func printCorners(s solution) {
	dim := len(s) - 1
	product := 1
	product *= s[0][0].Id
	product *= s[0][dim].Id
	product *= s[dim][0].Id
	product *= s[dim][dim].Id
	fmt.Println(product)
}

func main() {
	input := readInput()

	ts := readTiles(input)
	// fmt.Println(ts[2311])

	tes := listTileEdges(ts)
	// for _, te := range tes {
	// 	if te.Id == 2729 {
	// 		fmt.Println(te)
	// 	}
	// }

	em := buildEdgeMap(tes)
	// fmt.Println(em[397][0])

	gridSize := calculateGridsize(ts)

	s := findSolution(ts, gridSize, tes, em)
	if s != nil {
		printCorners(s)
	} else {
		fmt.Println("No solutions found")
	}
}
