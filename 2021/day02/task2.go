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

func getPosition(input []string) (int, int) {
	var horizontal, aim, depth int

	for _, s := range input {
		command := strings.Split(s, " ")

		num, err := strconv.Atoi(command[1])
		if err != nil {
			panic(err)
		}

		switch command[0] {
		case "forward":
			horizontal += num
			depth += num * aim
		case "down":
			aim += num
		case "up":
			aim -= num
		default:
			panic("Unknown command direction")
		}
	}

	return horizontal, depth
}

func main() {
	input := readInput()

	h, d := getPosition(input)
	fmt.Println(h * d)
}
