package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getALUIndex(s string) int {
	switch s {
	case "w":
		return 0
	case "x":
		return 1
	case "y":
		return 2
	case "z":
		return 3
	default:
		err := fmt.Errorf("Invalid variable \"%s\"", s)
		panic(err)
	}
}

type instruction struct {
	command         func(*ALU, int, int)
	variable, param int
}

func NewInstruction(s []string) instruction {
	varIndex := getALUIndex(s[1])
	if len(s) == 2 { // inp
		return instruction{inpVariable, varIndex, 0}
	}

	if param, err := strconv.Atoi(s[2]); err == nil {
		switch s[0] {
		case "add":
			return instruction{addValue, varIndex, param}
		case "mul":
			return instruction{mulValue, varIndex, param}
		case "div":
			return instruction{divValue, varIndex, param}
		case "mod":
			return instruction{modValue, varIndex, param}
		case "eql":
			return instruction{eqlValue, varIndex, param}
		default:
			panic("invalid ALU instruction")
		}
	} else {
		paramIndex := getALUIndex(s[2])
		switch s[0] {
		case "add":
			return instruction{addVariable, varIndex, paramIndex}
		case "mul":
			return instruction{mulVariable, varIndex, paramIndex}
		case "div":
			return instruction{divVariable, varIndex, paramIndex}
		case "mod":
			return instruction{modVariable, varIndex, paramIndex}
		case "eql":
			return instruction{eqlVariable, varIndex, paramIndex}
		default:
			panic("invalid ALU instruction")
		}
	}
}

type ALUProgram struct {
	instructions []instruction
	// At the start of next program, a variable is overwritten with inp
	// Store this for deduplication purposes
	overwriteIndex int
}

func NewALUProgram() *ALUProgram {
	return &ALUProgram{[]instruction{}, -1}
}

func (ap *ALUProgram) addInstruction(s []string) {
	ap.instructions = append(ap.instructions, NewInstruction(s))
}

type ALU struct {
	val []int
}

func NewALU() *ALU { return &ALU{[]int{0, 0, 0, 0}} }
func (alu *ALU) makeCopy() *ALU {
	aluCopy := ALU{make([]int, len(alu.val))}
	copy(aluCopy.val, alu.val)
	return &aluCopy
}

func inpVariable(alu *ALU, varIndex, param int)      { alu.val[varIndex] = param }
func addVariable(alu *ALU, varIndex, paramIndex int) { alu.val[varIndex] += alu.val[paramIndex] }
func addValue(alu *ALU, varIndex, param int)         { alu.val[varIndex] += param }
func mulVariable(alu *ALU, varIndex, paramIndex int) { alu.val[varIndex] *= alu.val[paramIndex] }
func mulValue(alu *ALU, varIndex, param int)         { alu.val[varIndex] *= param }
func divVariable(alu *ALU, varIndex, paramIndex int) { alu.val[varIndex] /= alu.val[paramIndex] }
func divValue(alu *ALU, varIndex, param int)         { alu.val[varIndex] /= param }
func modVariable(alu *ALU, varIndex, paramIndex int) { alu.val[varIndex] %= alu.val[paramIndex] }
func modValue(alu *ALU, varIndex, param int)         { alu.val[varIndex] %= param }
func eqlVariable(alu *ALU, varIndex, paramIndex int) {
	if alu.val[varIndex] == alu.val[paramIndex] {
		alu.val[varIndex] = 1
	} else {
		alu.val[varIndex] = 0
	}
}
func eqlValue(alu *ALU, varIndex, param int) {
	if alu.val[varIndex] == param {
		alu.val[varIndex] = 1
	} else {
		alu.val[varIndex] = 0
	}
}

func (alu *ALU) execute(ap *ALUProgram, inp int) {
	inpIns := ap.instructions[0]
	inpIns.command(alu, inpIns.variable, inp)

	for _, ins := range ap.instructions[1:] {
		ins.command(alu, ins.variable, ins.param)
	}
}

type ALUState struct {
	w, x, y, z int
}

func (alu *ALU) state(overwriteIndex int) ALUState {
	if overwriteIndex == -1 {
		return ALUState{alu.val[0], alu.val[1], alu.val[2], alu.val[3]}
	}

	var temp int

	alu.val[overwriteIndex], temp = temp, alu.val[overwriteIndex]
	als := ALUState{alu.val[0], alu.val[1], alu.val[2], alu.val[3]}
	alu.val[overwriteIndex], temp = temp, alu.val[overwriteIndex]

	return als
}

var empty struct{}

type Memoizer struct {
	cache []map[ALUState]struct{}
}

func NewMemoizer(length int) *Memoizer {
	cache := make([]map[ALUState]struct{}, length)
	for i := range cache {
		cache[i] = make(map[ALUState]struct{})
	}

	return &Memoizer{cache}
}

func (m *Memoizer) check(step int, als ALUState) bool {
	_, ok := m.cache[step][als]
	return ok
}

// Set something as a dead end
func (m *Memoizer) markDead(step int, als ALUState) {
	m.cache[step][als] = empty
}

func numString(digits []int) string {
	var b strings.Builder
	for _, d := range digits {
		b.WriteRune(rune('0' + d))
	}
	return b.String()
}

type modelNumFinder struct {
	aps        []*ALUProgram
	memoizer   *Memoizer
	completion int
}

func NewModelNumFinder(aps []*ALUProgram) *modelNumFinder {
	return &modelNumFinder{
		aps:        aps,
		memoizer:   NewMemoizer(len(aps)),
		completion: 0,
	}
}

func (mnf *modelNumFinder) find() (string, error) {
	success, digits := mnf.recursiveFind(NewALU(), 0)
	if success {
		return numString(digits), nil
	} else {
		return "", fmt.Errorf("Failed to find model number")
	}
}

func (mnf *modelNumFinder) recursiveFind(alu *ALU, step int) (bool, []int) {
	if step == len(mnf.aps) {
		if alu.val[3] == 0 {
			return true, []int{}
		} else {
			return false, nil
		}
	}

	als := alu.state(mnf.aps[step].overwriteIndex)
	if mnf.memoizer.check(step, als) {
		return false, nil
	}

	for nextDigit := 1; nextDigit <= 9; nextDigit++ {
		nextAlu := alu.makeCopy()
		nextAlu.execute(mnf.aps[step], nextDigit)

		if success, digits := mnf.recursiveFind(nextAlu, step+1); success {
			return true, append([]int{nextDigit}, digits...)
		}
	}

	mnf.memoizer.markDead(step, als)

	if step == 2 {
		mnf.completion++
		// fmt.Printf("Finished searching %d%% of search space\n", mnf.completion)
	}

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
				ap.overwriteIndex = getALUIndex(splits[1])
			}

			ap = NewALUProgram()
			aps = append(aps, ap)
		}

		ap.addInstruction(splits)
	}

	return aps
}

func main() {
	input := readInput()
	aps := parseInput(input)

	finder := NewModelNumFinder(aps)
	if result, err := finder.find(); err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
}
