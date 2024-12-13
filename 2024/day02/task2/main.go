package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

type report []int

func (r report) safe() bool {
	return (r.increasing() || r.decreasing()) && r.adjacent()
}

func (r report) tolerantSafe() bool {
	if r.safe() {
		return true
	}

	subreport := make(report, len(r)-1)
	for i := range r {
		copy(subreport, r[:i])
		copy(subreport[i:], r[i+1:])
		if subreport.safe() {
			return true
		}
	}
	return false
}

func (r report) increasing() bool {
	for i, v := range r[1:] {
		if v < r[i] {
			return false
		}
	}
	return true
}

func (r report) decreasing() bool {
	for i, v := range r[1:] {
		if v > r[i] {
			return false
		}
	}
	return true
}

func (r report) adjacent() bool {
	for i, v := range r[1:] {
		if diff := v - r[i]; diff == 0 {
			return false
		} else if diff*diff > 9 {
			return false
		}
	}
	return true
}

func parseInput(input []string) []report {
	ints := util.ParseInts(input)
	reports := make([]report, 0, len(input))
	for _, i := range ints {
		reports = append(reports, i)
	}
	return reports
}

func solve(input []string) (output string) {
	reports := parseInput(input)
	safe := 0
	for _, r := range reports {
		if r.tolerantSafe() {
			safe++
		}
	}
	return fmt.Sprintf("%d", safe)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
