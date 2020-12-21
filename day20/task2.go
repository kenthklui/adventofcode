package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

type tile struct {
	Id   int
	Size int
	Data string
}

func NewTile(id int, data string) *tile {
	t := &tile{Id: id, Size: 10, Data: data}

	return t
}

func (t *tile) Flip() *tile {
	flipped := make([]byte, len(t.Data))
	for i := 0; i < t.Size; i++ {
		for j := 0; j < t.Size; j++ {
			oldIndex := i*t.Size + j
			newIndex := j*t.Size + i
			flipped[newIndex] = t.Data[oldIndex]
		}
	}

	return &tile{t.Id, t.Size, string(flipped)}
}

func (t *tile) Rotate(rotate int) *tile {
	rotated := make([]byte, len(t.Data))
	for i := 0; i < t.Size; i++ {
		for j := 0; j < t.Size; j++ {
			var newI, newJ int
			switch rotate {
			case 1:
				newI = j
				newJ = t.Size - 1 - i
			case 2:
				newI = t.Size - 1 - i
				newJ = t.Size - 1 - j
			case 3:
				newI = t.Size - 1 - j
				newJ = i
			}
			oldIndex := i*t.Size + j
			newIndex := newI*t.Size + newJ

			rotated[newIndex] = t.Data[oldIndex]
		}
	}

	return &tile{t.Id, t.Size, string(rotated)}
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
	Id      int
	Edges   []int
	Flipped bool
	Rotated int
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
			if t.Data[p] == '#' {
				e[j] += significance
			}
		}
	}

	return &tileEdges{t.Id, e, false, 0}
}

func (te *tileEdges) Flip() *tileEdges {
	flipped := make([]int, 4)
	flipped[0] = te.Edges[3]
	flipped[1] = te.Edges[2]
	flipped[2] = te.Edges[1]
	flipped[3] = te.Edges[0]

	return &tileEdges{te.Id, flipped, !te.Flipped, te.Rotated}
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

	return &tileEdges{te.Id, reversed, te.Flipped, te.Rotated}
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

	return &tileEdges{te.Id, rotated, te.Flipped, (te.Rotated + rotate) % 4}
}

func listTileEdges(ts tiles) []*tileEdges {
	// We need to sort tiles by ID to get deterministic results x_x
	tileIds := make([]int, 0, len(ts))
	for id := range ts {
		tileIds = append(tileIds, id)
	}
	sort.Ints(tileIds)

	tes := make([]*tileEdges, 0, len(ts)*8)
	for _, id := range tileIds {
		t := ts[id]
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

type canvas struct {
	Data  string
	Width int
}

func (c canvas) String() string {
	var b strings.Builder
	for i := 0; i < len(c.Data); i += c.Width {
		fmt.Fprintf(&b, "\n%s", c.Data[i:i+c.Width])
	}

	return b.String()[1:]
}

func (c canvas) Flip() canvas {
	flipped := make([]byte, len(c.Data))
	for i := 0; i < c.Width; i++ {
		for j := 0; j < c.Width; j++ {
			oldIndex := i*c.Width + j
			newIndex := j*c.Width + i
			flipped[newIndex] = c.Data[oldIndex]
		}
	}

	return canvas{string(flipped), c.Width}
}

func (c canvas) Search(monster [][]int) int {
	count := 0
	for i := 0; i < c.Width-1-len(monster); i++ {
		for j := 0; j < c.Width; j++ {
			monsterFound := true
			for k, mRow := range monster {
				for _, l := range mRow {
					column := j + l
					if column > c.Width {
						monsterFound = false
						break
					}
					index := (i+k)*c.Width + column
					if c.Data[index] != '#' {
						monsterFound = false
						break
					}
				}
			}
			if monsterFound {
				count += 1
			}
		}
	}

	return count
}

func stitchCanvas(s solution, ts tiles) canvas {
	data := make([]byte, 8*8*len(ts))
	width := 8 * len(s)
	for i, row := range s {
		for j, te := range row {
			t := ts[te.Id]
			if te.Flipped {
				t = t.Flip()
			}
			if te.Rotated > 0 {
				t = t.Rotate(te.Rotated)
			}
			offset := 11
			for k := 0; k < 8; k++ {
				dataIndex := offset + (k * 10)
				canvasIndex := ((i*8)+k)*width + (j * 8)
				for l := 0; l < 8; l++ {
					data[canvasIndex+l] = t.Data[dataIndex+l]
				}
			}
		}
	}

	return canvas{string(data), width}
}

func getMonster(rotate, canvasWidth int) [][]int {
	monsterStr := "                  # \n#    ##    ##    ###\n #  #  #  #  #  #   "
	monsterLines := strings.Split(monsterStr, "\n")

	switch rotate {
	case 0:
	case 1:
		temp := make([]string, len(monsterLines[0]))
		for i := 0; i < len(monsterLines[0]); i++ {
			var b strings.Builder
			for j := len(monsterLines) - 1; j >= 0; j-- {
				fmt.Fprintf(&b, monsterLines[j][i:i+1])
			}
			temp[i] = b.String()
		}
		monsterLines = temp
	case 2:
		temp := make([]string, len(monsterLines))
		for i := len(monsterLines) - 1; i >= 0; i-- {
			var b strings.Builder
			for j := len(monsterLines[i]) - 1; j >= 0; j-- {
				fmt.Fprintf(&b, monsterLines[i][j:j+1])
			}
			temp[len(monsterLines)-1-i] = b.String()
		}
		monsterLines = temp
	case 3:
		temp := make([]string, len(monsterLines[0]))
		for i := len(monsterLines[0]) - 1; i >= 0; i-- {
			var b strings.Builder
			for j := 0; j < len(monsterLines); j++ {
				fmt.Fprintf(&b, monsterLines[j][i:i+1])
			}
			temp[len(monsterLines[0])-1-i] = b.String()
		}
		monsterLines = temp
	}
	// for _, line := range monsterLines {
	// 	fmt.Printf("%s\n", line)
	// }
	// fmt.Println("")

	monster := make([][]int, len(monsterLines))
	for i, line := range monsterLines {
		monster[i] = make([]int, 0, strings.Count(line, "#"))
		for j, r := range line {
			if r == '#' {
				monster[i] = append(monster[i], j)
			}
		}
	}

	return monster
}

func countMonsters(c canvas) int {
	count := 0
	f := c.Flip()
	for i := 0; i < 4; i++ {
		monster := getMonster(i, c.Width)
		count += c.Search(monster)
		count += f.Search(monster)
	}

	return count
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

	c := stitchCanvas(s, ts)
	// fmt.Println(c)

	count := countMonsters(c)
	// monster has 15 #'s - let's assume no monsters overlap in the ocean

	fmt.Println(strings.Count(c.Data, "#") - 15*count)
}
