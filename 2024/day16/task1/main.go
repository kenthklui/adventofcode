package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type point struct {
	x, y      int
	neighbors [4]*point
}

type move struct {
	p      *point
	facing int
	score  int

	// Index is needed for the priority queue
	index int
}

type moveQueue []*move

func (mq moveQueue) Len() int           { return len(mq) }
func (mq moveQueue) Less(i, j int) bool { return mq[i].score < mq[j].score }
func (mq moveQueue) Swap(i, j int) {
	mq[i], mq[j] = mq[j], mq[i]
	mq[i].index = i
	mq[j].index = j
}

func (mq *moveQueue) Push(x any) {
	n := len(*mq)
	p := x.(*move)
	p.index = n
	*mq = append(*mq, p)
}

func (mq *moveQueue) Pop() any {
	old := *mq
	n := len(old)
	p := old[n-1]
	old[n-1] = nil
	p.index = -1
	*mq = old[:n-1]

	return p
}

type maze struct {
	height, width              int
	startX, startY, endX, endY int
	grid                       [][]*point
}

func newMaze(input []string) *maze {
	m := &maze{height: len(input), width: len(input[0])}
	grid := make([][]*point, len(input))
	for y, line := range input {
		grid[y] = make([]*point, len(line))
		for x, char := range line {
			switch char {
			case '#':
				continue
			case 'S':
				m.startX, m.startY = x, y
			case 'E':
				m.endX, m.endY = x, y
			}

			p := point{x, y, [4]*point{}}
			if x > 0 && grid[y][x-1] != nil {
				p.neighbors[3] = grid[y][x-1]
				grid[y][x-1].neighbors[1] = &p
			}
			if y > 0 && grid[y-1][x] != nil {
				p.neighbors[2] = grid[y-1][x]
				grid[y-1][x].neighbors[0] = &p
			}
			grid[y][x] = &p
		}
	}
	m.grid = grid

	return m
}

func (m maze) isEnd(p *point) bool { return p.x == m.endX && p.y == m.endY }

var COSTS = []int{1, 1000}

type moveTracker [][][4]*move

func newMoveTracker(height, width int) moveTracker {
	mt := make([][][4]*move, height)
	for y := range mt {
		mt[y] = make([][4]*move, width)
	}
	return mt
}

func (mt moveTracker) createOrUpdate(p *point, facing int, newScore int) (m *move, isNew bool) {
	m = mt[p.y][p.x][facing]
	if m == nil {
		m = &move{
			p:      p,
			facing: facing,
			score:  newScore,
		}
		mt[p.y][p.x][facing] = m
		isNew = true
	} else if m.score > newScore {
		m.score = newScore
	}
	return
}

func (m *maze) solve() int {
	mt := newMoveTracker(m.height, m.width)

	initial, _ := mt.createOrUpdate(m.grid[m.startY][m.startX], 1, 0)
	mq := moveQueue{initial}
	for len(mq) > 0 {
		head := heap.Pop(&mq).(*move)
		if m.isEnd(head.p) {
			return head.score
		}

		if forward := head.p.neighbors[head.facing]; forward != nil {
			newScore := head.score + COSTS[0]
			if nextMove, isNew := mt.createOrUpdate(forward, head.facing, newScore); isNew {
				heap.Push(&mq, nextMove)
			}
		}

		for i := 1; i <= 3; i++ {
			if newFacing := (head.facing + i) % 4; head.p.neighbors[newFacing] != nil {
				turnScore := head.score + (2-(i%2))*COSTS[1]
				if nextMove, isNew := mt.createOrUpdate(head.p, newFacing, turnScore); isNew {
					heap.Push(&mq, nextMove)
				}
			}
		}
	}

	return -1
}

func solve(input []string) (output string) {
	m := newMaze(input)
	return strconv.Itoa(m.solve())
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
