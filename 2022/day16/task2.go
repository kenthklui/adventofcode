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

type state struct {
	flow, score, time int
	key               uint64
	nextAvail         []int
	curr, unopened    []*valve
}

func NewState(actors, time int, valves map[string]*valve) *state {
	start := valves["AA"]
	s := state{
		flow:      0,
		score:     0,
		time:      time,
		key:       0,
		nextAvail: make([]int, actors),
		curr:      make([]*valve, actors),
		unopened:  make([]*valve, 0, len(valves)),
	}

	for i := 0; i < actors; i++ {
		s.curr[i] = start
		s.nextAvail[i] = time
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
		flow:      s.flow,
		score:     s.score,
		time:      s.time,
		key:       s.key,
		nextAvail: make([]int, len(s.nextAvail)),
		curr:      make([]*valve, len(s.curr)),
		unopened:  make([]*valve, len(s.unopened)),
	}

	copy(ns.nextAvail, s.nextAvail)
	copy(ns.curr, s.curr)
	copy(ns.unopened, s.unopened)

	return ns
}

func (s *state) moveActor(actor, unopenIndex int) *state {
	dest := s.unopened[unopenIndex]

	timeCost := s.curr[actor].travelCost[dest] + 1
	eta := s.time - timeCost
	if eta <= 0 { // out of time
		return nil
	}

	ns := s.Copy()
	ns.curr[actor] = dest
	ns.nextAvail[actor] = eta

	end := len(s.unopened) - 1
	ns.unopened[unopenIndex], ns.unopened[end] = ns.unopened[end], ns.unopened[unopenIndex]
	ns.unopened = ns.unopened[:end]

	return ns
}

func (s *state) nextAction() ([]int, int) {
	var nextTime int
	actors := make([]int, 0)
	for i, t := range s.nextAvail {
		if t > nextTime {
			actors = []int{i}
			nextTime = t
		} else if t == nextTime {
			actors = append(actors, i)
		}
	}
	return actors, nextTime
}

func (s *state) openValve(v *valve) {
	// Use a bitmask to store valve open status
	if (s.key>>v.id)&1 == 0 {
		s.flow += v.flow
		s.key |= uint64(1) << v.id
	}
}

type trace struct {
	time, actor, nextAvail int
	valveName              string
}

func (s *state) recursiveOpen(path []trace, maxScore, maxFlow int) int {
	// Which actors are active (aka. not enroute) at this time stamp?
	// And what is the time?
	actors, nextTime := s.nextAction()

	// Update score pressure
	duration := s.time - nextTime
	s.score += s.flow * duration
	s.time = nextTime

	// Shortcut: if it's impossible to beat the already known max score, exit
	if s.score+s.time*maxFlow <= maxScore {
		return maxScore
	}

	// Open any pending valves
	for _, actor := range actors {
		s.openValve(s.curr[actor])
	}

	for _, actor := range actors {
		// Check every possible move by the actor
		for unopenIndex, dest := range s.unopened {
			if ns := s.moveActor(actor, unopenIndex); ns != nil {
				newPath := append(path, trace{ns.time, actor, ns.nextAvail[actor], dest.name})
				if score := ns.recursiveOpen(newPath, maxScore, maxFlow); score > maxScore {
					maxScore = score
				}
			}
		}

		// Check for actor waiting 1 minute doing nothing
		// Sometimes, the best strategy is waiting for the other guy to do something better
		// This is likely necessary until we are running out of valves to open
		if len(s.unopened) < len(s.curr)*2 {
			ns := s.Copy()
			ns.nextAvail[actor]--
			if score := ns.recursiveOpen(path, maxScore, maxFlow); score > maxScore {
				maxScore = score
			}
		}
	}

	// Check the "wait until end of time" case
	if waitScore := s.score + s.time*s.flow; waitScore > maxScore {
		maxScore = waitScore
	}

	// Debug printout
	// if len(path) <= 2 {
	// 	fmt.Println(path)
	// }
	return maxScore
}

func openValves(valves map[string]*valve) int {
	s := NewState(2, 26, valves)

	maxFlow := 0
	for _, v := range valves {
		maxFlow += v.flow
	}

	return s.recursiveOpen([]trace{}, 0, maxFlow)
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
	score := openValves(valves)

	fmt.Println(score)
}
