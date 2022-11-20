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
	min, max int
	cups     *list.List
}

func (c *Circle) String() string {
	return CupString(c.cups, 0)
}

func readCircle(input []string) *Circle {
	cups := list.New()
	min, max := 9, 1

	for i := range input[0] {
		if num, err := strconv.Atoi(input[0][i : i+1]); err == nil {
			if min > num {
				min = num
			}
			if max < num {
				max = num
			}

			cups.PushBack(num)
		} else {
			panic(err)
		}
	}

	return &Circle{min, max, cups}
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

	dest := current
	for dest.Value.(int) != destValue {
		dest = dest.Next()
	}

	for i := 0; i < 3; i++ {
		dest = c.cups.InsertAfter(pickup[i], dest)
	}

	c.cups.PushBack(c.cups.Remove(c.cups.Front()))
}

func (c *Circle) Solution() string {
	var b strings.Builder

	l := list.New()
	l.PushBackList(c.cups)
	l.PushBackList(c.cups)

	start := l.Front()
	for start.Value.(int) != 1 {
		start = start.Next()
	}

	for e := start.Next(); e.Value.(int) != 1; e = e.Next() {
		fmt.Fprintf(&b, "%d", e.Value.(int))
	}

	return b.String()
}

func main() {
	moves := 100

	input := readInput()
	c := readCircle(input)

	for i := 0; i < moves; i++ {
		c.Move()
	}

	fmt.Println(c.Solution())
}
