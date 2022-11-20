package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type op struct {
	Executed  bool
	Operation string
	Value     int
}

type machine struct {
	Ops          []op
	Acc, Current int
}

func (m *machine) Execute(stopOnRepeat bool) bool {
	currentOp := m.Ops[m.Current]

	if stopOnRepeat && currentOp.Executed {
		return false
	} else {
		m.Ops[m.Current].Executed = true
	}

	switch currentOp.Operation {
	case "nop":
		m.Current++
	case "acc":
		m.Acc += currentOp.Value
		m.Current++
	case "jmp":
		m.Current += currentOp.Value
	}

	return true
}

func runOps() int {
	ops := make([]op, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var operation string
		var value int

		n, err := fmt.Sscanf(scanner.Text(), "%s %d", &operation, &value)
		if err != nil || n != 2 {
			log.Fatal("Failed to parse line: %q", scanner.Text())
		}

		ops = append(ops, op{false, operation, value})
	}

	m := machine{ops, 0, 0}

	for {
		if !m.Execute(true) {
			break
		}
	}

	return m.Acc
}

func main() {
	fmt.Println(runOps())
}
