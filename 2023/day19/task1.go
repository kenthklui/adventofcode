package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type part struct {
	x, m, a, s int
}

func parseParts(partStrs []string) []part {
	parts := make([]part, 0, len(partStrs))
	for _, line := range partStrs {
		var p part
		fmt.Sscanf(line, "{x=%d,m=%d,a=%d,s=%d]", &p.x, &p.m, &p.a, &p.s)
		parts = append(parts, p)
	}
	return parts
}

func (p part) ratingSum() int {
	return p.x + p.m + p.a + p.s
}

type rule struct {
	op, typ byte
	value   int
	dest    string
}

type workflow struct {
	name  string
	rules []rule
	dest  string
}

func (w workflow) process(p part) string {
	for _, r := range w.rules {
		var partValue int
		switch r.typ {
		case 'x':
			partValue = p.x
		case 'm':
			partValue = p.m
		case 'a':
			partValue = p.a
		case 's':
			partValue = p.s
		}
		switch r.op {
		case '<':
			if partValue < r.value {
				return r.dest
			}
		case '>':
			if partValue > r.value {
				return r.dest
			}
		}
	}
	return w.dest
}

func parseWorkflows(workflowStrs []string) map[string]workflow {
	workflows := make(map[string]workflow)
	for _, line := range workflowStrs {
		name, after, _ := strings.Cut(line, "{")

		tokens := strings.Split(strings.Trim(after, "}"), ",")

		last := len(tokens) - 1
		dest := tokens[last]
		tokens = tokens[:last]

		rules := make([]rule, 0, len(tokens))
		for _, t := range tokens {
			var r rule
			r.typ = t[0]
			r.op = t[1]
			colonIndex := strings.Index(t, ":")
			r.dest = t[colonIndex+1:]
			if value, err := strconv.Atoi(t[2:colonIndex]); err == nil {
				r.value = value
			}
			rules = append(rules, r)
		}

		w := workflow{name: name, rules: rules, dest: dest}
		workflows[name] = w
	}
	return workflows
}

func main() {
	input := util.StdinReadlines()
	empty := slices.Index(input, "")

	parts := parseParts(input[empty+1:])

	workflows := parseWorkflows(input[:empty])

	sum := 0
	for _, p := range parts {
		workflowName := "in"
		for workflowName != "A" && workflowName != "R" {
			workflowName = workflows[workflowName].process(p)
		}
		if workflowName == "A" {
			sum += p.ratingSum()
		}
	}
	fmt.Println(sum)
}
