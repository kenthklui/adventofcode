package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type beforeMap map[int][]int

func parsebeforeMap(input []string) beforeMap {
	var a, b int
	bm := make(beforeMap)
	for _, line := range input {
		fmt.Sscanf(line, "%d|%d", &a, &b)
		bm[a] = append(bm[a], b)
	}
	return bm
}

func parsePages(input []string) [][]int {
	pages := make([][]int, len(input))
	for i, line := range input {
		tokens := strings.Split(line, ",")
		pages[i] = make([]int, len(tokens))
		for j, token := range tokens {
			pages[i][j], _ = strconv.Atoi(token)
		}
	}
	return pages
}

func validatePage(page []int, bm beforeMap) bool {
	for i := len(page) - 1; i > 0; i-- {
		for _, front := range page[:i] {
			for _, banned := range bm[page[i]] {
				if banned == front {
					return false
				}
			}
		}
	}
	return true
}

func pageMiddle(page []int) int {
	return page[len(page)/2]
}

func solve(input []string) (output string) {
	index := -1
	for i, line := range input {
		if line == "" {
			index = i
		}
	}
	if index == -1 {
		panic("Invalid input")
	}

	beforeMap := parsebeforeMap(input[:index])
	pages := parsePages(input[index+1:])

	sum := 0
	for _, page := range pages {
		if validatePage(page, beforeMap) {
			sum += pageMiddle(page)
		}
	}
	return strconv.Itoa(sum)
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
