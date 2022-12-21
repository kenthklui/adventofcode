package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type monkey interface {
	name() string
	value() int
	setParent(monkey)
	setChildren(monkeys map[string]monkey)
	parentMonkey() monkey
	// Given a goal value, and indicating the mutable child, what value should the child have?
	requirement(int, monkey) int
}

type opMonkey struct {
	nameStr               string
	leftStr, rightStr, op string
	parent, left, right   monkey
}

func (om *opMonkey) name() string { return om.nameStr }
func (om *opMonkey) value() int {
	switch om.op {
	case "+":
		return om.left.value() + om.right.value()
	case "-":
		return om.left.value() - om.right.value()
	case "*":
		return om.left.value() * om.right.value()
	case "/":
		return om.left.value() / om.right.value()
	default:
		panic("Invalid op")
	}
}
func (om *opMonkey) setParent(m monkey) {
	if om.parent != nil {
		panic("Multiple parents")
	}
	om.parent = m
}
func (om *opMonkey) setChildren(monkeys map[string]monkey) {
	om.left = monkeys[om.leftStr]
	om.left.setParent(om)
	om.right = monkeys[om.rightStr]
	om.right.setParent(om)
}
func (om *opMonkey) parentMonkey() monkey { return om.parent }

func (om *opMonkey) requirement(goal int, child monkey) int {
	if child.name() == om.left.name() {
		switch om.op {
		case "+":
			return goal - om.right.value()
		case "-":
			return goal + om.right.value()
		case "*":
			return goal / om.right.value()
		case "/":
			return goal * om.right.value()
		case "=":
			return om.right.value()
		default:
			panic("Invalid op")
		}
	} else if child.name() == om.right.name() {
		switch om.op {
		case "+":
			return goal - om.left.value()
		case "-":
			return om.left.value() - goal
		case "*":
			return goal / om.left.value()
		case "/":
			return om.left.value() / goal
		case "=":
			return om.left.value()
		default:
			panic("Invalid op")
		}
	}
	fmt.Println(child.name(), om.left.name(), om.right.name())
	panic("Invalid monkey child")
}

type valMonkey struct {
	nameStr string
	val     int
	parent  monkey
}

func (vm *valMonkey) name() string { return vm.nameStr }
func (vm *valMonkey) value() int   { return vm.val }
func (vm *valMonkey) setParent(m monkey) {
	if vm.parent != nil {
		panic("Multiple parents")
	}
	vm.parent = m
}
func (vm *valMonkey) setChildren(monkeys map[string]monkey) {}
func (vm *valMonkey) parentMonkey() monkey                  { return vm.parent }

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
func (vm *valMonkey) requirement(goal int, child monkey) int {
	return goal
}

func parseInput(input []string) map[string]monkey {
	monkeys := make(map[string]monkey)
	for _, line := range input {
		splits := strings.Split(line, " ")
		name := strings.TrimRight(splits[0], ":")
		if len(splits) == 2 {
			value, _ := strconv.Atoi(splits[1])
			monkeys[name] = &valMonkey{nameStr: name, val: value}
		} else {
			m := opMonkey{
				nameStr:  name,
				leftStr:  splits[1],
				rightStr: splits[3],
				op:       splits[2],
			}
			if name == "root" {
				m.op = "="
			}
			monkeys[name] = &m
		}
	}
	for _, monkey := range monkeys {
		monkey.setChildren(monkeys)
	}

	return monkeys
}

func main() {
	input := readInput()
	monkeys := parseInput(input)

	// Travel up the chain to find the path to root
	human := monkeys["humn"]
	path := make([]monkey, 0)
	for curr := human; curr.parentMonkey() != nil; curr = curr.parentMonkey() {
		path = append(path, curr)
	}

	// Travel down the chain to reverse engineer the required value
	curr := monkeys["root"]
	goal := 0
	for i := len(path) - 1; i >= 0; i-- {
		goal = curr.requirement(goal, path[i])
		curr = path[i]
	}

	fmt.Println(goal)
}
