package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type location interface {
	String() string
	name() string
	pathCount(extraVisit int) int
	visitable(extraVisit int) bool
	addConnection(l location)
}

//

type cave struct {
	nameStr        string
	large, visited bool
	destinations   []location
}

func (c *cave) String() string {
	var size string
	if c.large {
		size = "large"
	} else {
		size = "small"
	}

	dests := make([]string, len(c.destinations))
	for i, d := range c.destinations {
		dests[i] = d.name()
	}

	return fmt.Sprintf("{%s (%s), goes to: [%s]}", c.nameStr, size, strings.Join(dests, ", "))
}

func (c *cave) name() string { return c.nameStr }

func (c *cave) pathCount(extraVisit int) int {
	count := 0
	if !c.large && c.visited {
		extraVisit--
		for _, d := range c.destinations {
			if d.visitable(extraVisit) {
				count += d.pathCount(extraVisit)
			}
		}
	} else {
		c.visited = true
		for _, d := range c.destinations {
			if d.visitable(extraVisit) {
				count += d.pathCount(extraVisit)
			}
		}
		c.visited = false
	}
	return count
}

func (c *cave) visitable(extraVisit int) bool {
	return c.large || !c.visited || (extraVisit > 0)
}

func (c *cave) addConnection(l location) {
	c.destinations = append(c.destinations, l)
}

//

type start struct {
	destinations []location
}

func (s *start) String() string {
	dests := make([]string, len(s.destinations))
	for i, d := range s.destinations {
		dests[i] = d.name()
	}

	return fmt.Sprintf("{start, goes to: [%s]}", strings.Join(dests, ", "))
}
func (s *start) name() string { return "start" }
func (s *start) pathCount(extraVisit int) int {
	count := 0
	for _, d := range s.destinations {
		count += d.pathCount(extraVisit)
	}

	return count
}
func (s *start) visitable(extraVisit int) bool { return true } // just for fulfilling the interface
func (s *start) addConnection(l location) {
	s.destinations = append(s.destinations, l)
}

//

type end struct{}

func (e *end) String() string                { return "{end}" }
func (e *end) name() string                  { return "end" }
func (e *end) pathCount(extraVisit int) int  { return 1 }
func (e *end) visitable(extraVisit int) bool { return true }
func (e *end) addConnection(l location)      {}

//

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

func parseInput(input []string) map[string]location {
	caverns := make(map[string]location)

	s := &start{destinations: make([]location, 0)}
	e := &end{}

	caverns["start"] = s
	caverns["end"] = e

	for _, line := range input {
		edge := strings.Split(line, "-")
		from, to := edge[0], edge[1]

		fromCave, ok := caverns[from]
		if !ok {
			large := (from == strings.ToUpper(from))
			fromCave = &cave{
				nameStr:      from,
				large:        large,
				destinations: make([]location, 0),
			}
			caverns[from] = fromCave
		}

		toCave, ok := caverns[to]
		if !ok {
			large := (to == strings.ToUpper(to))
			toCave = &cave{
				nameStr:      to,
				large:        large,
				destinations: make([]location, 0),
			}
			caverns[to] = toCave
		}

		if toCave.name() != "start" {
			fromCave.addConnection(toCave)
		}
		if fromCave.name() != "start" {
			toCave.addConnection(fromCave)
		}
	}

	return caverns
}

func main() {
	input := readInput()
	caverns := parseInput(input)

	// fmt.Println(caverns)
	fmt.Println(caverns["start"].pathCount(1))
}
