package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// {ore, clay, obsidian, geodes}
type costList [3]uint8
type blueprint struct {
	costs    [4]costList
	maxCosts [3]uint8
}

func uint8Max(a, b uint8) uint8 {
	if a > b {
		return a
	} else {
		return b
	}
}

func NewBlueprint(params ...uint8) *blueprint {
	var bp blueprint

	bp.costs[0][0] = params[0]
	bp.costs[1][0] = params[1]
	bp.costs[2][0] = params[2]
	bp.costs[2][1] = params[3]
	bp.costs[3][0] = params[4]
	bp.costs[3][2] = params[5]

	bp.maxCosts[0] = uint8Max(uint8Max(uint8Max(params[0], params[1]), params[2]), params[4])
	bp.maxCosts[1] = params[3]
	bp.maxCosts[2] = params[5]

	return &bp
}

type inventory struct {
	// ore, clay, obsidian, geode
	resources, bots [4]uint8
	time            uint8
}

func NewInventory() *inventory {
	return &inventory{
		resources: [4]uint8{0, 0, 0, 0},
		bots:      [4]uint8{1, 0, 0, 0},
		time:      0,
	}
}

func (inv inventory) dup() *inventory {
	return &inventory{
		resources: inv.resources,
		bots:      inv.bots,
		time:      inv.time,
	}
}

func (inv *inventory) gather() {
	for i, b := range inv.bots {
		inv.resources[i] += b
	}
}

func (inv *inventory) build(bp *blueprint, buildOrder int) {
	for i, cost := range bp.costs[buildOrder] {
		inv.resources[i] -= cost
	}
	inv.bots[buildOrder]++
}

func (inv *inventory) minute(bp *blueprint, buildOrder int) *inventory {
	newInv := inv.dup()
	newInv.gather()
	if buildOrder >= 0 {
		newInv.build(bp, buildOrder)
	}
	newInv.pruneExcessResources(bp)
	newInv.time++
	return newInv
}

func (inv *inventory) pruneExcessResources(bp *blueprint) {
	for i := 0; i < 3; i++ {
		if inv.bots[i] >= bp.maxCosts[i] && inv.resources[i] > bp.maxCosts[i] {
			inv.resources[i] = bp.maxCosts[i]
		}
	}
}

func (inv inventory) buildOptions(bp *blueprint) []int {
	// Always build geode bots when you can
	if inv.checkAfford(bp, 3) {
		return []int{3}
	}

	options := make([]int, 0)
	// Build a bot only when we can use more of the resource, and can afford it
	for botIndex := 0; botIndex < 3; botIndex++ {
		if inv.bots[botIndex] < bp.maxCosts[botIndex] && inv.checkAfford(bp, botIndex) {
			options = append(options, botIndex)
		}
	}

	// Save only if something was unaffordable
	if len(options) < 3 {
		options = append(options, -1)
	}

	return options
}

func (inv inventory) checkAfford(bp *blueprint, botIndex int) bool {
	for i, cost := range bp.costs[botIndex] {
		if inv.resources[i] < cost {
			return false
		}
	}
	return true
}

type invkey uint64
type void struct{}
type memo map[invkey]void

var empty void

func (inv *inventory) key(bp *blueprint) invkey {
	k := uint64(inv.bots[3])<<8 + uint64(inv.resources[3])
	for i, c := range bp.maxCosts {
		k *= uint64(c + 1)
		k += uint64(inv.bots[i])
		k *= uint64(c*2 + 1)
		k += uint64(inv.resources[i])
	}
	return invkey(k)
}

func maxGeodes(bp *blueprint, minutes uint8) int {
	m := make(memo, 1<<16)
	queue := list.New()

	inv := NewInventory()
	queue.PushBack(inv)

	var maxGeodes uint8
	var next *list.Element
	for e := queue.Front(); e != nil; e = next {
		inv = e.Value.(*inventory)

		if inv.time == minutes-1 {
			geodes := inv.resources[3] + inv.bots[3]
			if maxGeodes < geodes {
				maxGeodes = geodes
			}
		} else {
			buildOptions := inv.buildOptions(bp)
			for _, buildOrder := range buildOptions {
				newInv := inv.minute(bp, buildOrder)
				key := newInv.key(bp)
				if _, ok := m[key]; !ok {
					m[key] = empty
					queue.PushBack(newInv)
				}
			}
		}

		next = e.Next()
		queue.Remove(e)
	}

	return int(maxGeodes)
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

func parseInput(input []string) []*blueprint {
	bps := make([]*blueprint, 0, len(input))

	params := make([]uint8, 6)
	tokenPos := []int{6, 12, 18, 21, 27, 30}
	for _, line := range input {
		tokens := strings.Split(line, " ")
		for i, pos := range tokenPos {
			if paramInt, err := strconv.Atoi(tokens[pos]); err == nil {
				params[i] = uint8(paramInt)
			} else {
				panic(err)
			}
		}

		bp := NewBlueprint(params...)
		bps = append(bps, bp)
	}
	return bps
}

const duration = 24

func main() {
	input := readInput()
	bps := parseInput(input)

	var qualitySum int
	for i, bp := range bps {
		id := i + 1
		geodes := maxGeodes(bp, duration)
		quality := geodes * id
		qualitySum += quality
	}
	fmt.Println(qualitySum)
}
