package main

import (
	"fmt"

	"github.com/kenthklui/adventofcode/util"
)

func binomial(row, col int) int {
	n, r := row, col
	if r+r > n {
		r = n - r
	}
	product := 1
	for i := 0; i < r; i++ {
		product *= (n - i)
	}
	for i := r; i > 0; i-- {
		product /= i
	}
	if col%2 == 0 {
		return product
	} else {
		return -product
	}
}

func history(intLine []int) int {
	order := 1
	for ; order < len(intLine); order++ {
		allZero := true
		for i := order; i <= order+1; i++ {
			sum := 0
			for j := 0; j <= order; j++ {
				sum += intLine[i-j] * binomial(order, j)
			}
			if sum != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			break
		}
	}

	nextValue := 0
	for i := 1; i <= order; i++ {
		nextValue -= intLine[len(intLine)-i] * binomial(order, i)
	}

	return nextValue
}

func main() {
	input := util.StdinReadlines()
	ints := util.ParseInts(input)

	sum := 0
	for _, intLine := range ints {
		sum += history(intLine)
	}
	fmt.Println(sum)
}
