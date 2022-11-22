package main

import (
	"bufio"
	"fmt"
	"os"
)

type octopus struct {
	energy    int
	flashed   bool
	neighbors []*octopus
}

func (o octopus) String() string {
	f := 0
	if o.flashed {
		f = 1
	}
	return fmt.Sprintf("{%d %d %d}", o.energy, f, len(o.neighbors))
}

func (o *octopus) Flash() int {
	o.flashed = true

	flashes := 1
	for _, n := range o.neighbors {
		flashes += n.Increment()
	}

	return flashes
}

func (o *octopus) Increment() int {
	if o.flashed {
		return 0
	}

	o.energy++
	if o.energy > 9 {
		return o.Flash()
	}

	return 0
}

func (o *octopus) Reset() {
	if o.flashed {
		o.energy = 0
		o.flashed = false
	}
}

func (o *octopus) AddNeighbor(neighbor *octopus) {
	o.neighbors = append(o.neighbors, neighbor)
}

type swarm struct {
	octopi   []*octopus
	row, col int
}

func (s *swarm) Print() {
	fmt.Println(s.octopi[:s.row])
}

func (s *swarm) Size() int {
	return len(s.octopi)
}

func (s *swarm) Step() int {
	flashCount := 0
	for _, o := range s.octopi {
		flashCount += o.Increment()
	}

	for _, o := range s.octopi {
		o.Reset()
	}

	return flashCount
}

func NewSwarm(input []string) *swarm {
	row := len(input)
	col := len(input[0])

	octopi := make([]*octopus, 0)
	for i, line := range input {
		for j, val := range line {
			energy := int(val - '0')
			o := &octopus{energy: energy}

			neighbors := make([]*octopus, 0)
			if i != 0 {
				neighbors = append(neighbors, octopi[(i-1)*col+j])

				if j != 0 {
					neighbors = append(neighbors, octopi[(i-1)*col+j-1])
				}

				if j != col-1 {
					neighbors = append(neighbors, octopi[(i-1)*col+j+1])
				}
			}
			if j != 0 {
				neighbors = append(neighbors, octopi[i*col+j-1])
			}

			for _, n := range neighbors {
				n.AddNeighbor(o)
				o.AddNeighbor(n)
			}

			octopi = append(octopi, o)
		}
	}

	return &swarm{octopi, row, col}
}

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

func main() {
	input := readInput()
	s := NewSwarm(input)

	for i := 1; i < 10000; i++ {
		if s.Size() == s.Step() {
			fmt.Println(i)
			break
		}
	}
}
