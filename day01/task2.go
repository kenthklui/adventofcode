package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readNums() []int {
	nums := make([]int, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		failOnErr(err)

		nums = append(nums, i)
	}
	failOnErr(scanner.Err())

	return nums
}

func findSum(goal, count int, sorted []int) ([]int, error) {
	if count == 1 {
		for i := 0; i < len(sorted); i++ {
			if sorted[i] == goal {
				return []int{sorted[i]}, nil
			} else if sorted[i] > goal {
				return []int{}, fmt.Errorf("Nope")
			}
		}
	}

	for i := 0; i < len(sorted); i++ {
		if sorted[i] >= goal {
			return []int{}, fmt.Errorf("Nope")
		}

		newGoal := goal - sorted[i]
		if solution, err := findSum(newGoal, count-1, sorted[i+1:]); err == nil {
			solution = append(solution, sorted[i])
			return solution, nil
		}
	}

	return []int{}, fmt.Errorf("Nope")
}

func printSolution(solution []int) {
	product := solution[0]
	output := fmt.Sprintf("%d", solution[0])

	for i := 1; i < len(solution); i++ {
		output += fmt.Sprintf(" * %d", solution[i])
		product *= solution[i]
	}
	output += fmt.Sprintf(" = %d\n", product)

	fmt.Print(output)
}

func main() {
	nums := readNums()
	sort.Ints(nums)

	if solution, err := findSum(2020, 3, nums); err == nil {
		printSolution(solution)
	} else {
		fmt.Print("Nope!\n")
	}
}
