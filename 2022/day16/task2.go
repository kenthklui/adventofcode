package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const workerCount = 8

type valve struct {
	id, flow int
	name     string
}

func (v valve) String() string { return v.name }

type state struct {
	flow, score, time, currId int
	key                       uint64
	unopened                  []int
}

func NewState(time, startId int, goodValves valveList) *state {
	s := state{
		currId:   startId,
		flow:     0,
		score:    0,
		time:     time,
		key:      uint64(0),
		unopened: make([]int, 0, len(goodValves)),
	}

	for _, v := range goodValves {
		s.unopened = append(s.unopened, v.id)
	}

	return &s
}

func (s *state) Copy() *state {
	ns := &state{
		currId:   s.currId,
		flow:     s.flow,
		score:    s.score,
		time:     s.time,
		key:      s.key,
		unopened: make([]int, len(s.unopened)),
	}
	copy(ns.unopened, s.unopened)

	return ns
}

type valveList []*valve

func (vl valveList) Len() int           { return len(vl) }
func (vl valveList) Less(i, j int) bool { return vl[i].flow > vl[j].flow }
func (vl valveList) Swap(i, j int)      { vl[i], vl[j] = vl[j], vl[i] }

type cave struct {
	maxFlow            int
	start              *valve
	valves, goodValves valveList
	travelCost         [][]int
}

func NewCave(valves map[string]*valve, tunnels map[string][]string) cave {
	var c cave

	c.start = valves["AA"]

	// Store valves in a slice for faster access
	// Also create list of useful tunnels with non-zero flow
	c.valves = make(valveList, len(valves))
	c.goodValves = make(valveList, 0, len(valves))
	for _, v := range valves {
		c.valves[v.id] = v
		if v.flow > 0 {
			c.goodValves = append(c.goodValves, v)
			c.maxFlow += v.flow
		}
	}
	sort.Sort(c.goodValves)

	// Floyd-Warshall for setting travel costs
	cost := make([][]int, len(valves))
	for i := range cost {
		cost[i] = make([]int, len(valves))
		for j := range cost[i] {
			cost[i][j] = len(valves)
		}
		cost[i][i] = 0
	}
	for source, dests := range tunnels {
		for _, dest := range dests {
			sId, dId := valves[source].id, valves[dest].id
			cost[sId][dId] = 1
		}
	}
	for mid := range cost {
		for from := range cost {
			for to := range cost {
				if cost[from][to] > cost[from][mid]+cost[mid][to] {
					cost[from][to] = cost[from][mid] + cost[mid][to]
				}
			}
		}
	}
	c.travelCost = cost

	return c
}

func (c cave) openValve(s *state, unopenedId int) *state {
	destId := s.unopened[unopenedId]

	timeCost := c.travelCost[s.currId][destId] + 1
	eta := s.time - timeCost
	if eta <= 0 {
		return nil
	}

	ns := s.Copy()
	ns.score += ns.flow * timeCost
	ns.flow += c.valves[destId].flow
	ns.time = eta
	ns.currId = destId

	ns.key += uint64(1) << destId

	end := len(s.unopened) - 1
	ns.unopened[unopenedId], ns.unopened[end] = ns.unopened[end], ns.unopened[unopenedId]
	ns.unopened = ns.unopened[:end]

	return ns
}

func (c *cave) recursiveFill(s *state, memo map[uint64]*state) {
	prevState, ok := memo[s.key]
	if ok {
		maxPossibleScore := s.score + c.maxFlow*(s.time-1)
		if maxPossibleScore <= prevState.score {
			return
		}
	}

	for unopenedId := range s.unopened {
		if ns := c.openValve(s, unopenedId); ns != nil {
			c.recursiveFill(ns, memo)
		}
	}

	waitScore := s.score + s.time*s.flow
	if !ok || waitScore > prevState.score {
		s.score = waitScore
		memo[s.key] = s
	}
}

func (c *cave) recursiveOpen(s *state, maxScore *int) {
	maxPossibleScore := s.score + c.maxFlow*(s.time-1)
	if maxPossibleScore <= *maxScore {
		return
	}

	waitScore := s.score + s.time*s.flow
	if waitScore > *maxScore {
		*maxScore = waitScore
	}

	for unopenedId := range s.unopened {
		if ns := c.openValve(s, unopenedId); ns != nil {
			c.recursiveOpen(ns, maxScore)
		}
	}
}

func (c *cave) openValves() int {
	maxTime := 26

	manState := NewState(maxTime, c.start.id, c.goodValves)
	memo := make(map[uint64]*state)
	c.recursiveFill(manState, memo)

	// Multithreaded elephant search to speed things up
	stateCh := make(chan *state, workerCount)
	go func(sCh chan<- *state) {
		for _, s := range memo {
			s.time = maxTime
			s.currId = c.start.id
			s.flow = 0
			sCh <- s
		}
		close(sCh)
	}(stateCh)

	var wg sync.WaitGroup
	maxScores := make([]int, workerCount)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(sCh <-chan *state, maxScore *int) {
			defer wg.Done()
			for s := range sCh {
				c.recursiveOpen(s, maxScore)
			}
		}(stateCh, &maxScores[i])
	}
	wg.Wait()

	maxScore := 0
	for _, s := range maxScores {
		if s > maxScore {
			maxScore = s
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

func parseInput(input []string) cave {
	tunnels := make(map[string][]string)
	valves := make(map[string]*valve)
	for i, line := range input {
		v := valve{id: i}
		if n, err := fmt.Sscanf(line, "Valve %s has flow rate=%d", &v.name, &v.flow); err != nil {
			panic(err)
		} else if n != 2 {
			panic("Failed")
		}

		split := strings.Split(line, " to valve")
		tunnels[v.name] = strings.Split(strings.Trim(split[1], "s "), ", ")

		valves[v.name] = &v
	}

	return NewCave(valves, tunnels)
}

func main() {
	input := readInput()
	caverns := parseInput(input)
	score := caverns.openValves()
	fmt.Println(score)
}
