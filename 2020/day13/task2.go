package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GCD(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(a, b uint64) uint64 {
	return a * b / GCD(a, b)
}

func findEarliest() uint64 {
	scanner := bufio.NewScanner(os.Stdin)

	// Skip a line
	scanner.Scan()
	scanner.Scan()
	constraints := make(map[uint64]uint64)
	for i, idStr := range strings.Split(scanner.Text(), ",") {
		if idStr == "x" {
			continue
		}

		if busID, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			constraints[uint64(i)] = busID
		} else {
			panic("Failed to convert bus ID")
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	var candidate, step uint64 = 0, 1
	for i, busID := range constraints {
		for (candidate+i)%busID != 0 {
			candidate += step
		}

		step = LCM(step, busID)
	}

	return candidate
}

func main() {
	fmt.Println(findEarliest())
}
