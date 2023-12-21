package main

import (
	"fmt"
	"slices"

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

var up = vec{0, -1}
var right = vec{1, 0}
var down = vec{0, 1}
var left = vec{-1, 0}
var dirs = [4]vec{up, right, down, left}

func (v vec) add(dir vec) vec    { return vec{v.x + dir.x, v.y + dir.y} }
func (v vec) mod2(steps int) int { return (v.x + v.y + steps) % 2 }

type garden struct {
	width, height, evens, odds int
	start                      vec
	grid                       map[vec]byte

	stepsCache map[vec]steps
	edges      []vec
}

type steps struct {
	node             map[vec]int
	edgeMin, edgeMax int
}

func parseGarden(input []string) garden {
	width, height := len(input[0]), len(input)
	var start vec
	grid := make(map[vec]byte)
	odds, evens := 0, 0
	for y, row := range input {
		for x, c := range row {
			if c == 'S' {
				start = vec{x, y}
				c = '.'
			}
			grid[vec{x, y}] = byte(c)
		}
	}
	g := garden{
		width:      width,
		height:     height,
		evens:      evens,
		odds:       odds,
		start:      start,
		grid:       grid,
		stepsCache: make(map[vec]steps),
	}
	for v := range g.nodeSteps(g.start).node {
		if v.mod2(0) == 0 {
			g.evens++
		} else {
			g.odds++
		}
	}

	return g
}

func (g garden) inBounds(v vec) bool {
	return v.x >= 0 && v.y >= 0 && v.x < g.width && v.y < g.height
}

func (g garden) crossBorder(oob vec) vec {
	if oob.x < 0 {
		oob.x += g.width
		return oob
	} else if oob.x >= g.width {
		oob.x -= g.width
		return oob
	}
	if oob.y < 0 {
		oob.y += g.height
		return oob
	} else if oob.y >= g.height {
		oob.y -= g.height
		return oob
	}
	panic("Not out of bounds")
}

type agent struct {
	v     vec
	steps int
}

func (g garden) nodeSteps(start vec) steps {
	s, ok := g.stepsCache[start]
	if !ok {
		travelled := make(map[vec]int)
		queue := make([]agent, 0, g.width*g.height)
		queue = append(queue, agent{start, 0})

		stepsToNode := make(map[vec]int)
		// Interestingly all inputs have open edge nodes
		for len(queue) > 0 {
			a := queue[0]
			queue = queue[1:]

			stepsToNode[a.v] = a.steps
			travelled[a.v] = 2

			for _, d := range dirs {
				if n := a.v.add(d); g.inBounds(n) {
					if travelled[n] > 0 {
						continue
					}
					if g.grid[n] == '#' {
						continue
					}

					travelled[n] = 1
					queue = append(queue, agent{n, a.steps + 1})
				}
			}
		}

		edgeMin, edgeMax := g.width+g.height, 0
		for _, e := range g.getEdges() {
			s, _ := stepsToNode[e]
			if s < edgeMin {
				edgeMin = s
			}
			if s > edgeMax {
				edgeMax = s
			}
		}

		s = steps{stepsToNode, edgeMin, edgeMax}
		g.stepsCache[start] = s
	}

	return s
}

func (g garden) getEdges() []vec {
	if g.edges == nil {
		edgeNodes := make([]vec, 0, (g.width+g.height-2)*2)
		for i := 0; i < g.width; i++ {
			edgeNodes = append(edgeNodes, vec{i, 0})
		}
		for i := 1; i < g.height; i++ {
			edgeNodes = append(edgeNodes, vec{g.width - 1, i})
		}
		for i := g.width - 2; i >= 0; i-- {
			edgeNodes = append(edgeNodes, vec{i, g.height - 1})
		}
		for i := g.height - 2; i >= 1; i-- {
			edgeNodes = append(edgeNodes, vec{0, i})
		}
		g.edges = edgeNodes
	}

	return g.edges
}

// Corners in clockwise, from top left
func (g garden) corners() []vec {
	w, h := g.width-1, g.height-1
	return []vec{vec{0, 0}, vec{w, 0}, vec{w, h}, vec{0, h}}
}

func (g garden) dirEdge(dirIndex int) []vec {
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
		edge := make([]vec, g.height)
		copy(edge, g.getEdges()[start:])
		edge[len(edge)-1] = g.getEdges()[0]
		return edge
	}
	return slices.Clone(g.getEdges()[start : start+length])
}

func (g garden) diagonalReachable(cornerIndex, stepsRemain int) int {
	if stepsRemain < 0 {
		return 0
	}
	// Assume height == width
	corner := g.corners()[cornerIndex]
	cornerSteps := g.nodeSteps(corner)
	if cornerSteps.edgeMax > g.width+g.height-2 {
		panic("Non-trivial square!")
	}
	oppositeCorner := g.corners()[(cornerIndex+2)%4]
	oppositeCornerDistance := cornerSteps.node[oppositeCorner]

	reachableSum := 0

	// Assume height == width
	q, r := stepsRemain/g.height, stepsRemain%g.height
	for r < oppositeCornerDistance {
		perDiagonal := g.localReachable([]agent{{corner, r}})
		diagonalCount := q + 1

		// fmt.Println(cornerIndex, diagonalCount, perDiagonal)
		reachableSum += diagonalCount * perDiagonal

		q--
		r += g.height
	}
	q++

	initial, alternative := g.odds, g.evens
	if stepsRemain%2 == 0 {
		initial, alternative = g.evens, g.odds
	}

	if interweave := (g.height%2 == 1); interweave {
		reachableSum += sumOddsUpTo(q) * initial
		reachableSum += sumEvensUpTo(q) * alternative
	} else {
		reachableSum += sumUpTo(q) * initial
	}

	return reachableSum
}

func (g garden) straightReachable(agents []agent, dirIndex int) int {
	agentSteps := make(map[vec]steps)

	type edgeTracker struct {
		edge                       vec
		closest, closestAgentIndex int
	}
	edgeNodes := g.dirEdge(dirIndex)
	oppositeEdge := make([]edgeTracker, len(edgeNodes))
	for i, e := range edgeNodes {
		oppositeEdge[i].edge = e
	}

	canCover, canCross := false, false
	for i, a := range agents {
		agentSteps[a.v] = g.nodeSteps(a.v)
		if a.steps >= agentSteps[a.v].edgeMax {
			canCover, canCross = true, true
		}
		for j, e := range oppositeEdge {
			crossing := agentSteps[a.v].node[e.edge]
			if e.closest == 0 || crossing < e.closest {
				oppositeEdge[j].closest = crossing
				oppositeEdge[j].closestAgentIndex = i
			}
			if diff := agents[i].steps - crossing; diff > 0 {
				canCross = true
			}
		}
	}

	newAgentIndices := make(map[int]void)
	for _, e := range oppositeEdge {
		newAgentIndices[e.closestAgentIndex] = nul
	}
	newAgents := make([]agent, 0, len(agents))
	crossings := make([]int, len(newAgentIndices))
	for i := range newAgentIndices {
		a := agents[i]
		oppositeDir := dirs[(dirIndex+2)%4]
		oppositeSide := g.crossBorder(a.v.add(oppositeDir))
		crossing := agentSteps[a.v].node[oppositeSide] + 1
		newSteps := a.steps - crossing
		if newSteps > 0 {
			crossings[len(newAgents)] = crossing
			newAgents = append(newAgents, agent{a.v, newSteps})
		}
	}

	reachableSum := 0
	if canCover {
		// Since all borders are open, any cell should be fine for even/odd computation
		var currentSquareSum int
		if mod2 := agents[0].v.mod2(agents[0].steps); mod2 == 0 {
			currentSquareSum = g.evens
		} else {
			currentSquareSum = g.odds
		}

		canSkip := (len(agents) == len(newAgents))
		for i, na := range newAgents {
			if na.steps < agentSteps[na.v].edgeMax {
				canSkip = false
				break
			}
			if na.v != agents[i].v {
				canSkip = false
				break
			}
		}

		if canSkip {
			stepsRemain := make([]int, len(newAgents))
			for i, a := range newAgents {
				stepsRemain[i] = a.steps
			}
			leastSteps := slices.Min(stepsRemain)
			leastStepsAgentIndex := slices.IndexFunc(newAgents, func(a agent) bool {
				return a.steps == leastSteps
			})
			leastStepsAgent := newAgents[leastStepsAgentIndex]

			h := slices.Min(crossings)
			q, r := leastSteps/h, leastSteps%h
			for q > 0 && r+h < agentSteps[leastStepsAgent.v].edgeMax {
				q--
				r += h
			}
			// fmt.Println(dirIndex, "Shortcutting:", newAgents, h, q, r)
			// fmt.Println(dirIndex, "Agents are:", agents, newAgents)
			// fmt.Println(dirIndex, "Opposite edge:", oppositeEdge)

			if interweave := (g.height%2 == 1); interweave {
				// fmt.Println(dirIndex, "Fast forward", q, g.evens, g.odds, currentSquareSum)
				reachableSum += (g.evens + g.odds) * (q / 2)
				if q%2 == 1 {
					reachableSum += currentSquareSum
				}
			} else {
				reachableSum += q * currentSquareSum
			}

			for i := range newAgents {
				newAgents[i].steps = agents[i].steps - q*h
			}
			newSum := g.straightReachable(newAgents, dirIndex)
			// fmt.Println(dirIndex, newAgents, newSum)
			reachableSum += newSum

			return reachableSum
		}

		if mod2 := agents[0].v.mod2(agents[0].steps); mod2 == 0 {
			reachableSum += g.evens
		} else {
			reachableSum += g.odds
		}
	} else {
		newSum := g.localReachable(agents)
		// fmt.Println(dirIndex, "Sum:", newSum, "Agents are:", agents)
		reachableSum += newSum
	}

	// Build new agent list and recurse
	if canCross {
		reachableSum += g.straightReachable(newAgents, dirIndex)
	}

	return reachableSum
}

func (g garden) localReachable(agents []agent) int {
	agentSteps := make(map[vec]steps)

	for _, a := range agents {
		agentSteps[a.v] = g.nodeSteps(a.v)
		mod2 := a.v.mod2(a.steps)
		if a.steps >= agentSteps[a.v].edgeMax {
			if mod2 == 0 {
				return g.evens
			} else {
				return g.odds
			}
		}
	}

	reachables := make(map[vec]void)
	for _, a := range agents {
		mod2 := a.v.mod2(a.steps)
		for dest, stepsToReach := range agentSteps[a.v].node {
			if a.steps >= stepsToReach && dest.mod2(0) == mod2 {
				reachables[dest] = nul
			}
		}
	}

	return len(reachables)
}

func (g garden) neighborReachable(agents []agent) [4]int {
	var reachables [4]int
	for i, dir := range dirs {
		crossBorderAgents := make([]agent, 0)
		for _, a := range agents {
			n := a.v.add(dir)
			if g.inBounds(n) {
				continue
			}
			n = g.crossBorder(n)
			crossBorderAgents = append(crossBorderAgents, agent{n, a.steps - 1})
		}
		reachables[i] = g.localReachable(crossBorderAgents)
	}
	return reachables
}

func (g garden) multiReachable(node vec, steps int) int {
	reachableSum := 0
	reachableSum += g.localReachable([]agent{{g.start, steps}})

	startSteps := g.nodeSteps(g.start)
	if startSteps.edgeMin < steps {
		// Handle 4 corners
		corners := g.corners()
		for i, corner := range corners {
			stepsToCorner := startSteps.node[corner]
			oppositeIndex := (i + 2) % 4
			sum := g.diagonalReachable(oppositeIndex, steps-stepsToCorner-2)
			// fmt.Println("Diagonal", i, sum)
			reachableSum += sum
		}

		// Handle going straight in 4 directions
		for dirIndex, dir := range dirs {
			edge := g.dirEdge(dirIndex)
			agents := make([]agent, 0, len(edge))
			for _, e := range edge {
				newSteps := steps - startSteps.node[e] - 1
				if newSteps > 0 {
					otherSide := g.crossBorder(e.add(dir))
					agents = append(agents, agent{otherSide, newSteps})
				}
			}

			sum := g.straightReachable(agents, dirIndex)
			// fmt.Println("Straight:", dirIndex, sum)
			reachableSum += sum
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
		fmt.Println(g.multiReachable(g.start, steps))
	}
}
