package main

import (
	"bufio"
	"fmt"
	"os"
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

type template struct {
	first, last string
	pairs       map[string]int
}

func parseInput(input []string) (template, map[string]string) {
	// Order in the template doesn't matter, only frequency of pair occurrence...
	templatePairs := make(map[string]int)
	for i := range input[0][1:] {
		pair := input[0][i : i+2]
		if _, ok := templatePairs[pair]; ok {
			templatePairs[pair]++
		} else {
			templatePairs[pair] = 1
		}
	}
	// ...except the first and last characters!

	temp := template{
		first: input[0][:1],
		last:  input[0][len(input[0])-1:],
		pairs: templatePairs,
	}

	rules := make(map[string]string)

	var pair, insert string
	for _, line := range input[2:] {
		if _, err := fmt.Sscanf(line, "%s -> %s", &pair, &insert); err == nil {
			rules[pair] = insert
		} else {
			panic(err)
		}
	}

	return temp, rules
}

func polymerStep(temp template, rules map[string]string) template {
	newPairs := make(map[string]int)
	for pair, v := range temp.pairs {
		if insert, ok := rules[pair]; ok {
			pair1 := pair[:1] + insert

			if _, ok := newPairs[pair1]; ok {
				newPairs[pair1] += v
			} else {
				newPairs[pair1] = v
			}

			pair2 := insert + pair[1:]

			if _, ok := newPairs[pair2]; ok {
				newPairs[pair2] += v
			} else {
				newPairs[pair2] = v
			}
		} else {
			newPairs[pair] = v
		}
	}

	return template{
		first: temp.first,
		last:  temp.last,
		pairs: newPairs,
	}

}

func elementDiff(temp template) int {
	occurrence := make(map[byte]int)

	// Observation: every pair causes characters to be counted double...
	for pair, v := range temp.pairs {
		if _, ok := occurrence[pair[0]]; ok {
			occurrence[pair[0]] += v
		} else {
			occurrence[pair[0]] = v
		}
		if _, ok := occurrence[pair[1]]; ok {
			occurrence[pair[1]] += v
		} else {
			occurrence[pair[1]] = v
		}
	}

	// ...except for the first and last characters which were counted once
	// Add the one, and then divide by two at the end
	occurrence[temp.first[0]]++
	occurrence[temp.last[0]]++

	var min, max int
	for _, v := range occurrence {
		if min == 0 || min > v {
			min = v
		}
		if max == 0 || max < v {
			max = v
		}
	}

	// Don't forget to divide by two to fix double counting
	return (max - min) / 2
}

func main() {
	input := readInput()
	template, rules := parseInput(input)

	for i := 0; i < 40; i++ {
		template = polymerStep(template, rules)
	}

	fmt.Println(elementDiff(template))
}
