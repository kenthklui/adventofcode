package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type monkey interface {
	value() int
	setMonkeys(monkeys map[string]monkey)
}

type opMonkey struct {
	leftStr, rightStr string
	op                string
	left, right       monkey
}

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
func (om *opMonkey) setMonkeys(monkeys map[string]monkey) {
	om.left = monkeys[om.leftStr]
	om.right = monkeys[om.rightStr]
}

type valMonkey struct {
	val int
}

func (vm *valMonkey) value() int {
	return vm.val
}
func (vm *valMonkey) setMonkeys(monkeys map[string]monkey) {}

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

func parseInput(input []string) map[string]monkey {
	monkeys := make(map[string]monkey)
	for _, line := range input {
		splits := strings.Split(line, " ")
		name := strings.TrimRight(splits[0], ":")
		if len(splits) == 2 {
			value, _ := strconv.Atoi(splits[1])
			monkeys[name] = &valMonkey{value}
		} else {
			monkeys[name] = &opMonkey{
				leftStr:  splits[1],
				rightStr: splits[3],
				op:       splits[2],
			}
		}
	}
	for _, monkey := range monkeys {
		monkey.setMonkeys(monkeys)
	}

	return monkeys
}

func main() {
	input := readInput()
	monkeys := parseInput(input)
	fmt.Println(monkeys["root"].value())
}
