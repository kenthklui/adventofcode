package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

const firstWorkflowName = "in"
const acceptedName = "A"
const rejectedName = "R"

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
	op         byte
	typ, value int
	dest       string
}

type workflow struct {
	name  string
	rules []rule
	dest  string
}

type partsRange struct {
	values [4][2]int
}

func (pr partsRange) volume() int {
	product := 1
	for _, r := range pr.values {
		product *= r[1] - r[0] + 1
	}
	return product
}

func defaultPartsRange() partsRange {
	defRange := [2]int{1, 4000}
	return partsRange{[4][2]int{defRange, defRange, defRange, defRange}}
}

func (r rule) applyRange(pr partsRange) (*partsRange, *partsRange) {
	diverted, remains := pr, pr
	switch r.op {
	case '<':
		if r.value <= pr.values[r.typ][0] {
			return nil, &remains
		} else if r.value <= pr.values[r.typ][1] {
			diverted.values[r.typ][1] = r.value - 1
			remains.values[r.typ][0] = r.value
			return &diverted, &remains
		} else {
			return &diverted, nil
		}
	case '>':
		if r.value >= pr.values[r.typ][1] {
			return nil, &remains
		} else if r.value >= pr.values[r.typ][0] {
			diverted.values[r.typ][0] = r.value + 1
			remains.values[r.typ][1] = r.value
			return &diverted, &remains
		} else {
			return &diverted, nil
		}
	default:
		panic("Invalid op")
	}
}

func (w workflow) ranges(pr partsRange) map[partsRange]string {
	rangeMap := make(map[partsRange]string)
	for _, r := range w.rules {
		diverted, remains := r.applyRange(pr)
		if diverted != nil {
			rangeMap[*diverted] = r.dest
		}
		if remains != nil {
			pr = *remains
		} else {
			break
		}
	}
	rangeMap[pr] = w.dest

	return rangeMap
}

func countAccepted(workflows map[string]workflow) int {
	entry := workflows[firstWorkflowName]
	return entry.recurseCountAccepted(workflows, defaultPartsRange())
}

func (w workflow) recurseCountAccepted(workflows map[string]workflow, pr partsRange) int {
	sum := 0
	for subPr, destName := range w.ranges(pr) {
		switch destName {
		case acceptedName:
			sum += subPr.volume()
		case rejectedName:
			continue
		default:
			newWorkflow := workflows[destName]
			sum += newWorkflow.recurseCountAccepted(workflows, subPr)
		}
	}
	return sum
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
			switch t[0] {
			case 'x':
				r.typ = 0
			case 'm':
				r.typ = 1
			case 'a':
				r.typ = 2
			case 's':
				r.typ = 3
			}
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

	// parts := parseParts(input[empty+1:])
	workflows := parseWorkflows(input[:empty])
	fmt.Println(countAccepted(workflows))
}
