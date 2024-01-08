package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type vec struct{ x, y int }

var up = vec{0, -1}
var right = vec{1, 0}
var down = vec{0, 1}
var left = vec{-1, 0}
var dirs = [4]vec{up, right, down, left}

func (v vec) add(d vec) vec { return vec{v.x + d.x, v.y + d.y} }

type area struct {
	start, end    vec
	width, height int
	aMap          [][]byte
}

func parseArea(input []string) area {
	slopes := []string{"^", ">", "v", "<"}
	for i := range input {
		for _, s := range slopes {
			input[i] = strings.ReplaceAll(input[i], s, ".")
		}
	}
	aMap := make([][]byte, len(input))
	for i := range input {
		aMap[i] = []byte(input[i])
	}

	a := area{
		start:  vec{1, 0},
		end:    vec{len(input[0]) - 2, len(input) - 1},
		width:  len(input[0]),
		height: len(input),
		aMap:   aMap,
	}
	a.fillBorder()

	return a
}

func (a area) inBounds(v vec) bool {
	return v.x >= 0 && v.y >= 0 && v.x < a.width && v.y < a.height
}

func (a area) neighbors(v vec) []vec {
	n := make([]vec, 0, len(dirs))
	for _, dir := range dirs {
		next := v.add(dir)
		if !a.inBounds(next) {
			continue
		}
		if a.aMap[next.y][next.x] == '.' {
			n = append(n, next)
		}
	}
	return n
}

func (a area) nextFork(prev, curr vec, stepsAway int) (bool, vec, int) {
	neighbors := a.neighbors(curr)
	if len(neighbors) == 1 { // Dead end
		return false, curr, stepsAway
	} else if len(neighbors) == 2 { // path
		if neighbors[0] != prev {
			return a.nextFork(curr, neighbors[0], stepsAway+1)
		} else {
			return a.nextFork(curr, neighbors[1], stepsAway+1)
		}
	} else { // fork in the road
		return true, curr, stepsAway
	}
}

const borderByte byte = '$'

func (a area) fillBorder() {
	a.recurseFill(vec{0, 0})
	a.recurseFill(vec{2, 0})
}

func (a area) recurseFill(v vec) {
	if a.aMap[v.y][v.x] != '#' {
		return
	}
	a.aMap[v.y][v.x] = borderByte
	for _, d := range dirs {
		if d = v.add(d); a.inBounds(d) {
			a.recurseFill(d)
		}
	}
}

func (a area) print() {
	for _, row := range a.aMap {
		fmt.Println(string(row))
	}
}

func (a area) longestHike() int {
	nm := buildNodeMap(a)
	return nm.longestHike()
}

type node struct {
	v                 vec
	visited, edgeNode bool
	adjacent          []int
	endDist           int
}

func NewNode(v vec) *node {
	return &node{v: v, adjacent: make([]int, 0)}
}

type nodeMap struct {
	nodes           []*node
	nodeDistances   [][]int
	countPerEndDist []int
}

const (
	startIndex = 0
	endIndex   = 1
)

func buildNodeMap(a area) *nodeMap {
	nodes := []*node{NewNode(a.start), NewNode(a.end)}
	for y, row := range a.aMap {
		for x, c := range row {
			if c != '.' {
				continue
			}

			if len(a.neighbors(vec{x, y})) > 2 {
				nodes = append(nodes, NewNode(vec{x, y}))
			}
		}
	}

	nm := &nodeMap{nodes: nodes}
	nm.detectEdgeNodes(a)
	nm.computeNodeDistances(a)
	nm.setEndDistance()

	return nm
}

func (nm *nodeMap) detectEdgeNodes(a area) {
	nm.nodes[startIndex].edgeNode, nm.nodes[endIndex].edgeNode = true, true
	rest := nm.nodes[2:]
	for i, n := range rest {
		for j, d1 := range dirs {
			d1 = n.v.add(d1)
			d2 := dirs[(j+1)%4]
			d2 = d1.add(d2)
			if a.aMap[d1.y][d1.x] == borderByte || a.aMap[d2.y][d2.x] == borderByte {
				rest[i].edgeNode = true
				break
			}
		}
	}
}

func (nm *nodeMap) computeNodeDistances(a area) {
	nodeIndices := make(map[vec]int)
	for i, n := range nm.nodes {
		nodeIndices[n.v] = i
	}

	nodeDistances := make([][]int, len(nm.nodes))
	for i := range nodeDistances {
		nodeDistances[i] = make([]int, len(nm.nodes))
	}
	for sourceIndex, sourceNode := range nm.nodes {
		var destIndex int
		for _, neighbor := range a.neighbors(sourceNode.v) {
			isFork, dest, stepsAway := a.nextFork(sourceNode.v, neighbor, 1)
			if isFork {
				destIndex = nodeIndices[dest]
			} else if dest == a.start {
				destIndex = startIndex
			} else if dest == a.end {
				destIndex = endIndex
			} else {
				continue // ignore dead ends
			}

			// Prevent loopbacks
			if sourceIndex == destIndex {
				continue
			}

			// Prevent adding the same adjacent node multiple times
			if nodeDistances[sourceIndex][destIndex] == 0 {
				sourceNode.adjacent = append(sourceNode.adjacent, destIndex)
				destNode := nm.nodes[destIndex]
				destNode.adjacent = append(destNode.adjacent, sourceIndex)
			}

			// Only store the max distance between two nodes
			if stepsAway > nodeDistances[sourceIndex][destIndex] {
				nodeDistances[sourceIndex][destIndex] = stepsAway
				nodeDistances[destIndex][sourceIndex] = stepsAway
			}
		}
	}
	nm.nodeDistances = nodeDistances
}

func (nm *nodeMap) setEndDistance() {
	nm.nodes[endIndex].visited = true
	queue := []int{endIndex}
	for len(queue) > 0 {
		n := nm.nodes[queue[0]]
		queue = queue[1:]
		for _, i := range n.adjacent {
			if !nm.nodes[i].visited {
				nm.nodes[i].visited = true
				nm.nodes[i].endDist = n.endDist + 1
				queue = append(queue, i)
			}
		}
	}
	nm.countPerEndDist = make([]int, nm.nodes[startIndex].endDist+1)
	for i, n := range nm.nodes {
		nm.countPerEndDist[n.endDist]++
		nm.nodes[i].visited = false
	}
}

func (nm *nodeMap) longestHike() int {
	return nm.recurseHike(startIndex, endIndex, nm.nodes[startIndex].endDist)
}

func (nm *nodeMap) recurseHike(fromIndex, toIndex, maxEndDist int) int {
	if fromIndex == toIndex {
		return 0
	}

	from := nm.nodes[fromIndex]
	from.visited = true
	// If every node of distance k to the end is visited, remaining path must only use nodes of
	// distance k-1 or less. Otherwise, the path will be cut off and the end cannot be reached
	nm.countPerEndDist[from.endDist]--
	if nm.countPerEndDist[from.endDist] == 0 {
		maxEndDist = from.endDist - 1
	}

	bestHike := -1
	for _, aIndex := range from.adjacent {
		adjacent := nm.nodes[aIndex]
		if adjacent.visited {
			continue
		}
		if adjacent.endDist > maxEndDist {
			continue
		}
		// If current node is an edge node, disallow moving along the edge away from the end
		if from.edgeNode && adjacent.edgeNode && from.endDist < adjacent.endDist {
			continue
		}
		if remain := nm.recurseHike(aIndex, toIndex, maxEndDist); remain >= 0 {
			if newHike := nm.nodeDistances[fromIndex][aIndex] + remain; newHike > bestHike {
				bestHike = newHike
			}
		}
	}

	nm.countPerEndDist[from.endDist]++
	from.visited = false

	return bestHike
}

func main() {
	input := util.StdinReadlines()
	a := parseArea(input)
	fmt.Println(a.longestHike())
}
