package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

type Cups []int

type Circle struct {
	min, max, front int
	cups            Cups
}

func readCircle(input []string, realMax int) *Circle {
	cups := make([]int, realMax+1)

	min, max := 9, 1
	var front, prev int
	for i := range input[0] {
		num, err := strconv.Atoi(input[0][i : i+1])
		if err != nil {
			panic(err)
		}

		if min > num {
			min = num
		}
		if max < num {
			max = num
		}

		if prev == 0 {
			front = num
		} else {
			cups[prev] = num
		}
		prev = num
	}

	if realMax > max {
		cups[prev] = max + 1
		for i := max + 1; i < realMax; i++ {
			cups[i] = i + 1
		}
		cups[realMax] = front
	} else {
		cups[prev] = front
	}

	return &Circle{min, realMax, front, cups}
}

func (c *Circle) NextDest(dest int) int {
	dest--
	if dest < c.min {
		dest = c.max
	}
	return dest
}

func (c *Circle) Move() {
	pickup := make([]int, 3)
	pickup[0] = c.cups[c.front]
	pickup[1] = c.cups[pickup[0]]
	pickup[2] = c.cups[pickup[1]]

	dest := c.NextDest(c.front)
	for dest == pickup[0] || dest == pickup[1] || dest == pickup[2] {
		dest = c.NextDest(dest)
	}

	c.cups[c.front] = c.cups[pickup[2]]
	c.cups[pickup[2]] = c.cups[dest]
	c.cups[dest] = pickup[0]
	c.front = c.cups[c.front]
}

func (c *Circle) Solution() string {
	value1 := c.cups[1]
	value2 := c.cups[value1]

	return fmt.Sprintf("%d", value1*value2)
}

func main() {
	realMax := 1000000
	moves := 10000000

	input := readInput()
	c := readCircle(input, realMax)

	for i := 0; i < moves; i++ {
		c.Move()
	}

	fmt.Println(c.Solution())
}
