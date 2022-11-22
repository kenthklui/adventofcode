package main

import (
	"bufio"
	"fmt"
	"os"
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

func parseInput(input []string) (string, map[string]string) {
	template := input[0]
	rules := make(map[string]string)

	var pair, insert string
	for _, line := range input[2:] {
		if _, err := fmt.Sscanf(line, "%s -> %s", &pair, &insert); err == nil {
			rules[pair] = insert
		} else {
			panic(err)
		}
	}

	return template, rules
}

func polymerStep(template string, rules map[string]string) string {
	var newTemplate strings.Builder

	// First char never changes
	if err := newTemplate.WriteByte(template[0]); err != nil {
		panic(err)
	}

	for i, r := range template[1:] {
		pair := template[i : i+2]
		if insert, ok := rules[pair]; ok {
			if _, err := newTemplate.WriteString(insert); err != nil {
				panic(err)
			}
		}

		if _, err := newTemplate.WriteRune(r); err != nil {
			panic(err)
		}
	}

	return newTemplate.String()
}

func elementDiff(template string) int {
	occurrence := make(map[rune]int)

	for _, r := range template {
		if _, ok := occurrence[r]; ok {
			occurrence[r]++
		} else {
			occurrence[r] = 1
		}
	}

	var min, max int
	for _, v := range occurrence {
		if min == 0 || min > v {
			min = v
		}
		if max == 0 || max < v {
			max = v
		}
	}

	return max - min
}

func main() {
	input := readInput()
	template, rules := parseInput(input)

	s := template
	for i := 0; i < 10; i++ {
		s = polymerStep(s, rules)
	}

	fmt.Println(elementDiff(s))
}
