package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type vec2 struct {
	x, y int
}

type point struct {
	x, y int

	neighbors [4]*point
}

type move struct {
	p      *point
	facing int

	score   int
	visited bool

	// Index is needed for the priority queue
	index int
}

func (m *move) inQueue() bool {
	return m.index > -1
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

func (m *maze) solve() int {
	tracker := make([][][4]*move, m.height)
	for y := range tracker {
		tracker[y] = make([][4]*move, m.width)
	}

	initial := &move{
		p:      m.grid[m.startY][m.startX],
		facing: 1,
	}
	tracker[initial.p.y][initial.p.x][initial.facing] = initial

	mq := moveQueue{initial}
	for len(mq) > 0 {
		head := heap.Pop(&mq).(*move)
		if m.isEnd(head.p) {
			return head.score
		}
		head.visited = true

		if forward := head.p.neighbors[head.facing]; forward != nil {
			newScore := head.score + COSTS[0]
			if tracker[forward.y][forward.x][head.facing] == nil {
				nextMove := &move{
					p:      forward,
					facing: head.facing,
					score:  newScore,
				}
				tracker[forward.y][forward.x][head.facing] = nextMove
				heap.Push(&mq, nextMove)
			} else if tracker[forward.y][forward.x][head.facing].score > newScore {
				tracker[forward.y][forward.x][head.facing].score = newScore
			}
		}

		for i := 1; i <= 3; i++ {
			if newFacing := (head.facing + i) % 4; head.p.neighbors[newFacing] != nil {
				turnScore := head.score + (2-(i%2))*COSTS[1]
				if tracker[head.p.y][head.p.x][newFacing] == nil {
					nextMove := &move{
						p:      head.p,
						facing: newFacing,
						score:  turnScore,
					}
					tracker[head.p.y][head.p.x][newFacing] = nextMove
					heap.Push(&mq, nextMove)
				} else if tracker[head.p.y][head.p.x][head.facing].score > turnScore {
					tracker[head.p.y][head.p.x][head.facing].score = turnScore
				}

				if forward := head.p.neighbors[head.facing]; forward != nil {
					newScore := turnScore + COSTS[0]
					if tracker[forward.y][forward.x][head.facing] == nil {
						nextMove := &move{
							p:      forward,
							facing: head.facing,
							score:  newScore,
						}
						tracker[forward.y][forward.x][head.facing] = nextMove
						heap.Push(&mq, nextMove)
					} else if tracker[forward.y][forward.x][head.facing].score > newScore {
						tracker[forward.y][forward.x][head.facing].score = newScore
					}
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
