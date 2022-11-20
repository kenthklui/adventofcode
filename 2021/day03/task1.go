package main

import (
	"bufio"
	"fmt"
	"os"
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

func getGammaEpsilon(input []string) (int, int) {
	digits := len(input[0])
	onesOverZeros := make([]int, digits)

	for _, s := range input {
		for i, c := range s {
			switch c {
			case '0':
				onesOverZeros[i]--
			case '1':
				onesOverZeros[i]++
			default:
				panic("Non binary digit")
			}
		}
	}

	var gamma, epsilon int
	for _, n := range onesOverZeros {
		gamma *= 2
		epsilon *= 2

		if n > 0 {
			gamma++
		} else {
			epsilon++
		}
	}

	return gamma, epsilon
}

func main() {
	input := readInput()
	gamma, epsilon := getGammaEpsilon(input)

	fmt.Println(gamma * epsilon)
}
