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
	re := regexp.MustCompile(`([0-9]+) \+ ([0-9]+)`)
	for m := re.FindStringSubmatch(line); m != nil; m = re.FindStringSubmatch(line) {
		first, _ := strconv.Atoi(m[1])
		second, _ := strconv.Atoi(m[2])
		strValue := strconv.Itoa(first + second)
		line = strings.Replace(line, m[0], strValue, 1)
	}
	re = regexp.MustCompile(`([0-9]+) \* ([0-9]+)`)
	for m := re.FindStringSubmatch(line); m != nil; m = re.FindStringSubmatch(line) {
		first, _ := strconv.Atoi(m[1])
		second, _ := strconv.Atoi(m[2])
		strValue := strconv.Itoa(first * second)
		line = strings.Replace(line, m[0], strValue, 1)
	}
	value, _ := strconv.Atoi(line)

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
