package main

import (
	"bufio"
	"fmt"
	"os"
)

func failOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type combo struct {
	FirstIndex, SecondIndex int
}

func readPreambleSize() (int, error) {
	var preambleSize int
	if len(os.Args) < 2 {
		preambleSize = 25
	} else {
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

func checkXmas(preambleSize int) (int, error) {
	values := readValues()

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
			fmt.Printf(
				"Old sum: %d[%d] = %d[%d] + %d[%d}",
				value, i,
				values[c.FirstIndex], c.FirstIndex,
				values[c.SecondIndex], c.SecondIndex,
			)
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

func main() {
	preambleSize, err := readPreambleSize()
	failOnErr(err)

	value, err := checkXmas(preambleSize)
	failOnErr(err)
	fmt.Println(value)
}
