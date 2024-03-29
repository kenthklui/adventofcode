package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const defaultPreambleSize = 25

func failOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type combo struct {
	FirstIndex, SecondIndex int
}

func readPreambleSize() (int, error) {
	preambleSize := defaultPreambleSize
	if len(os.Args) >= 2 {
		_, err := fmt.Sscanf(os.Args[1], "%d", &preambleSize)
		if err != nil {
			err := fmt.Errorf("Failed to parse preamble size: %q, reason: %q",
				os.Args[1], err.Error())
			return 0, err
		}
	}
	return preambleSize, nil
}

func readValues() []int {
	values := make([]int, 0)

	var value int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n, err := fmt.Sscanf(scanner.Text(), "%d", &value)
		if err != nil || n != 1 {
			panic(fmt.Errorf("Failed to parse line: %q", scanner.Text()))
		}

		values = append(values, value)
	}
	failOnErr(scanner.Err())

	return values
}

func findNonSum(values []int, preambleSize int) (int, error) {
	sums := make(map[int]combo)
	for i := 0; i < preambleSize; i++ {
		for j := 1; j < preambleSize; j++ {
			sums[values[i]+values[j]] = combo{i, j}
		}
	}

	for i := preambleSize; i < len(values); i++ {
		value := values[i]
		if c, ok := sums[value]; !ok { // not a sum
			return value, nil
		} else if i-c.FirstIndex > preambleSize { // outdated sum
			/*
				fmt.Printf(
					"Old sum: %d[%d] = %d[%d] + %d[%d]\n",
					value, i,
					values[c.FirstIndex], c.FirstIndex,
					values[c.SecondIndex], c.SecondIndex,
				)
			*/
			return value, nil
		}

		for j := i - preambleSize + 1; j < i; j++ {
			prevValue := values[j]
			sum := value + prevValue
			if c, ok := sums[sum]; !ok || c.FirstIndex < j {
				sums[sum] = combo{j, i}
			}
		}
	}

	return 0, fmt.Errorf("Not found")
}

func findContSum(values []int, target int) (combo, error) {
	runs := make([]int, len(values))

	for i, iValue := range values {
		runs[i] = iValue

		for j, run := range runs[:i] {
			newRun := run + iValue
			if newRun == target {

				// fmt.Printf("Target %d is a sum from [%d] to [%d]\n", target, j, i)
				return combo{j, i}, nil
			}

			runs[j] = newRun
		}
	}

	return combo{}, fmt.Errorf("Not found")
}

func computeWeakness(values []int, c combo) int {
	min := slices.Min(values[c.FirstIndex : c.SecondIndex+1])
	max := slices.Max(values[c.FirstIndex : c.SecondIndex+1])

	return min + max
}

func main() {
	values := readValues()

	preambleSize, err := readPreambleSize()
	failOnErr(err)

	nonSum, err := findNonSum(values, preambleSize)
	failOnErr(err)

	c, err := findContSum(values, nonSum)
	failOnErr(err)

	weakness := computeWeakness(values, c)
	fmt.Println(weakness)
}
