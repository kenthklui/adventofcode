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
	value, err := strconv.Atoi(tokens[0])
	if err != nil {
		panic(err)
	}

	for i := 1; i < len(tokens); i += 2 {
		nextValue, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			panic(err)
		}

		switch tokens[i] {
		case "+":
			value += nextValue
		case "*":
			value *= nextValue
		default:
			panic("Invalid token")
		}
	}

	return value
}

func solveLine(line string) int {
	bracketRegex := regexp.MustCompile(`\([^\(\)]+\)`)
	for strings.Index(line, "(") != -1 {
		newLine := line
		submatches := bracketRegex.FindAllStringSubmatchIndex(line, -1)

		for _, match := range submatches {
			exp := line[match[0]+1 : match[1]-1]
			strValue := strconv.Itoa(solveLine(exp))
			newLine = strings.ReplaceAll(newLine, line[match[0]:match[1]], strValue)
		}

		line = newLine
	}

	addRegex := regexp.MustCompile(`[0-9]+ \+ [0-9]+`)
	for strings.Index(line, "+") != -1 {
		submatch := addRegex.FindString(line)

		strValue := strconv.Itoa(solveBasic(submatch))
		line = strings.Replace(line, submatch, strValue, 1)
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
