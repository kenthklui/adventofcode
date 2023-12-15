package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

func hash(s string) int {
	value := 0
	for _, r := range s {
		value += int(r)
		value *= 17
		value %= 256
	}
	return value
}

func main() {
	input := util.StdinReadlines()
	for _, line := range input {
		sum := 0
		for _, s := range strings.Split(line, ",") {
			sum += hash(s)
		}
		fmt.Println(sum)
	}
}
