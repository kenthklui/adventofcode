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

func seed(input []string, target int) ([]int, []int) {
	numStr := strings.Split(input[0], ",")
	nums := make([]int, len(numStr))

	lastOccurence := make([]int, target)
	for i := range lastOccurence {
		lastOccurence[i] = -1
	}

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

func iterate(nums []int, lastOccurence []int, target int) int {
	prev := nums[len(nums)-1]
	for i := len(nums); i < target; i++ {
		var newNum int

		if lastOccur := lastOccurence[prev]; lastOccur != -1 {
			newNum = i - 1 - lastOccur
		}

		lastOccurence[prev] = i - 1
		prev = newNum
	}

	return prev
}

func main() {
	input := readInput()
	target := 30000000
	nums, lastOccurence := seed(input, target)
	num := iterate(nums, lastOccurence, target)
	fmt.Println(num)
}
