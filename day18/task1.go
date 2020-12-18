package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func solveBasic(line string) int {
	tokens := strings.Split(line, " ")
	value, _ := strconv.Atoi(tokens[0])
	for i := 1; i < len(tokens); i += 2 {
		if nextValue, err := strconv.Atoi(tokens[i+1]); err == nil {
			switch tokens[i] {
			case "+":
				value += nextValue
			case "*":
				value *= nextValue
			}
		} else {
			panic(err)
		}
	}

	return value
}

func solveBasicStr(line string) string {
	line = strings.TrimPrefix(line, "(")
	line = strings.TrimSuffix(line, ")")
	return strconv.Itoa(solveBasic(line))
}

func solveLine(line string) int {
	re := regexp.MustCompile(`\(([^\(\)]+)\)`)
	for re.MatchString(line) {
		line = re.ReplaceAllStringFunc(line, solveBasicStr)
	}

	return solveBasic(line)
}

func main() {
	input := readInput()
	sum := 0
	for _, line := range input {
		sum += solveLine(line)
	}

	fmt.Println(sum)
}
