package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type valve struct {
	flow       int
	name       string
	travelCost map[*valve]int
}

func (v valve) String() string { return v.name }

type state struct {
	flow, score, time int
	curr              *valve
	unopened          []*valve
}

func NewState(time int, valves map[string]*valve) *state {
	start := valves["AA"]
	s := state{
		flow:     0,
		score:    0,
		time:     time,
		curr:     valves["AA"],
		unopened: make([]*valve, 0, len(valves)),
	}

	for _, v := range valves {
		if v == start {
			continue
		}
		s.unopened = append(s.unopened, v)
	}

	return &s
}

func (s *state) Copy() *state {
	ns := &state{
		flow:     s.flow,
		score:    s.score,
		time:     s.time,
		curr:     s.curr,
		unopened: make([]*valve, len(s.unopened)),
	}

	copy(ns.unopened, s.unopened)

	return ns
}

func (s *state) openValve(unopenIndex int) *state {
	dest := s.unopened[unopenIndex]

	timeCost := s.curr.travelCost[dest] + 1
	eta := s.time - timeCost
	if eta <= 0 {
		return nil
	}

	ns := s.Copy()
	ns.score += ns.flow * timeCost
	ns.flow += dest.flow
	ns.time = eta
	ns.curr = dest

	end := len(s.unopened) - 1
	ns.unopened[unopenIndex], ns.unopened[end] = ns.unopened[end], ns.unopened[unopenIndex]
	ns.unopened = ns.unopened[:end]

	return ns
}

type trace struct {
	time      int
	valveName string
}

func (s *state) recursiveOpen(path []trace) int {
	var maxScore int
	for unopenIndex, dest := range s.unopened {
		if ns := s.openValve(unopenIndex); ns != nil {
			newPath := append(path, trace{ns.time, dest.name})
			if score := ns.recursiveOpen(newPath); score > maxScore {
				maxScore = score
			}
		}
	}
	if maxScore == 0 {
		maxScore = s.score + s.time*s.flow
	}

	return maxScore
}

func openValves(valves map[string]*valve) int {
	s := NewState(30, valves)
	return s.recursiveOpen([]trace{})
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

func parseInput(input []string) map[string]*valve {
	destinations := make(map[string][]string)
	valves := make(map[string]*valve)
	for _, line := range input {
		var v valve
		if n, err := fmt.Sscanf(line, "Valve %s has flow rate=%d", &v.name, &v.flow); err != nil {
			panic(err)
		} else if n != 2 {
			panic("Failed")
		}

		split := strings.Split(line, " to valve")
		destinations[v.name] = strings.Split(strings.Trim(split[1], "s "), ", ")

		valves[v.name] = &v
	}

	for _, v := range valves {
		v.travelCost = make(map[*valve]int)
		for _, destName := range destinations[v.name] {
			dest := valves[destName]
			v.travelCost[dest] = 1
		}
	}

	for _, source := range valves {
		// Use BFS to compute distance between every pair of valves
		queue := destinations[source.name]
		for index := 0; index < len(queue); index++ {
			edgeName := queue[index]
			edge := valves[edgeName]
			edgeTravelTime := source.travelCost[edge]

			for _, neighborName := range destinations[edgeName] {
				neighbor := valves[neighborName]
				if _, ok := source.travelCost[neighbor]; !ok {
					source.travelCost[neighbor] = edgeTravelTime + 1
					queue = append(queue, neighborName)
				}
			}
		}
	}

	// Purge useless destinations with no flow
	for k, v1 := range valves {
		if v1.flow == 0 && v1.name != "AA" {
			delete(valves, k)
			for _, v2 := range valves {
				if _, ok := v2.travelCost[v1]; ok {
					delete(v2.travelCost, v1)
				}
			}
		}
	}

	return valves
}

func main() {
	input := readInput()
	valves := parseInput(input)

	fmt.Println(openValves(valves))
}
