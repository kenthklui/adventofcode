package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

func solve(input []string) (output string) {
	ints := util.ParseInts(input)
	list1, list2 := make([]int, len(ints)), make(map[int]int)
	for i := range ints {
		list1[i] = ints[i][0]
		if _, ok := list2[ints[i][1]]; ok {
			list2[ints[i][1]]++
		} else {
			list2[ints[i][1]] = 1
		}
	}
	distance := 0
	for _, k := range list1 {
		if v, ok := list2[k]; ok {
			distance += (k * v)
		}
	}
	return strconv.Itoa(distance)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
