package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func add(a, b int) int    { return a + b }
func mul(a, b int) int    { return a * b }
func square(a, b int) int { return a * a }

type operation struct {
	op    func(int, int) int
	param int
}

type monkey struct {
	items     []int
	inspected int

	op                  operation
	trueDest, falseDest int
	divisibleTest       int
}

func (m *monkey) inspectItems(divisor int) ([]int, []int) {
	trueItems, falseItems := []int{}, []int{}
	for _, item := range m.items {
		item = m.op.op(item, m.op.param)
		item %= divisor

		m.inspected++

		if item%m.divisibleTest == 0 {
			trueItems = append(trueItems, item)
		} else {
			falseItems = append(falseItems, item)
		}
	}
	m.items = nil

	return trueItems, falseItems
}

type monkeyPack struct {
	monkeys []*monkey
	gcd     int
}

func (mp *monkeyPack) divisor() int {
	if mp.gcd == 0 {
		mp.gcd = 1
		for _, m := range mp.monkeys {
			mp.gcd *= m.divisibleTest
		}
	}

	return mp.gcd
}

func (mp *monkeyPack) round() {
	for _, m := range mp.monkeys {
		trueItems, falseItems := m.inspectItems(mp.divisor())

		trueDest, falseDest := &mp.monkeys[m.trueDest].items, &mp.monkeys[m.falseDest].items
		*trueDest = append(*trueDest, trueItems...)
		*falseDest = append(*falseDest, falseItems...)
	}
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

func parseInput(input []string) *monkeyPack {
	mp := monkeyPack{make([]*monkey, 0), 0}
	for i := 0; i < len(input); i += 7 {
		itemStr := strings.TrimPrefix(input[i+1], "  Starting items: ")
		itemTokens := strings.Split(itemStr, ", ")
		items := make([]int, 0, len(itemTokens))
		for _, it := range itemTokens {
			item, _ := strconv.Atoi(it)
			items = append(items, item)
		}

		var op operation
		opStr := strings.TrimPrefix(input[i+2], "  Operation: new = ")
		opTokens := strings.Split(opStr, " ")
		switch opTokens[1] {
		case "*":
			if opTokens[2] == "old" {
				op.op = square
			} else {
				op.op = mul
				op.param, _ = strconv.Atoi(opTokens[2])
			}
		case "+":
			op.op = add
			op.param, _ = strconv.Atoi(opTokens[2])
		default:
			fmt.Println(opStr)
		}

		var divisibleTest int
		if _, err := fmt.Sscanf(input[i+3], "  Test: divisible by %d", &divisibleTest); err != nil {
			panic(err)
		}

		var trueDest, falseDest int
		if _, err := fmt.Sscanf(input[i+4], "    If true: throw to monkey %d", &trueDest); err != nil {
			panic(err)
		}
		if _, err := fmt.Sscanf(input[i+5], "    If false: throw to monkey %d", &falseDest); err != nil {
			panic(err)
		}
		mp.monkeys = append(mp.monkeys, &monkey{items, 0, op, trueDest, falseDest, divisibleTest})
	}

	return &mp
}

func main() {
	input := readInput()
	monkeys := parseInput(input)

	for i := 0; i < 10000; i++ {
		monkeys.round()
	}
	activity := make([]int, len(monkeys.monkeys))
	for _, m := range monkeys.monkeys {
		activity = append(activity, m.inspected)
	}

	sort.Ints(activity)
	n := len(activity) - 1
	fmt.Println(activity[n] * activity[n-1])
}
