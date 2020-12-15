package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readNums() []int {
	nums := make([]int, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		failOnErr(err)

		nums = append(nums, i)
	}
	failOnErr(scanner.Err())

	return nums
}

func find2020Sum(nums []int) {
	sort.Ints(nums)

	for i, j := 0, len(nums)-1; i < j; {
		sum := nums[i] + nums[j]

		if sum == 2020 {
			fmt.Printf("%d * %d = %d\n", nums[i], nums[j], nums[i]*nums[j])
			return
		} else if sum < 2020 {
			i++
		} else {
			j--
		}
	}
}

func main() {
	nums := readNums()

	find2020Sum(nums)
}
