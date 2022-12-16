package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type valve struct {
	id, flow   int
	name       string
	travelCost map[*valve]int
}

func (v valve) String() string { return v.name }

type state struct {
	flow, score, time int
	key               uint64
	curr              *valve
	unopened          []*valve
}

func NewState(time int, valves map[string]*valve) *state {
	start := valves["AA"]
	s := state{
		flow:     0,
		score:    0,
		time:     time,
		key:      uint64(0),
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
		key:      s.key,
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

	ns.key += uint64(1) << dest.id

	end := len(s.unopened) - 1
	ns.unopened[unopenIndex], ns.unopened[end] = ns.unopened[end], ns.unopened[unopenIndex]
	ns.unopened = ns.unopened[:end]

	return ns
}

type memo struct {
	states map[uint64]*state
}

func (s *state) recursiveFill(m *memo, maxFlow int) {
	prevState, prevStateOk := m.states[s.key]
	if prevStateOk {
		if s.score+maxFlow*s.time <= prevState.score {
			return
		}
	}
	for unopenIndex := range s.unopened {
		if ns := s.openValve(unopenIndex); ns != nil {
			ns.recursiveFill(m, maxFlow)
		}
	}

	score := s.score + s.time*s.flow
	if !prevStateOk || score > prevState.score {
		s.score = score
		m.states[s.key] = s
	}
}

func (s *state) recursiveOpen(maxScore, maxFlow int) int {
	if s.score+maxFlow*s.time <= maxScore {
		return maxScore
	}

	for unopenIndex := range s.unopened {
		if ns := s.openValve(unopenIndex); ns != nil {
			if score := ns.recursiveOpen(maxScore, maxFlow); score > maxScore {
				maxScore = score
			}
		}
	}
	if score := s.score + s.time*s.flow; score > maxScore {
		maxScore = score
	}

	return maxScore
}

func openValves(valves map[string]*valve) int {
	maxTime := 26

	maxFlow := 0
	for _, v := range valves {
		maxFlow += v.flow
	}

	s := NewState(maxTime, valves)
	m := memo{make(map[uint64]*state)}
	s.recursiveFill(&m, maxFlow)

	var maxScore int
	for _, s := range m.states {
		s.time = maxTime
		s.curr = valves["AA"]
		s.flow = 0

		if score := s.recursiveOpen(maxScore, maxFlow); score > maxScore {
			maxScore = score
		}
	}

	return maxScore
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
	for i, line := range input {
		v := valve{id: i}
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
