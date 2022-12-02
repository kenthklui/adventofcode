package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// List of location IDs, in order:
// 00-06 : Hallway left to right, skipping just outside doors
// 07-08 : Room A, top to bottom
// 09-10 : Room B, top to bottom
// 11-12 : Room C, top to bottom
// 13-14 : Room D, top to bottom

// Need this struct for hashing in a map
type burrowState struct {
	hallway1, hallway2, hallway3, hallway4, hallway5, hallway6, hallway7 uint8
	roomA1, roomA2, roomB1, roomB2, roomC1, roomC2, roomD1, roomD2       uint8
}

// Actual burrow object for performing moves
type burrow struct {
	b    []uint8
	cost int

	// fields for priority queue
	costFinalized bool
	index         int
}

func NewBurrow(size int) *burrow {
	return &burrow{
		b:             make([]uint8, size),
		cost:          0,
		costFinalized: false,
		index:         -1,
	}
}

func (b *burrow) placeAmphipod(index int, color uint8) {
	b.b[index] = color
}

func (b *burrow) state() burrowState {
	return burrowState{
		b.b[0], b.b[1], b.b[2], b.b[3], b.b[4], b.b[5], b.b[6],
		b.b[7], b.b[8], b.b[9], b.b[10], b.b[11], b.b[12], b.b[13], b.b[14],
	}
}

func (b *burrow) done() bool {
	return b.b[7] == 1 && b.b[8] == 1 && b.b[9] == 2 && b.b[10] == 2 &&
		b.b[11] == 3 && b.b[12] == 3 && b.b[13] == 4 && b.b[14] == 4
}

type move struct {
	from, to, cost int
}

func (b *burrow) candidateMoves(pm *pathMap) []move {
	moves := make([]move, 0)
	for i, color := range b.b {
		if color != 0 { // occupied
			moves = append(moves, b.desiredMoves(i, color, pm)...)
		}
	}

	return moves
}

func (b *burrow) take(m move) *burrow {
	nbb := make([]uint8, len(b.b))
	copy(nbb, b.b)
	nbb[m.to] = b.b[m.from]
	nbb[m.from] = 0

	nb := burrow{
		b:             nbb,
		cost:          b.cost + m.cost,
		costFinalized: false,
		index:         -1,
	}

	return &nb
}

func costPerStep(color uint8) int {
	switch color {
	case 1:
		return 1
	case 2:
		return 10
	case 3:
		return 100
	case 4:
		return 1000
	default:
		panic("Invalid color")
	}
}

func (b *burrow) pathClear(p path) bool {
	for _, i := range p {
		if b.b[i] != 0 {
			return false
		}
	}

	return true
}

func (b *burrow) desiredMoves(locId int, color uint8, pm *pathMap) []move {
	moves := make([]move, 0, 6) // I don't think you can have more than 6?
	hallwayLength := 7

	// Check if amphipod can go directly to desginated room
	// Hack for getting room ID per color
	roomId1 := (int(color) + 3) * 2
	roomId2 := roomId1 - 1

	// Amphipod cannot visit room occupied by different colored member
	if b.b[roomId1] > 0 { // Bottom room not clear
		if b.b[roomId1] == color { // Try top room if same color
			if p := pm.path(locId, roomId2); p != nil && b.pathClear(p) {
				cost := costPerStep(color) * pm.dist(locId, roomId2)
				moves = append(moves, move{locId, roomId2, cost})
			}
		}
	} else { // Bottom room is clear
		if p := pm.path(locId, roomId1); p != nil && b.pathClear(p) {
			cost := costPerStep(color) * pm.dist(locId, roomId1)
			moves = append(moves, move{locId, roomId1, cost})
		}
		// No reason to move to top room if bottom room is open
	}

	// If in a room, amphipod can go to hallway
	if locId >= hallwayLength {
		for hallwayId := range b.b[:hallwayLength] {
			if p := pm.path(locId, hallwayId); p != nil {
				if b.pathClear(p) {
					cost := costPerStep(color) * pm.dist(locId, hallwayId)
					moves = append(moves, move{locId, hallwayId, cost})
				}
			}
		}
	}

	return moves
}

// Tools to compute travel costs and blockers for moving between each position

func intMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}
func absDiff(a, b int) int {
	min, max := intMinMax(a, b)
	return max - min
}

type loc struct {
	x, y int
}

func manhattan(l1, l2 loc) int { return absDiff(l1.x, l2.x) + absDiff(l1.y, l2.y) }

type locIdMap map[int]map[int]int

func getLocIdMap(locations []loc) locIdMap {
	locToId := make(locIdMap)

	var xMin, xMax int
	for _, l := range locations {
		xMin, _ = intMinMax(xMin, l.x)
		_, xMax = intMinMax(xMax, l.x)
	}
	for i := xMin; i <= xMax; i++ {
		locToId[i] = make(map[int]int)
	}
	for i, l := range locations {
		locToId[l.x][l.y] = i
	}

	return locToId
}

type path []int
type pathMap struct {
	blockers [][]path
	dists    [][]int
}

func NewPathMap(locations []loc, locToId locIdMap) *pathMap {
	blockers := make([][]path, len(locations))
	for i := range blockers {
		blockers[i] = make([]path, len(locations))
	}
	dists := make([][]int, len(locations))
	for i := range dists {
		dists[i] = make([]int, len(locations))
	}

	// Compute path between every location. Valid/chosen paths will depend on amphipod
	for i, from := range locations {
		for j, to := range locations {
			if i == j {
				continue
			}
			blockers[i][j], dists[i][j] = pathToDest(from, to, locToId)
		}
	}

	return &pathMap{blockers, dists}
}

func pathToDest(from, to loc, locToId map[int]map[int]int) ([]int, int) {
	if from == to {
		return nil, 0
	}

	steps := 0
	blockers := make([]int, 0)

	// Perhaps it's faster to check from destination backward?
	x, y := from.x, from.y
	xInc := 1
	if x > to.x {
		xInc = -1
	}

	for y != 1 {
		y--
		steps++
		if id, ok := locToId[x][y]; ok {
			blockers = append(blockers, id)
		}
	}
	for x != to.x {
		x += xInc
		steps++
		if id, ok := locToId[x][y]; ok {
			blockers = append(blockers, id)
		}
	}
	for y < to.y {
		y++
		steps++
		if id, ok := locToId[x][y]; ok {
			blockers = append(blockers, id)
		}
	}

	return blockers, steps
}

func (pm pathMap) path(from, to int) path   { return pm.blockers[from][to] }
func (pm pathMap) pathList(from int) []path { return pm.blockers[from] }
func (pm pathMap) dist(from, to int) int    { return pm.dists[from][to] }

// Solution searcher

type burrowQueue []*burrow

func (bq burrowQueue) Len() int           { return len(bq) }
func (bq burrowQueue) Less(i, j int) bool { return bq[i].cost < bq[j].cost }
func (bq burrowQueue) Swap(i, j int) {
	bq[i], bq[j] = bq[j], bq[i]
	bq[i].index, bq[j].index = i, j
}
func (bq *burrowQueue) Push(x any) {
	n := len(*bq)
	b := x.(*burrow)
	b.index = n
	*bq = append(*bq, b)
}
func (bq *burrowQueue) Pop() any {
	old := *bq
	n := len(old)
	b := old[n-1]
	old[n-1] = nil
	b.index = -1
	*bq = old[:n-1]

	return b
}

func solve(b *burrow, pm *pathMap) int {
	bs := make(map[burrowState]*burrow)
	bs[b.state()] = b

	bq := make(burrowQueue, 0)
	heap.Push(&bq, b)

	checked := 0
	for len(bq) > 0 {
		curr := heap.Pop(&bq).(*burrow)
		if curr.done() {
			return curr.cost
		}
		curr.costFinalized = true

		checked++

		for _, move := range curr.candidateMoves(pm) {
			nb := curr.take(move)
			newState := nb.state()

			cached, ok := bs[newState]
			if ok {
				if cached.costFinalized {
					continue
				}

				if nb.cost < cached.cost {
					cached.cost = nb.cost
					heap.Fix(&bq, cached.index)
				}
			} else {
				bs[newState] = nb
				heap.Push(&bq, nb)
			}
		}
	}

	fmt.Printf("Failed after checking %d nodes\n", checked)
	return -1
}

// Reading and parsing input

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

func parseInput(input []string, locations []loc, locToId locIdMap) *burrow {
	b := NewBurrow(len(locations))

	for i, line := range input {
		for j, r := range line {
			switch r {
			case 'A':
				fallthrough
			case 'B':
				fallthrough
			case 'C':
				fallthrough
			case 'D':
				color := uint8(r - '@')
				id := locToId[j][i]
				b.placeAmphipod(id, color)
			case '#':
				fallthrough
			case ' ':
				fallthrough
			case '.':
				continue
			default:
				err := fmt.Errorf("Unknown character %q in input", r)
				panic(err)
			}
		}
	}

	return b
}

// main

func main() {
	locations := []loc{
		// Hallway left to right
		{1, 1}, {2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1}, {11, 1},
		// Rooms
		{3, 2}, {3, 3}, {5, 2}, {5, 3}, {7, 2}, {7, 3}, {9, 2}, {9, 3},
	}
	locIdMap := getLocIdMap(locations)
	pm := NewPathMap(locations, locIdMap)

	input := readInput()
	b := parseInput(input, locations, locIdMap)

	solution := solve(b, pm)
	fmt.Println(solution)
}
