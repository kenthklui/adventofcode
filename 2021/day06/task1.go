package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput() []string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

func parseInput(input []string) []int {
	fish := make([]int, 9)

	for _, s := range strings.Split(input[0], ",") {
		if n, err := strconv.Atoi(s); err == nil {
			fish[n]++
		}
	}

	return fish
}

func iterateFish(fish []int) []int {
	newFish := make([]int, 9)

	for i, n := range fish[1:] {
		newFish[i] = n
	}
	newFish[8] += fish[0]
	newFish[6] += fish[0]

	return newFish
}

func sumFish(fish []int) int {
	sum := 0
	for _, n := range fish {
		sum += n
	}

	return sum
}

func main() {
	input := readInput()
	fish := parseInput(input)

	for i := 0; i < 18; i++ {
		fish = iterateFish(fish)
	}

	fmt.Println("18 days:", sumFish(fish))

	for i := 0; i < 62; i++ {
		fish = iterateFish(fish)
	}

	fmt.Println("80 days:", sumFish(fish))
}
