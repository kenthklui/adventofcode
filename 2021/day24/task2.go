package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type instruction struct {
	command, variable, param string
}

type ALUProgram struct {
	instructions []instruction
	// At the start of next program, a variable is overwritten with inp
	// Store this for deduplication purposes
	overwriteVariable string
}

func NewALUProgram(inpVar string) *ALUProgram {
	ins := instruction{command: "inp", variable: inpVar}
	ap := ALUProgram{[]instruction{ins}, ""}
	return &ap
}

func (ap *ALUProgram) addInstruction(s []string) {
	ap.instructions = append(ap.instructions, instruction{s[0], s[1], s[2]})
}

type ALU struct {
	w, x, y, z int
}

func NewALU() *ALU { return &ALU{} }

func (alu *ALU) makeCopy() *ALU {
	return &ALU{alu.w, alu.x, alu.y, alu.z}
}

func (alu *ALU) getVar(s string) *int {
	switch s {
	case "w":
		return &alu.w
	case "x":
		return &alu.x
	case "y":
		return &alu.y
	case "z":
		return &alu.z
	default:
		err := fmt.Errorf("Invalid variable \"%s\"", s)
		panic(err)
	}
}

func (alu *ALU) get(s string) *int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return &i
	} else {
		return alu.getVar(s)
	}
}

func (alu *ALU) execute(ap *ALUProgram, inp uint8) {
	for _, ins := range ap.instructions {
		v := alu.getVar(ins.variable)

		if ins.command == "inp" {
			*v = int(inp)
			continue
		}

		p := alu.get(ins.param)
		switch ins.command {
		case "add":
			*v += *p
		case "mul":
			*v *= *p
		case "div":
			/*
				if *p == 0 {
					panic("div by 0")
				}
			*/
			*v /= *p
		case "mod":
			/*
				if *v < 0 {
					panic("mod negative")
				} else if *p <= 0 {
					panic("mod by non-positive")
				}
			*/
			*v %= *p
		case "eql":
			if *v == *p {
				*v = 1
			} else {
				*v = 0
			}
		default:
			panic("invalid ALU instruction")
		}
	}
}

var empty struct{}

type Memoizer struct {
	cache []map[ALU]struct{}
}

func NewMemoizer(length int) *Memoizer {
	cache := make([]map[ALU]struct{}, length)
	for i := range cache {
		cache[i] = make(map[ALU]struct{})
	}

	return &Memoizer{cache}
}

func (m *Memoizer) check(step int, alu *ALU) bool {
	_, ok := m.cache[step][*alu]
	return ok
}

// Set something as a dead end
func (m *Memoizer) kill(step int, alu *ALU) {
	m.cache[step][*alu] = empty
}

func numString(digits []uint8) string {
	var b strings.Builder
	for _, d := range digits {
		b.WriteRune(rune('0' + d))
	}
	return b.String()
}

func findModelNum(aps []*ALUProgram) string {
	memoizer := NewMemoizer(len(aps))

	alu := NewALU()
	success, digits := recursiveFindModelNum(aps, memoizer, alu, 0)
	if success {
		return numString(digits)
	} else {
		return "Failed to find model number"
	}
}

func recursiveFindModelNum(aps []*ALUProgram, memoizer *Memoizer, alu *ALU, step int) (bool, []uint8) {
	if step == len(aps) {
		if alu.z == 0 {
			return true, []uint8{}
		} else {
			return false, nil
		}
	}

	checkAlu := alu.makeCopy()
	if aps[step].overwriteVariable != "" {
		overwriteVar := checkAlu.getVar(aps[step].overwriteVariable)
		*overwriteVar = 0
	}

	if memoizer.check(step, checkAlu) {
		return false, nil
	}

	for nextDigit := uint8(1); nextDigit <= 9; nextDigit++ {
		nextAlu := alu.makeCopy()
		nextAlu.execute(aps[step], nextDigit)

		success, digits := recursiveFindModelNum(aps, memoizer, nextAlu, step+1)
		if success {
			return true, append([]uint8{nextDigit}, digits...)
		}
	}

	memoizer.kill(step, checkAlu)

	return false, nil
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

func parseInput(input []string) []*ALUProgram {
	// Break up the program at each inp command
	aps := make([]*ALUProgram, 0, 14)

	var ap *ALUProgram
	for _, line := range input {
		splits := strings.Split(line, " ")
		if splits[0] == "inp" {
			if ap != nil {
				ap.overwriteVariable = splits[1]
			}

			ap = NewALUProgram(splits[1])
			aps = append(aps, ap)
		} else {
			ap.addInstruction(splits)
		}
	}

	return aps
}

func main() {
	input := readInput()
	aps := parseInput(input)

	modelNum := findModelNum(aps)
	fmt.Println(modelNum)
}
