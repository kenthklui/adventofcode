package main

import (
	"container/heap"
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

const maxStraight = 3

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

type node struct {
	loc                  vec
	dir, steps, heatLoss int
	index                int
}

func (n *node) inQueue() bool { return n.index > -1 }

type nodeQueue []*node

func (nq nodeQueue) Len() int           { return len(nq) }
func (nq nodeQueue) Less(i, j int) bool { return nq[i].heatLoss < nq[j].heatLoss }
func (nq nodeQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
	nq[i].index, nq[j].index = i, j
}
func (nq *nodeQueue) Push(x any) {
	n := len(*nq)
	p := x.(*node)
	p.index = n
	*nq = append(*nq, p)
}
func (nq *nodeQueue) Pop() any {
	old := *nq
	n := len(old)
	p := old[n-1]
	old[n-1] = nil
	p.index = -1
	*nq = old[:n-1]
	return p
}

type agent struct {
	width, height int
	heatLoss      [][]int
	nodes         [][][4][maxStraight]*node
	searchQueue   nodeQueue
}

func makeAgent(input []string) *agent {
	width := len(input[0])
	height := len(input)

	heatLoss := make([][]int, height)
	nodes := make([][][4][maxStraight]*node, height)
	for y, line := range input {
		heatLoss[y] = make([]int, width)
		nodes[y] = make([][4][maxStraight]*node, width)
		for x, c := range line {
			heatLoss[y][x] = int(c - '0')
		}
	}
	searchQueue := make(nodeQueue, 0, width*height*maxStraight)

	return &agent{width, height, heatLoss, nodes, searchQueue}
}

func (a *agent) exit(loc vec) bool { return loc.x == a.width-1 && loc.y == a.height-1 }

func (a *agent) search() int {
	for len(a.searchQueue) > 0 {
		n := heap.Pop(&a.searchQueue).(*node)
		if a.exit(n.loc) {
			return n.heatLoss
		}
		a.iterate(n)
	}
	return -1
}

func (a *agent) addToSearch(loc vec, dir, steps, heatLoss int) {
	n := a.nodes[loc.y][loc.x][dir][steps-1]
	if n == nil {
		n = &node{loc: loc, dir: dir, steps: steps, heatLoss: heatLoss}
		a.nodes[loc.y][loc.x][dir][steps-1] = n
		heap.Push(&a.searchQueue, n)
	} else if heatLoss < n.heatLoss {
		n.heatLoss = heatLoss
		if n.inQueue() {
			heap.Fix(&a.searchQueue, n.index)
		} else {
			heap.Push(&a.searchQueue, n)
		}
	}
}

func (a *agent) iterate(n *node) {
	leftDir, rightDir := (n.dir+3)%4, (n.dir+1)%4
	left, right := n.loc.add(dirs[leftDir]), n.loc.add(dirs[rightDir])
	if left.inBounds(a.width, a.height) {
		newHeatLoss := n.heatLoss + a.heatLoss[left.y][left.x]
		a.addToSearch(left, leftDir, 1, newHeatLoss)
	}
	if right.inBounds(a.width, a.height) {
		newHeatLoss := n.heatLoss + a.heatLoss[right.y][right.x]
		a.addToSearch(right, rightDir, 1, newHeatLoss)
	}

	if n.steps < maxStraight {
		forward := n.loc.add(dirs[n.dir])
		if forward.inBounds(a.width, a.height) {
			newHeatLoss := n.heatLoss + a.heatLoss[forward.y][forward.x]
			newSteps := n.steps + 1
			a.addToSearch(forward, n.dir, newSteps, newHeatLoss)
		}
	}
}

func traverse(input []string) int {
	a := makeAgent(input)
	origin := vec{0, 0}
	heap.Push(&a.searchQueue, &node{loc: origin, dir: 1, steps: 0, heatLoss: 0})
	heap.Push(&a.searchQueue, &node{loc: origin, dir: 2, steps: 0, heatLoss: 0})
	return a.search()
}

func main() {
	input := util.StdinReadlines()
	fmt.Println(traverse(input))
}
