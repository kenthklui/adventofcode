package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
)

func mix(nums []int, key, cycles int) []int {
	l := list.New()
	elements := make([]*list.Element, 0, len(nums))

	for _, n := range nums {
		newVal := n * key
		e := l.PushBack(newVal)
		elements = append(elements, e)
	}

	for cycle := 0; cycle < cycles; cycle++ {
		for _, e := range elements {
			val := e.Value.(int) % (len(nums) - 1)
			curr := e

			if val > 0 {
				var next *list.Element
				for i := val; i > 0; i-- {
					if curr.Next() != nil {
						next = curr.Next()
					} else {
						next = l.Front()
					}
					curr = next
				}
				if curr == l.Back() {
					l.MoveBefore(e, l.Front())
				} else {
					l.MoveAfter(e, curr)
				}
			} else if val < 0 {
				var prev *list.Element
				for i := val; i < 0; i++ {
					if curr.Prev() != nil {
						prev = curr.Prev()
					} else {
						prev = l.Back()
					}
					curr = prev
				}
				if curr == l.Front() {
					l.MoveAfter(e, l.Back())
				} else {
					l.MoveBefore(e, curr)
				}
			}
		}
	}

	newNums := make([]int, 0, len(nums))
	for e := l.Front(); e != nil; e = e.Next() {
		newNums = append(newNums, e.Value.(int))
	}

	return newNums
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

func parseInput(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, line := range input {
		if i, err := strconv.Atoi(line); err == nil {
			nums = append(nums, i)
		} else {
			panic(err)
		}
	}
	return nums
}

func zeroIndex(nums []int) int {
	zeroIndex := -1
	for i, n := range nums {
		if n == 0 {
			zeroIndex = i
			break
		}
	}
	if zeroIndex == -1 {
		panic("No zero!")
	}
	return zeroIndex
}

const decryptionKey = 811589153
const mixCycles = 10

func main() {
	input := readInput()
	nums := parseInput(input)
	mixed := mix(nums, decryptionKey, mixCycles)

	zIndex := zeroIndex(mixed)
	sum := 0
	for i := 1000; i <= 3000; i += 1000 {
		index := (zIndex + i) % len(mixed)
		sum += mixed[index]
	}
	fmt.Println(sum)
}
