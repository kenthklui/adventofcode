package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type void struct{}

var nul void

type node struct {
	name      string
	edges     map[*edge]*node
	travelled bool
}

func newNode(name string) *node { return &node{name: name, edges: make(map[*edge]*node)} }

type nodeMap map[*node]void

type edge struct {
	travelled bool
}

type graph struct {
	nodes map[string]*node
	edges []*edge
}

func (g *graph) resetNodes() {
	for _, n := range g.nodes {
		n.travelled = false
	}
}

func (g *graph) resetEdges() {
	for _, e := range g.edges {
		e.travelled = false
	}
}

func (g *graph) bfs(source *node, dest *node) (bool, nodeMap) {
	type queueItem struct {
		e    *edge
		n    *node
		prev *queueItem
	}

	queue := make([]*queueItem, 0, len(g.nodes))
	queue = append(queue, &queueItem{n: source})

	reachedNodes := make(nodeMap)
	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		reachedNodes[curr.n] = nul
		if curr.n == dest {
			for itr := curr; itr.e != nil; itr = itr.prev {
				itr.e.travelled = true
			}
			found = true
			break
		}

		for e, n := range curr.n.edges {
			if e.travelled || n.travelled {
				continue
			}
			n.travelled = true
			queue = append(queue, &queueItem{e, n, curr})
		}
	}
	g.resetNodes()
	return found, reachedNodes
}

func (g *graph) cutPaths(source, dest *node, pathNum int) bool {
	complete := true
	for i := 0; i < pathNum; i++ {
		if found, _ := g.bfs(source, dest); !found {
			complete = false
			break
		}
	}
	return complete
}

func (g *graph) split(cuts int) (nodeMap, nodeMap) {
	g1, g2 := make(nodeMap), make(nodeMap)

	var source *node
	for _, n := range g.nodes {
		source = n
		break
	}
	g1[source] = nul

	for _, dest := range g.nodes {
		if source == dest {
			continue
		}
		_, ok1 := g1[dest]
		_, ok2 := g2[dest]
		if ok1 || ok2 {
			continue
		}

		if g.cutPaths(source, dest, cuts+1) {
			g1[dest] = nul
		} else {
			// Use disconnected graph to categorize as many nodes as possible
			_, g1nodes := g.bfs(source, nil)
			for n := range g1nodes {
				g1[n] = nul
			}
			_, g2nodes := g.bfs(dest, nil)
			for n := range g2nodes {
				g2[n] = nul
			}
		}
		g.resetEdges()
	}
	return g1, g2
}

func parse(input []string) *graph {
	nodes := make(map[string]*node)
	for _, line := range input {
		name, _, _ := strings.Cut(line, ": ")
		nodes[name] = newNode(name)
	}
	edges := make([]*edge, 0)
	for _, line := range input {
		sourceName, destNames, _ := strings.Cut(line, ": ")
		source := nodes[sourceName]
		for _, destName := range strings.Split(destNames, " ") {
			dest, ok := nodes[destName]
			if !ok {
				dest = newNode(destName)
				nodes[destName] = dest
			}
			newEdge := &edge{}
			edges = append(edges, newEdge)
			source.edges[newEdge] = dest
			dest.edges[newEdge] = source
		}
	}
	return &graph{nodes, edges}
}

const cutCount = 3

func main() {
	input := util.StdinReadlines()
	g := parse(input)
	g1, g2 := g.split(cutCount)
	fmt.Println(len(g1) * len(g2))
}
