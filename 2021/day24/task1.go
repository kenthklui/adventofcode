package main

import (
	"bufio"
	"context"
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
	m            Memoizer
	// At the start of next program, a variable is overwritten with inp
	// Store this for deduplication purposes
	nextOverwrite int
}

func NewALUProgram() *ALUProgram {
	return &ALUProgram{
		instructions:  make([]instruction, 0),
		m:             make(Memoizer),
		nextOverwrite: -1,
	}
}

func (ap *ALUProgram) addInstruction(s []string) {
	ap.instructions = append(ap.instructions, NewInstruction(s))
}

func (ap *ALUProgram) execute(alu *ALU, inp int) *ALU {
	nextAlu := alu.child(inp)

	inpIns := ap.instructions[0]
	inpIns.command(nextAlu, inpIns.variable, inp)

	for _, ins := range ap.instructions[1:] {
		ins.command(nextAlu, ins.variable, ins.param)
	}

	if ap.nextOverwrite >= 0 {
		nextAlu.val[ap.nextOverwrite] = 0
	}

	return nextAlu
}

type ALU struct {
	val    [4]int
	parent *ALU
	digit  int
}

func (alu *ALU) child(digit int) *ALU {
	newAlu := new(ALU)
	newAlu.val = alu.val
	newAlu.parent = alu
	newAlu.digit = digit

	return newAlu
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

type ALUKey [3]int

func (alu ALU) key() ALUKey {
	return [3]int{alu.val[0] - alu.val[1], alu.val[1] - alu.val[2], alu.val[2] - alu.val[3]}
}

type void struct{}
type Memoizer map[ALUKey]void

var empty void

func numString(digits []int) string {
	var b strings.Builder
	for _, d := range digits {
		b.WriteRune(rune('0' + d))
	}
	return b.String()
}

type modelNumFinder struct {
	aps []*ALUProgram
}

func NewModelNumFinder(aps []*ALUProgram) *modelNumFinder {
	return &modelNumFinder{aps: aps}
}

func (mnf *modelNumFinder) find() (string, error) {
	steps := len(mnf.aps)

	aluChs := make([]chan *ALU, steps+1)
	aluChs[0] = make(chan *ALU)
	ctx, cancel := context.WithCancel(context.Background())
	for i, ap := range mnf.aps {
		aluChs[i+1] = make(chan *ALU, 2048)
		go chFind(aluChs[i], aluChs[i+1], ap, ctx)
	}

	aluChs[0] <- new(ALU)
	close(aluChs[0])

	for alu := range aluChs[steps] {
		if alu.val[3] == 0 {
			cancel()
			return mnf.digitString(alu), nil
		}
	}

	return "", fmt.Errorf("Failed to find model number")
}

func chFind(inCh <-chan *ALU, outCh chan<- *ALU, ap *ALUProgram, ctx context.Context) {
Loop:
	for alu := range inCh {
		for nextDigit := 9; nextDigit > 0; nextDigit-- {
			nextAlu := ap.execute(alu, nextDigit)

			key := nextAlu.key()
			if _, ok := ap.m[key]; !ok {
				ap.m[key] = empty
				select {
				case <-ctx.Done():
					break Loop
				case outCh <- nextAlu:
				}
			}
		}
	}

	close(outCh)
	for _ = range inCh {
	}
}

func (mnf *modelNumFinder) digitString(alu *ALU) string {
	generations := 0
	for ancestor := alu; ancestor.parent != nil; ancestor = ancestor.parent {
		generations++
	}

	dig := make([]int, generations)
	for curr := alu; generations > 0; curr = curr.parent {
		generations--
		dig[generations] = curr.digit
	}

	return numString(dig)
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
				ap.nextOverwrite = getALUIndex(splits[1])
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
