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

func findDistribution(adapters []int) []int {

	distribution := make([]int, 4)

	for i, adapter := range adapters[1:] {
		diff := adapter - adapters[i]
		distribution[diff]++
	}

	return distribution
}

func main() {
	adapters := readAdapters()
	distribution := findDistribution(adapters)

	fmt.Println(distribution[1] * distribution[3])
}
