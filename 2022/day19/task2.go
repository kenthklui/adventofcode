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

func uint8Min(a, b uint8) uint8 {
	if a < b {
		return a
	} else {
		return b
	}
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
	time, nextBuild uint8
}

func NewInventory() *inventory {
	return &inventory{
		resources: [4]uint8{0, 0, 0, 0},
		bots:      [4]uint8{1, 0, 0, 0},
		time:      0,
	}
}

func (inv *inventory) dup() *inventory {
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

func (inv *inventory) build(bp *blueprint) {
	for i, cost := range bp.costs[inv.nextBuild] {
		inv.resources[i] -= cost
	}
	inv.bots[inv.nextBuild]++
}

func (inv *inventory) minute(bp *blueprint) bool {
	canBuild := (inv.checkAfford(bp, inv.nextBuild) == 1)
	inv.gather()
	if canBuild {
		inv.build(bp)
	}
	inv.pruneExcessResources(bp)
	inv.time++

	return canBuild
}

func (inv *inventory) withBuildTarget(nextBot uint8) *inventory {
	newInv := inv.dup()
	newInv.nextBuild = nextBot
	return newInv
}

func (inv *inventory) pruneExcessResources(bp *blueprint) {
	for i := 0; i < 3; i++ {
		if inv.bots[i] >= bp.maxCosts[i] {
			inv.resources[i] = uint8Min(inv.resources[i], bp.maxCosts[i])
		}
	}
}

func (inv inventory) buildOptions(bp *blueprint) []uint8 {
	// Always build geode bots when you can
	options := make([]uint8, 0)
	if affordGeode := inv.checkAfford(bp, 3); affordGeode >= 0 {
		options = append(options, 3)
		if affordGeode == 1 {
			return options
		}
	}

	// Build a bot only when we can use more of the resource, and can afford it in the future
	for botIndex := uint8(0); botIndex < 3; botIndex++ {
		if inv.canUseBot(bp, botIndex) && inv.checkAfford(bp, botIndex) >= 0 {
			options = append(options, botIndex)
		}
	}

	return options
}

func (inv inventory) canUseBot(bp *blueprint, botIndex uint8) bool {
	timeRemains := int(duration - inv.time - 1)
	shortfall := timeRemains * (int(inv.bots[botIndex]) - int(bp.maxCosts[botIndex]))
	projection := shortfall + int(inv.resources[botIndex])
	return projection < 0
}

// 1: Can buy now, 0: Can buy if we wait, -1: Can't buy without building another bot first
func (inv inventory) checkAfford(bp *blueprint, bot uint8) int8 {
	afford := 0
	future := 0
	for i, cost := range bp.costs[bot] {
		if inv.resources[i] >= cost {
			afford++
		} else if inv.bots[i] > 0 {
			future++
		}
	}
	if afford == 3 {
		return 1
	} else if afford+future == 3 {
		return 0
	} else {
		return -1
	}
}

type invkey uint64
type void struct{}
type memo map[invkey]void

var empty void

func (inv *inventory) key(bp *blueprint) invkey {
	k := uint64(inv.bots[3])<<8 + uint64(inv.resources[3])
	k = k*4 + uint64(inv.nextBuild)
	for i, c := range bp.maxCosts {
		k *= uint64(c + 1)
		k += uint64(inv.bots[i])
		k *= uint64(c*2 + 1)
		k += uint64(inv.resources[i])
	}
	return invkey(k)
}

func maxGeodes(bp *blueprint, minutes uint8) int {
	m := make(memo, 1<<18)
	queue := list.New()

	inv := NewInventory()
	queue.PushBack(inv.withBuildTarget(0))
	queue.PushBack(inv.withBuildTarget(1))

	var maxGeodes uint8
	var next *list.Element
	for e := queue.Front(); e != nil; e = next {
		inv = e.Value.(*inventory)
		botBuilt := inv.minute(bp)

		if inv.time == minutes-1 {
			geodes := inv.resources[3] + inv.bots[3]
			if maxGeodes < geodes {
				maxGeodes = geodes
			}
		} else {
			key := inv.key(bp)
			if _, ok := m[key]; !ok {
				m[key] = empty
				if botBuilt {
					for _, nextBot := range inv.buildOptions(bp) {
						queue.PushBack(inv.withBuildTarget(nextBot))
					}
				} else {
					queue.PushBack(inv)
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

const duration = 32

func main() {
	input := readInput()
	bps := parseInput(input)
	if len(bps) > 3 {
		bps = bps[:3]
	}

	maxGeodeProduct := 1
	for _, bp := range bps {
		geodes := maxGeodes(bp, duration)
		maxGeodeProduct *= geodes
	}
	fmt.Println(maxGeodeProduct)
}
