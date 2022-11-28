package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	value               int
	left, right, parent *node
}

func add(s1, s2 *node, printSteps bool) *node {
	n := &node{
		left:  s1.dupe(),
		right: s2.dupe(),
	}
	n.left.parent = n
	n.right.parent = n

	if printSteps {
		fmt.Println("after addition:", n)
	}

	for {
		exploded := n.reduce()
		if exploded == 0 {
			break
		} else if printSteps {
			switch exploded {
			case 1:
				fmt.Println("after explode: ", n)
			case 2:
				fmt.Println("after split:   ", n)
			}
		}
	}

	return n
}

func (n *node) String() string {
	if n.left != nil {
		return fmt.Sprintf("[%s,%s]", n.left.String(), n.right.String())
	} else {
		return fmt.Sprintf("%d", n.value)
	}
}

func (n *node) dupe() *node {
	if n.left != nil {
		left := n.left.dupe()
		right := n.right.dupe()

		dup := &node{
			value: n.value,
			left:  left,
			right: right,
		}
		left.parent = dup
		right.parent = dup

		return dup
	} else {
		return &node{value: n.value}
	}
}

func (n *node) magnitude() int {
	if n.left == nil {
		return n.value
	} else {
		return n.left.magnitude()*3 + n.right.magnitude()*2
	}
}

func (n *node) leftMostNode() *node {
	next := n
	for next.left != nil {
		next = next.left
	}

	return next
}

func (n *node) rightMostNode() *node {
	next := n
	for next.right != nil {
		next = next.right
	}

	return next
}

func (n *node) nextLeft() *node {
	curr := n
	for curr.parent != nil {
		if curr == curr.parent.right {
			return curr.parent.left.rightMostNode()
		} else {
			curr = curr.parent
		}
	}

	return nil
}

func (n *node) nextRight() *node {
	curr := n
	for curr.parent != nil {
		if curr == curr.parent.left {
			return curr.parent.right.leftMostNode()
		} else {
			curr = curr.parent
		}
	}

	return nil
}

func (n *node) reduce() int {
	if n.recursiveExplode(0) {
		return 1
	}
	if n.recursiveSplit() {
		return 2
	}

	return 0
}

func (n *node) recursiveExplode(level int) bool {
	if n.left != nil { // Not leaf node
		return n.left.recursiveExplode(level+1) || n.right.recursiveExplode(level+1)
	}

	if level > 4 {
		n.parent.explode()
		return true
	}

	return false
}

func (n *node) recursiveSplit() bool {
	if n.left != nil { // Not leaf node
		return n.left.recursiveSplit() || n.right.recursiveSplit()
	}

	if n.value >= 10 {
		n.split()
		return true
	}

	return false
}

func (n *node) explode() {
	if left := n.nextLeft(); left != nil {
		left.value += n.left.value
	}
	if right := n.nextRight(); right != nil {
		right.value += n.right.value
	}

	n.left = nil
	n.right = nil
	n.value = 0
}

func (n *node) split() {
	leftValue := n.value / 2
	rightValue := n.value - leftValue

	n.value = 0
	n.left = &node{value: leftValue, parent: n}
	n.right = &node{value: rightValue, parent: n}
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

func parseSnailfishString(s string) *node {
	if s[0] == '[' {
		stackLevel := -1
		n := &node{}
		for i, r := range s {
			switch r {
			case '[':
				stackLevel++
			case ']':
				stackLevel--
			case ',':
				if stackLevel == 0 {
					n.left = parseSnailfishString(s[1:i])
					n.left.parent = n
					n.right = parseSnailfishString(s[i+1 : len(s)-1])
					n.right.parent = n
					break
				}
			default:
				continue
			}
		}

		return n
	} else {
		var value int
		if _, err := fmt.Sscanf(s, "%d", &value); err == nil {
			return &node{value: value}
		} else {
			panic(err)
		}
	}
}

func parseInput(input []string) []*node {
	snailfish := make([]*node, len(input))

	for i, line := range input {
		snailfish[i] = parseSnailfishString(line)
	}

	return snailfish
}

func main() {
	input := readInput()
	snailfish := parseInput(input)

	var mag, largestMagnitude int
	for i, s1 := range snailfish[:len(snailfish)-1] {
		for _, s2 := range snailfish[i+1:] {
			mag = add(s1, s2, false).magnitude()
			if mag > largestMagnitude {
				largestMagnitude = mag
			}

			mag = add(s2, s1, false).magnitude()
			if mag > largestMagnitude {
				largestMagnitude = mag
			}
		}
	}

	fmt.Println(largestMagnitude)
}
