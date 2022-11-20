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

type rule struct {
	id    int
	combo [][]int
	char  string
	reStr string
	re    *regexp.Regexp
}

func (r rule) String() string {
	return r.reStr
}

func (r *rule) BuildRegex(rules map[int]*rule) string {
	if r.reStr != "" {
		return r.reStr
	}

	if len(r.combo) == 0 {
		r.reStr = "(" + r.char + ")"
	} else {
		subReStrs := make([]string, len(r.combo))
		for i, c := range r.combo {
			subReStr := ""
			for _, rn := range c {
				if rule, ok := rules[rn]; ok {
					subReStr += rule.BuildRegex(rules)
				} else {
					panic("Rule not found")
				}
			}

			subReStrs[i] = subReStr
		}
		r.reStr = fmt.Sprintf("(%s)", strings.Join(subReStrs, "|"))
	}
	r.re = regexp.MustCompile("^" + r.reStr + "$")

	return r.reStr
}

func (r *rule) Match(s string) bool {
	return r.re.MatchString(s)
}

func parseRules(input []string) (map[int]*rule, int) {
	var i int
	var line string

	rules := make(map[int]*rule)

	for i, line = range input {
		if line == "" {
			break
		}

		tokens := strings.Split(line, " ")
		ruleNumber, _ := strconv.Atoi(tokens[0][:len(tokens[0])-1])

		if tokens[1] == "\"a\"" || tokens[1] == "\"b\"" {
			rules[ruleNumber] = &rule{id: ruleNumber, char: tokens[1][1:2]}
			continue
		}

		subrules := make([][]int, 1)
		subrules[0] = make([]int, 0)
		for _, r := range tokens[1:] {
			if r == "|" {
				newSubrule := make([]int, 0)
				subrules = append(subrules, newSubrule)
			} else {
				subRuleNumber, _ := strconv.Atoi(r)
				subrules[len(subrules)-1] = append(subrules[len(subrules)-1], subRuleNumber)
			}
		}

		rules[ruleNumber] = &rule{id: ruleNumber, combo: subrules}
	}

	return rules, i
}

func main() {
	input := readInput()
	rules, offset := parseRules(input)
	rules[0].BuildRegex(rules)

	matched := 0
	for _, m := range input[offset+1:] {
		if rules[0].Match(m) {
			matched++
		}
	}

	fmt.Println(matched)
}
