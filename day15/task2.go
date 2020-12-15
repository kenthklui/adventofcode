package main

import (
	"bufio"
	"fmt"
	"os"
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

func seed(input []string) ([]int, map[int]int) {
	numStr := strings.Split(input[0], ",")
	nums := make([]int, len(numStr))
	lastOccurence := make(map[int]int)

	for i, s := range numStr {
		if num, err := strconv.Atoi(s); err == nil {
			nums[i] = num
			lastOccurence[num] = i
		} else {
			panic(err)
		}
	}

	return nums, lastOccurence
}

func iterate(nums []int, lastOccurence map[int]int, target int) []int {
	for len(nums) < target {
		prevIndex := len(nums) - 1
		prev := nums[prevIndex]

		lastOccur, ok := lastOccurence[prev]

		var newNum int
		if ok {
			newNum = prevIndex - lastOccur
		}

		lastOccurence[prev] = prevIndex
		nums = append(nums, newNum)
	}

	return nums
}

func main() {
	input := readInput()
	nums, lastOccurence := seed(input)
	nums = iterate(nums, lastOccurence, 30000000)
	fmt.Println(nums[len(nums)-1])
}
