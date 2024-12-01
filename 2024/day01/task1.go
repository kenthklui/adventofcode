package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

func solve(input []string) (output string) {
	ints := util.ParseInts(input)
	list1, list2 := make([]int, len(ints)), make([]int, len(ints))
	for i := range ints {
		list1[i] = ints[i][0]
		list2[i] = ints[i][1]
	}
	sort.Ints((list1))
	sort.Ints((list2))
	distance := 0
	for i := range list1 {
		if list1[i] > list2[i] {
			distance += list1[i] - list2[i]
		} else {
			distance += list2[i] - list1[i]
		}
	}
	return strconv.Itoa(distance)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
