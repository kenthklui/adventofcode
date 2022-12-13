package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type data interface {
	typ() int
}

type list struct {
	val []data
}

type packets []*list

func (ps packets) Len() int           { return len(ps) }
func (ps packets) Less(i, j int) bool { return inOrder(ps[i], ps[j]) == 1 }
func (ps packets) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

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

func parseInput(input []string) packets {
	ps := make(packets, 0)
	input = append(input, "[[2]]", "[[6]]")

	for _, line := range input {
		if line == "" {
			continue
		}
		ps = append(ps, parsePacket(line))
	}

	return ps
}

func main() {
	input := readInput()
	ps := parseInput(input)

	divider1, divider2 := ps[len(ps)-2], ps[len(ps)-1]
	sort.Sort(ps)

	var index1, index2 int
	for i, p := range ps {
		if p == divider1 {
			index1 = i + 1
		}
		if p == divider2 {
			index2 = i + 1
		}
	}
	fmt.Println(index1 * index2)
}
