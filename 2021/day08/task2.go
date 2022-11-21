package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
For our clarity, let's define the "correct" 7 segment display like this:
 aaa
b   c
b   c
 ddd
e   f
e   f
 ggg
*/

type segments []rune

func (s segments) Len() int {
	return len(s)
}

func (s segments) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s segments) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s segments) ReadAsNum() (int, error) {
	sort.Sort(s)
	str := string(s)

	var n int
	var err error

	switch str {
	case "abcefg":
		n = 0
	case "cf":
		n = 1
	case "acdeg":
		n = 2
	case "acdfg":
		n = 3
	case "bcdf":
		n = 4
	case "abdfg":
		n = 5
	case "abdefg":
		n = 6
	case "acf":
		n = 7
	case "abcdefg":
		n = 8
	case "abcdfg":
		n = 9
	default:
		err = fmt.Errorf("Invalid segments: %s", str)
	}

	return n, err
}

type Limitations map[rune]map[rune]rune

type puzzle struct {
	patterns []string
	output   []string

	limits Limitations
}

func NewPuzzle(patterns, output []string) *puzzle {
	p := puzzle{patterns: patterns, output: output}
	p.SetLimits()

	return &p
}

func (p *puzzle) SetLimits() {
	p.limits = make(Limitations)

	for r1 := 'a'; r1 <= 'g'; r1++ {
		l := make(map[rune]rune)
		for r2 := 'a'; r2 <= 'g'; r2++ {
			l[r2] = r2
		}
		p.limits[r1] = l
	}

	for _, pat := range p.patterns {
		disallowedSegments := []rune{}

		switch len(pat) {
		case 2: // 1
			disallowedSegments = []rune{'a', 'b', 'd', 'e', 'g'}
		case 3: // 7
			disallowedSegments = []rune{'b', 'd', 'e', 'g'}
		case 4: // 4
			disallowedSegments = []rune{'a', 'e', 'g'}
		case 5: // 2, 3, 5
			// nothing
		case 6: // 0, 6, 9
			// nothing
		case 7: // 8
			// nothing
		default:
			panic("Weird number of segments")
		}

		for _, disallow := range disallowedSegments {
			for _, patSeg := range pat {
				if _, ok := p.limits[patSeg][disallow]; ok {
					delete(p.limits[patSeg], disallow)
				}
			}
		}
	}
}

func (p *puzzle) SanityCheck(solution map[rune]rune) bool {
	for _, pat := range p.patterns {
		s := segments(pat)
		for i, r := range s {
			s[i] = solution[r]
		}

		if _, err := s.ReadAsNum(); err != nil {
			return false
		}
	}
	return true
}

func (p *puzzle) OutputNum() int {
	solution := p.Solve(map[rune]rune{}, map[rune]rune{}, 'a')

	outputNum := 0
	for _, o := range p.output {
		outputNum *= 10

		s := segments(o)
		for i, r := range s {
			s[i] = solution[r]
		}

		if n, err := s.ReadAsNum(); err == nil {
			outputNum += n
		}
	}

	return outputNum
}

func (p *puzzle) Solve(solved, taken map[rune]rune, next rune) map[rune]rune {
	if next > 'g' {
		if p.SanityCheck(solved) {
			return solved
		} else {
			return nil
		}
	}

	for candidate := range p.limits[next] {
		if _, ok := taken[candidate]; ok {
			continue
		}

		solved[next] = candidate
		taken[candidate] = next

		if solution := p.Solve(solved, taken, next+1); solution != nil {
			return solution
		} else {
			delete(solved, next)
			delete(taken, candidate)
		}
	}

	return nil
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

func parseInput(input []string) []*puzzle {
	puzzles := make([]*puzzle, len(input))
	for i, s := range input {
		splits := strings.Split(s, " | ")

		patterns := strings.Split(splits[0], " ")
		output := strings.Split(splits[1], " ")

		puzzles[i] = NewPuzzle(patterns, output)
	}

	return puzzles
}

func main() {
	input := readInput()
	puzzles := parseInput(input)

	sum := 0
	for _, p := range puzzles {
		sum += p.OutputNum()
	}

	fmt.Println(sum)
}
