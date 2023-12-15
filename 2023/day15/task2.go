package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

func hash(s string) int {
	value := 0
	for _, r := range s {
		value += int(r)
		value *= 17
		value %= 256
	}
	return value
}

type lens struct {
	label string
	focal int
	next  *lens
}

type action struct {
	op     byte
	boxNum int
	lens   *lens
}

func parseAction(s string) action {
	opIndex := strings.IndexAny(s, "-=")

	label := s[:opIndex]

	focal := 0
	op := s[opIndex]
	if op == '=' {
		focal = int(s[opIndex+1] - '0')
	}

	l := lens{label, focal, nil}

	return action{op, hash(label), &l}
}

type boxes struct {
	b [256]*lens
}

func (b *boxes) handle(a action) {
	l := b.b[a.boxNum]
	switch a.op {
	case '-':
		// Skip if box is empty
		if l == nil {
			return
		}

		// Remove and reset head value if at front
		if l.label == a.lens.label {
			b.b[a.boxNum] = l.next
			return
		}

		// Find and remove node if in middle onward
		for n := l.next; n != nil; l, n = l.next, n.next {
			if n.label == a.lens.label {
				l.next = n.next
				return
			}
		}
	case '=':
		// Insert at front if empty
		if l == nil {
			b.b[a.boxNum] = a.lens
			return
		}

		// Replace if present
		var prev *lens
		for l != nil {
			if l.label == a.lens.label {
				l.focal = a.lens.focal
				return
			}
			prev = l
			l = l.next
		}

		// Add to end if not present
		prev.next = a.lens
	}
}

func (b *boxes) focusingPower() int {
	sum := 0
	for i, l := range b.b {
		for lensNum := 1; l != nil; l = l.next {
			sum += (i + 1) * lensNum * l.focal
			lensNum++
		}
	}
	return sum
}

func main() {
	input := util.StdinReadlines()
	b := &boxes{}
	for _, line := range input {
		for _, s := range strings.Split(line, ",") {
			action := parseAction(s)
			b.handle(action)
		}
	}
	fmt.Println(b.focusingPower())
}
