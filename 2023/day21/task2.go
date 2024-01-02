package main

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

func sumOddsUpTo(q int) int { // Sum 1, 3, 5, ..., n where n <= q
	half := (q - 1) / 2
	return (half + 1) * (half + 1)
}

func sumEvensUpTo(q int) int { // Sum 2, 4, 6, ..., n where n <= q
	half := q / 2
	return half * (half + 1)
}

func sumUpTo(q int) int {
	return q * (q + 1) / 2
}

type vec struct {
	x, y int
}

var (
	up    = vec{0, -1}
	right = vec{1, 0}
	down  = vec{0, 1}
	left  = vec{-1, 0}
	dirs  = [4]vec{up, right, down, left}
)

func (v vec) add(dir vec) vec    { return vec{v.x + dir.x, v.y + dir.y} }
func (v vec) mod2(steps int) int { return (v.x + v.y + steps) % 2 }

type garden struct {
	width, height, evens, odds int
	start                      vec
	grid                       [][]byte
	stepsCache                 [][]*steps
	edges                      []vec
}

type steps struct {
	node     [][]int
	min, max int
}

func (s *steps) stepsTo(v vec) int { return s.node[v.y][v.x] }

func (s *steps) reachables() []vec {
	vecs := make([]vec, 0, len(s.node)*len(s.node[0]))
	for y, row := range s.node {
		for x, distance := range row {
			if distance >= 0 {
				vecs = append(vecs, vec{x, y})
			}
		}
	}
	return vecs
}

func parseGarden(input []string) *garden {
	width, height := len(input[0]), len(input)
	var start vec
	grid, stepsCache := make([][]byte, height), make([][]*steps, height)
	for y, row := range input {
		grid[y], stepsCache[y] = make([]byte, width), make([]*steps, width)
		for x, c := range row {
			if c == 'S' {
				start = vec{x, y}
			}
			grid[y][x] = byte(c)
		}
	}
	g := &garden{
		width:      width,
		height:     height,
		start:      start,
		grid:       grid,
		stepsCache: stepsCache,
	}

	g.precomputeDistances()
	for _, v := range g.nodeSteps(g.start).reachables() {
		if v.mod2(0) == 0 {
			g.evens++
		} else {
			g.odds++
		}
	}

	return g
}

func (g *garden) inBounds(v vec) bool {
	return v.x >= 0 && v.y >= 0 && v.x < g.width && v.y < g.height
}

type agent struct {
	v     vec
	steps int
}

func (g *garden) nodeSteps(start vec) *steps {
	cached := g.stepsCache[start.y][start.x]
	if cached == nil {
		stepsTo, queued := make([][]int, g.height), make([][]bool, g.height)
		for y := range stepsTo {
			stepsTo[y], queued[y] = make([]int, g.width), make([]bool, g.width)
			for x := range stepsTo[y] {
				stepsTo[y][x] = -1
			}
		}

		queue := make([]agent, 0, g.width*g.height)
		queue = append(queue, agent{start, 0})
		queued[start.y][start.x] = true

		// Interestingly, all inputs have open edge nodes...
		for len(queue) > 0 {
			a := queue[0]
			queue = queue[1:]

			stepsTo[a.v.y][a.v.x] = a.steps
			for _, d := range dirs {
				if n := a.v.add(d); g.inBounds(n) {
					if g.grid[n.y][n.x] != '#' && !queued[n.y][n.x] {
						queue = append(queue, agent{n, a.steps + 1})
						queued[n.y][n.x] = true
					}
				}
			}
		}

		min, max := g.width+g.height, 0
		for _, row := range stepsTo {
			if rowMin := slices.Min(row); rowMin < min {
				min = rowMin
			}
			if rowMax := slices.Max(row); rowMax > max {
				max = rowMax
			}
		}

		cached = &steps{stepsTo, min, max}
		g.stepsCache[start.y][start.x] = cached
	}

	return cached
}

func (g *garden) getEdges() []vec {
	if g.edges == nil {
		edgeNodes := make([]vec, 0, (g.width+g.height)*2)
		for i := 0; i < g.width; i++ {
			edgeNodes = append(edgeNodes, vec{i, 0})
		}
		for i := 1; i < g.height; i++ {
			edgeNodes = append(edgeNodes, vec{g.width - 1, i})
		}
		for i := g.width - 2; i >= 0; i-- {
			edgeNodes = append(edgeNodes, vec{i, g.height - 1})
		}
		for i := g.height - 2; i >= 0; i-- {
			edgeNodes = append(edgeNodes, vec{0, i})
		}
		g.edges = edgeNodes
	}

	return g.edges
}

func (g *garden) precomputeDistances() {
	startCh := make(chan vec)
	go func(ch chan<- vec) {
		ch <- g.start
		for _, e := range g.getEdges() {
			ch <- e
		}
		close(ch)
	}(startCh)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(ch <-chan vec) {
			defer wg.Done()
			for start := range ch {
				g.nodeSteps(start)
			}
		}(startCh)
	}
	wg.Wait()
}

// Corners in clockwise, from top left
func (g *garden) corners() [4]vec {
	w, h := g.width-1, g.height-1
	return [4]vec{vec{0, 0}, vec{w, 0}, vec{w, h}, vec{0, h}}
}

func (g *garden) dirEdge(dirIndex int) []vec {
	var start, length int
	switch dirIndex {
	case 0:
		start = 0
		length = g.width
	case 1:
		start = g.width - 1
		length = g.height
	case 2:
		start = g.width + g.height - 2
		length = g.width
	case 3:
		start = g.width + g.width + g.height - 3
		length = g.height
	}
	return slices.Clone(g.getEdges()[start : start+length])
}

func (g *garden) diagonalSum(cornerIndex, stepsRemain int) int {
	if stepsRemain < 0 {
		return 0
	}

	corner := g.corners()[cornerIndex]
	cornerSteps := g.nodeSteps(corner)
	if cornerSteps.max > g.width+g.height-2 {
		panic("Non-trivial square!")
	}
	oppositeCorner := g.corners()[(cornerIndex+2)%4]
	oppositeCornerDistance := cornerSteps.stepsTo(oppositeCorner)

	reachableSum := 0

	// Assume height == width
	q, r := stepsRemain/g.height, stepsRemain%g.height
	for r <= oppositeCornerDistance {
		perDiagonal := g.singleSum([]agent{{corner, r}})
		diagonalCount := q + 1
		reachableSum += diagonalCount * perDiagonal

		q--
		r += g.height
	}
	q++

	tileValue, altValue := g.odds, g.evens
	if stepsRemain%2 == 0 {
		tileValue, altValue = g.evens, g.odds
	}

	if interweave := (g.height%2 == 1); interweave {
		reachableSum += sumOddsUpTo(q) * tileValue
		reachableSum += sumEvensUpTo(q) * altValue
	} else {
		reachableSum += sumUpTo(q) * tileValue
	}

	return reachableSum
}

func (g *garden) borderStepsSig(borderSteps []int, mod int) string {
	minSteps := slices.Min(borderSteps)
	mod = minSteps - (minSteps % mod)

	// Assume input grid has height and width <= 255
	var b strings.Builder
	for _, s := range borderSteps {
		b.WriteByte(byte(s - mod))
	}
	return b.String()
}

func (g *garden) straightSum(borderSteps []int, dirIndex, depth int) int {
	if maxBorderSteps := slices.Max(borderSteps); maxBorderSteps < 0 {
		return 0
	}

	startDir, endDir := (dirIndex+2)%4, dirIndex
	startSide, endSide := g.dirEdge(startDir), g.dirEdge(endDir)
	slices.Reverse(startSide)

	crossCost := make([][]int, len(startSide))
	for i, start := range startSide {
		crossCost[i] = make([]int, len(endSide))
		nodeCost := g.nodeSteps(start)
		for j, end := range endSide {
			crossCost[i][j] = nodeCost.stepsTo(end) + 1
		}
	}

	mod := g.height
	if dirIndex%2 == 1 {
		mod = g.width
	}
	interweave := (mod%2 == 1)

	tileValue, altValue := g.evens, g.odds
	if tileMod2 := startSide[0].mod2(borderSteps[0]); tileMod2 == 1 {
		tileValue, altValue = g.odds, g.evens
	}

	sig := g.borderStepsSig(borderSteps, mod)
	stepsRemain, working := make([]int, len(borderSteps)), make([]int, len(borderSteps))

	reachableSum := 0
	for iteration, increment := 0, 1; true; iteration += increment {
		for i, ccRow := range crossCost {
			for j := range ccRow {
				working[j] = borderSteps[j] - crossCost[j][i]
			}
			stepsRemain[i] = slices.Max(working)
		}

		if minStepsRemain := slices.Min(stepsRemain); minStepsRemain < 0 {
			// Some edges couldn't cross, break
			a1, a2 := []agent{}, []agent{}
			for i, s := range startSide {
				if borderSteps[i] >= 0 {
					a1 = append(a1, agent{s, borderSteps[i]})
				}
				if stepsRemain[i] >= 0 {
					a2 = append(a2, agent{s, stepsRemain[i]})
				}
			}
			reachableSum += g.singleSum(a1)
			reachableSum += g.singleSum(a2)
			break
		}

		if newSig := g.borderStepsSig(stepsRemain, mod); newSig == sig {
			minJumpAmount := slices.Min(borderSteps)
			skipIterations := minJumpAmount / mod
			if skipIterations > 0 {
				if interweave {
					reachableSum += (tileValue + altValue) * (skipIterations / 2)
					if skipIterations%2 == 1 {
						if iteration%2 == 1 {
							reachableSum += altValue
						} else {
							reachableSum += tileValue
						}
					}
				} else {
					reachableSum += tileValue * skipIterations
				}

				skipValue := skipIterations * mod
				for i := range borderSteps {
					borderSteps[i] -= skipValue
				}

				increment = skipIterations
				continue
			}
		} else {
			sig = newSig
		}

		if interweave && iteration%2 == 1 {
			reachableSum += altValue
		} else {
			reachableSum += tileValue
		}

		copy(borderSteps, stepsRemain)
		increment = 1
	}

	return reachableSum
}

func (g *garden) singleSum(agents []agent) int {
	agentSteps := make([]*steps, len(agents))
	for i, a := range agents {
		agentSteps[i] = g.nodeSteps(a.v)
		if a.steps >= agentSteps[i].max {
			if a.v.mod2(a.steps) == 0 {
				return g.evens
			} else {
				return g.odds
			}
		}
	}

	reachables := make(map[vec]void)
	for i, a := range agents {
		mod2 := a.v.mod2(a.steps)
		for _, dest := range agentSteps[i].reachables() {
			if dest.mod2(0) == mod2 && a.steps >= agentSteps[i].stepsTo(dest) {
				reachables[dest] = nul
			}
		}
	}

	return len(reachables)
}

func (g *garden) sumAll(node vec, steps int) int {
	reachableSum := g.singleSum([]agent{{g.start, steps}})

	if startSteps := g.nodeSteps(g.start); startSteps.min < steps {
		// Handle 4 corners
		corners := g.corners()
		for i, corner := range corners {
			stepsToCorner := startSteps.stepsTo(corner)
			if stepsRemain := steps - stepsToCorner - 2; stepsRemain >= 0 {
				oppositeIndex := (i + 2) % 4
				reachableSum += g.diagonalSum(oppositeIndex, stepsRemain)
			}
		}

		// Handle going straight in 4 directions
		for dirIndex := range dirs {
			edge := g.dirEdge(dirIndex)
			borderSteps := make([]int, len(edge))
			for i, e := range edge {
				borderSteps[i] = steps - startSteps.stepsTo(e) - 1
			}
			reachableSum += g.straightSum(borderSteps, dirIndex, 0)
		}
	}

	return reachableSum
}

var maxSteps = []int{26501365}

// var maxSteps = []int{6, 10, 50, 100, 500, 1000, 5000}

func main() {
	input := util.StdinReadlines()
	g := parseGarden(input)
	for _, steps := range maxSteps {
		fmt.Println(g.sumAll(g.start, steps))
	}
}
