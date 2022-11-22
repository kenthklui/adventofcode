package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Basically djikstra's on a integer grid

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

type point struct {
	x, y    int
	risk    int
	visited bool

	distance int
	// Index is needed for the priority queue
	index int
}

func (p *point) inQueue() bool {
	return p.index > -1
}

func (p point) String() string {
	return fmt.Sprintf("[%d]", p.risk)
}

type pointQueue []*point

func (pq pointQueue) Len() int           { return len(pq) }
func (pq pointQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }
func (pq pointQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *pointQueue) Push(x any) {
	n := len(*pq)
	p := x.(*point)
	p.index = n
	*pq = append(*pq, p)
}

func (pq *pointQueue) Pop() any {
	old := *pq
	n := len(old)
	p := old[n-1]
	old[n-1] = nil
	p.index = -1
	*pq = old[:n-1]

	return p
}

type cave [][]*point

func (c cave) traverse() int {
	xEnd, yEnd := len(c)-1, len(c[0])-1

	c[0][0].distance = 0
	pq := make(pointQueue, 0)
	heap.Push(&pq, c[0][0])

	for {
		current := heap.Pop(&pq).(*point)
		if current.x == xEnd && current.y == yEnd {
			return current.distance
		}
		current.visited = true

		// fmt.Printf("Visiting (%d, %d)\n", current.x, current.y)

		for _, pair := range neighbors(current.x, current.y, xEnd, yEnd) {
			neighbor := c[pair[0]][pair[1]]
			if neighbor.visited {
				continue
			}

			neighborModified := false
			distance := current.distance + neighbor.risk
			if distance < neighbor.distance {
				neighbor.distance = distance
				neighborModified = true
			}

			if neighbor.inQueue() {
				if neighborModified {
					heap.Fix(&pq, neighbor.index)
				}
			} else {
				// fmt.Printf("Adding (%d, %d) to queue\n", neighbor.x, neighbor.y)
				heap.Push(&pq, neighbor)
			}
		}

	}
}

func neighbors(x, y, xEnd, yEnd int) [][]int {
	up := []int{x, y - 1}
	down := []int{x, y + 1}
	left := []int{x - 1, y}
	right := []int{x + 1, y}

	if x == 0 {
		if y == 0 {
			return [][]int{right, down}
		} else if y == yEnd {
			return [][]int{right, up}
		} else {
			return [][]int{right, down, up}
		}
	} else if x == xEnd {
		if y == 0 {
			return [][]int{down, left}
		} else if y == yEnd {
			// Should never happen...
			// return [][]int{up, left}
			panic("Calling neighbors at the end")
		} else {
			return [][]int{down, left, up}
		}
	} else {
		if y == 0 {
			return [][]int{right, down, left}
		} else if y == yEnd {
			return [][]int{right, left, up}
		} else {
			return [][]int{right, down, left, up}
		}
	}
}

func parseInput(input []string) cave {
	scaleFactor := 5
	c := make(cave, len(input) * scaleFactor)

	for i, line := range input {
		for m := 0; m < scaleFactor; m++ {
			x := m * len(input) + i
			c[x] = make([]*point, len(line) * scaleFactor)

			for j, r := range line {

				for n := 0; n < scaleFactor; n++ {
					y := n * len(line) + j

					risk := int(r - '0') + m + n
					if risk > 9 {
						risk -= 9
					}

					c[x][y] = &point{
						x:        x,
						y:        y,
						risk:     risk,
						distance: 9 * scaleFactor * scaleFactor * len(input) * len(line),
						visited:  false,
						index:    -1,
					}
				}
			}
		}
	}

	return c
}

func main() {
	input := readInput()
	c := parseInput(input)

	fmt.Println(c.traverse())
}
