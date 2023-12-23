package main

import (
	"fmt"
	"slices"
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
func (v vec) inBounds(width, height int) bool {
	return v.x >= 0 && v.y >= 0 && v.x < width && v.y < height
}

type area struct {
	start, end    vec
	width, height int
	aMap          []string
}

func parseArea(input []string) area {
	slopes := []string{"^", ">", "v", "<"}
	for i := range input {
		for _, s := range slopes {
			input[i] = strings.ReplaceAll(input[i], s, ".")
		}
	}

	return area{
		start:  vec{1, 0},
		end:    vec{len(input[0]) - 2, len(input) - 1},
		width:  len(input[0]),
		height: len(input),
		aMap:   input,
	}
}

func (a area) neighbors(v vec) []vec {
	n := make([]vec, 0, len(dirs))
	for _, dir := range dirs {
		next := v.add(dir)
		if !next.inBounds(a.width, a.height) {
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

func (a area) longestHike() int {
	nm := buildNodeMap(a)
	return nm.longestHike(a.start, a.end)
}

type node struct {
	index    int
	v        vec
	adjacent []int
}

type nodeMap struct {
	nodes         []*node
	nodeDistances [][]int
}

func buildNodeMap(a area) *nodeMap {
	nodes := make([]*node, 0)
	nodes = append(nodes, &node{0, a.start, make([]int, 0)})
	nodes = append(nodes, &node{1, a.end, make([]int, 0)})

	for y, row := range a.aMap {
		for x, c := range row {
			if c != '.' {
				continue
			}

			curr := vec{x, y}
			if len(a.neighbors(curr)) > 2 {
				index := len(nodes)
				nodes = append(nodes, &node{
					index:    index,
					v:        curr,
					adjacent: make([]int, 0),
				})
			}
		}
	}

	nm := &nodeMap{nodes: nodes}
	nm.computeNodeDistances(a)

	return nm
}

func (nm *nodeMap) computeNodeDistances(a area) {
	nodeIndices := make(map[vec]int)
	for _, n := range nm.nodes {
		nodeIndices[n.v] = n.index
	}

	nodeDistances := make([][]int, len(nm.nodes))
	for i := range nodeDistances {
		nodeDistances[i] = make([]int, len(nm.nodes))
	}
	for _, sourceNode := range nm.nodes {
		var destIndex int
		sourceIndex := sourceNode.index
		for _, neighbor := range a.neighbors(sourceNode.v) {
			isFork, dest, stepsAway := a.nextFork(sourceNode.v, neighbor, 1)
			if isFork {
				destIndex = nodeIndices[dest]
			} else if dest == a.start {
				destIndex = 0
			} else if dest == a.end {
				destIndex = 1
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

func (nm nodeMap) longestHike(from, to vec) int {
	fromIndex := slices.IndexFunc(nm.nodes, func(n *node) bool { return n.v == from })
	toIndex := slices.IndexFunc(nm.nodes, func(n *node) bool { return n.v == to })
	if fromIndex == -1 || toIndex == -1 {
		panic("Invalid start/stop locations for computing hike distance")
	}
	travelled := make([]bool, len(nm.nodes))
	return nm.recurseHike(fromIndex, toIndex, travelled)
}

func (nm nodeMap) recurseHike(fromIndex, toIndex int, travelled []bool) int {
	if fromIndex == toIndex {
		return 0
	}

	travelled[fromIndex] = true
	hikes := make([]int, 0, len(nm.nodes[fromIndex].adjacent))
	for _, aIndex := range nm.nodes[fromIndex].adjacent {
		if travelled[aIndex] {
			continue
		}
		if remain := nm.recurseHike(aIndex, toIndex, travelled); remain >= 0 {
			hikes = append(hikes, nm.nodeDistances[fromIndex][aIndex]+remain)
		}
	}
	travelled[fromIndex] = false

	if len(hikes) > 0 {
		return slices.Max(hikes)
	} else {
		return -1
	}
}

func main() {
	input := util.StdinReadlines()
	a := parseArea(input)
	fmt.Println(a.longestHike())
}
