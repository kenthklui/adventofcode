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

type state struct {
	Acc, Current int
}

type machine struct {
	Ops      []op
	Executed []bool
	State    state
}

func NewMachine(ops []op) machine {
	return machine{
		Ops:      ops,
		Executed: make([]bool, len(ops)),
		State:    state{0, 0},
	}
}

func (m *machine) ExecuteWithoutRepeat() bool {
	if m.Executed[m.State.Current] {
		return false
	} else {
		m.Executed[m.State.Current] = true
	}

	currentOp := m.Ops[m.State.Current]

	switch currentOp.Operation {
	case "nop":
		m.State.Current++
	case "acc":
		m.State.Acc += currentOp.Value
		m.State.Current++
	case "jmp":
		m.State.Current += currentOp.Value
	}

	return true
}

func (m *machine) Run() (int, error) {
	for {
		if m.ExecuteWithoutRepeat() {
			if m.State.Current == len(m.Ops) {
				return m.State.Acc, nil
			}
		} else {
			break
		}
	}
	return 0, fmt.Errorf("Failed to terminate, repeated instruction #%d", m.State.Current)
}

func (m *machine) Reset() {
	for i := 0; i < len(m.Executed); i++ {
		m.Executed[i] = false
	}

	m.State.Acc = 0
	m.State.Current = 0
}

func (m *machine) TryFixes() {
	for i := 0; i < len(m.Ops); i++ {
		currentOp := m.Ops[i]

		switch currentOp.Operation {
		case "nop":
			m.Ops[i].Operation = "jmp"

			if result, err := m.Run(); err == nil {
				log.Printf("Succeeded fixing %d, accumulator: %d\n", i, result)
				return
			}

			m.Ops[i].Operation = "nop"
			m.Reset()
		case "acc":
			// Do nothing
		case "jmp":
			m.Ops[i].Operation = "nop"

			if result, err := m.Run(); err == nil {
				log.Printf("Succeeded fixing %d, accumulator: %d\n", i, result)
				return
			}

			m.Ops[i].Operation = "jmp"
			m.Reset()
		}
	}

	log.Fatal("Failed to find a fix")
}

func runOps() {
	ops := make([]op, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var operation string
		var value int

		n, err := fmt.Sscanf(scanner.Text(), "%s %d", &operation, &value)
		if err != nil || n != 2 {
			log.Fatalf("Failed to parse line: %q", scanner.Text())
		}

		ops = append(ops, op{false, operation, value})
	}

	m := NewMachine(ops)
	m.TryFixes()
}

func main() {
	runOps()
}
