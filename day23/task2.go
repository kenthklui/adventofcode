package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"container/list"
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

func CupString(cups *list.List, startValue int) string {
	var b strings.Builder

	for c := cups.Front(); c != nil; c = c.Next() {
		fmt.Fprintf(&b, "%d", c.Value.(int))
	}

	return b.String()
}

type Circle struct {
	min, max  int
	cups      *list.List
	cupsIndex []*list.Element
}

func (c *Circle) String() string {
	return CupString(c.cups, 0)
}

func readCircle(input []string, realMax int) *Circle {
	cups := list.New()
	cupsIndex := make([]*list.Element, realMax+1)
	min, max := 9, 1

	for i := range input[0] {
		if num, err := strconv.Atoi(input[0][i : i+1]); err == nil {
			if min > num {
				min = num
			}
			if max < num {
				max = num
			}

			e := cups.PushBack(num)
			cupsIndex[num] = e
		} else {
			panic(err)
		}
	}

	for i := max + 1; i <= realMax; i++ {
		e := cups.PushBack(i)
		cupsIndex[i] = e
	}
	max = realMax

	return &Circle{min, max, cups, cupsIndex}
}

func (c *Circle) Move() {
	current := c.cups.Front()

	pickup := make([]int, 3)

	for i := 0; i < 3; i++ {
		next := current.Next()
		pickup[i] = c.cups.Remove(next).(int)
	}

	destValue := current.Value.(int)
	for {
		if destValue == c.min {
			destValue = c.max
		} else {
			destValue--
		}

		found := false
		for i := 0; i < 3; i++ {
			found = found || (pickup[i] == destValue)
		}

		if !found {
			break
		}
	}

	dest := c.cupsIndex[destValue]

	for i := 0; i < 3; i++ {
		dest = c.cups.InsertAfter(pickup[i], dest)
		c.cupsIndex[pickup[i]] = dest
	}

	front := c.cups.Remove(c.cups.Front())
	e := c.cups.PushBack(front)
	c.cupsIndex[front.(int)] = e
}

func (c *Circle) Solution() string {
	l := list.New()
	l.PushBackList(c.cups)
	l.PushBackList(c.cups)

	start := l.Front()
	for start.Value.(int) != 1 {
		start = start.Next()
	}

	start = start.Next()
	star1 := start.Value.(int)
	start = start.Next()
	star2 := start.Value.(int)

	return fmt.Sprintf("%d", star1*star2)
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
