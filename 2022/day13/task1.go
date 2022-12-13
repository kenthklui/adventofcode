package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type data interface {
	typ() int
}

type list struct {
	val []data
}

type pair [2]*list

func (l *list) typ() int { return 0 }

type integer struct {
	val int
}

func (i *integer) typ() int       { return 1 }
func (i *integer) String() string { return strconv.Itoa(i.val) }

func inOrder(a, b data) int {
	// fmt.Println("Comparing", a, b)
	switch a.typ() + b.typ() {
	case 2:
		aVal, bVal := a.(*integer).val, b.(*integer).val
		if aVal < bVal {
			return 1
		} else if aVal == bVal {
			return 0
		} else {
			return -1
		}
	case 0:
		aList, bList := a.(*list), b.(*list)
		for i := range aList.val {
			if i >= len(bList.val) {
				return -1
			}
			switch inOrder(aList.val[i], bList.val[i]) {
			case -1:
				return -1
			case 0:
				continue
			case 1:
				return 1
			}
		}
		if len(aList.val) < len(bList.val) {
			return 1
		}

		return 0
	case 1:
		if a.typ() == 0 {
			return inOrder(a, &list{val: []data{b}})
		} else {
			return inOrder(&list{val: []data{a}}, b)
		}
	default:
		panic("Invalid typ sum")
	}
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

func parsePacket(line string) *list {
	p := new(list)
	stack := []*list{p}
	curr := stack[len(stack)-1]
	index := 1
	for index < len(line) {
		// fmt.Println(line, index)
		switch line[index] {
		case '[':
			l := new(list)
			l.val = make([]data, 0)
			curr.val = append(curr.val, l)
			stack = append(stack, l)
			curr = l
			index++
		case ']':
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				curr = stack[len(stack)-1]
			}
			index++
		case ',':
			index++
		default:
			endIndex := strings.IndexAny(line[index:], ",]")
			n, err := strconv.Atoi(line[index : index+endIndex])
			if err != nil {
				panic(err)
			}
			curr.val = append(curr.val, &integer{n})
			index += endIndex
		}
	}

	return p
}

func parseInput(input []string) []pair {
	pairs := make([]pair, 0)
	for i := 0; i < len(input); i += 3 {
		var p pair
		for j := 0; j < 2; j++ {
			line := input[i+j]
			p[j] = parsePacket(line)
		}
		pairs = append(pairs, p)
	}

	return pairs
}

func main() {
	input := readInput()
	pairs := parseInput(input)

	ordered := 0
	for i, p := range pairs {
		if inOrder(p[0], p[1]) >= 0 {
			// fmt.Println("Packet", i+1, "in order:", p, "\n")
			ordered += (i + 1)
		} else {
			// fmt.Println("Packet", i+1, "not in order:", p, "\n")
		}
	}

	fmt.Println(ordered)
}
