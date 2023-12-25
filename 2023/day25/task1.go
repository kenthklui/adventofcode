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
	group     int
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

func (g *graph) bfs(source, dest *node) bool {
	type queueItem struct {
		e    *edge
		n    *node
		prev *queueItem
	}

	queue := make([]*queueItem, 0)
	queue = append(queue, &queueItem{n: source})

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.prev != nil {
			curr.n.group = curr.prev.n.group
		}
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
	return found
}

func (g *graph) cutPaths(source, dest *node, pathNum int) bool {
	complete := true
	for i := 0; i < pathNum; i++ {
		if !g.bfs(source, dest) {
			complete = false
			break
		}
	}
	return complete
}

func (g *graph) split(cuts int) (int, int) {
	var source *node
	for _, n := range g.nodes {
		source = n
		break
	}
	source.group = 1

	for _, dest := range g.nodes {
		if source == dest {
			continue
		}
		if dest.group > 0 {
			continue
		}

		if !g.cutPaths(source, dest, cuts+1) {
			// Use disconnected graph to categorize as many nodes as possible
			g.bfs(source, nil)
			dest.group = 2
			g.bfs(dest, nil)
		}
		g.resetEdges()
	}

	disconnected := 0
	for _, n := range g.nodes {
		if n.group != 1 {
			disconnected++
		}
	}
	return len(g.nodes) - disconnected, disconnected
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
	g1size, g2size := g.split(cutCount)
	fmt.Println(g1size * g2size)
}
