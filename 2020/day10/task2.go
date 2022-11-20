package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readAdapters() []int {
	// 1 for wall socket
	adapters := make([]int, 1)

	var adapter int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n, err := fmt.Sscanf(scanner.Text(), "%d", &adapter)
		if err != nil || n != 1 {
			log.Fatal("Failed to parse line: %q", scanner.Text())
		}

		adapters = append(adapters, adapter)
	}
	failOnErr(scanner.Err())

	sort.Ints(adapters)

	// your device
	adapters = append(adapters, adapters[len(adapters)-1]+3)

	return adapters
}

func countCombinations(adapters []int) int {
	combsWithN := make([]int, len(adapters))
	combsWithN[0] = 1
	combsWithN[1] = 1
	if adapters[2]-adapters[0] <= 3 {
		combsWithN[2] = 2
	} else {
		combsWithN[2] = 1
	}

	for i := 3; i < len(adapters); i++ {
		for j := 1; j <= 3; j++ {
			if adapters[i]-adapters[i-j] <= 3 {
				combsWithN[i] += combsWithN[i-j]
			}
		}
	}

	return combsWithN[len(adapters)-1]
}

func main() {
	adapters := readAdapters()

	fmt.Println(countCombinations(adapters))
}
