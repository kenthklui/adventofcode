package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readBitmask(bitmask string) (int, int) {
	var ones, zeroes int

	for i, r := range bitmask {
		if r == 'X' {
			continue
		}

		significance := 1 << (35 - i)
		switch r {
		case '1':
			ones += significance
		case '0':
			zeroes += significance
		}
	}

	// fmt.Println(ones, zeroes)
	return ones, zeroes
}

func readInput() int {
	scanner := bufio.NewScanner(os.Stdin)

	var ones, zeroes int
	values := make(map[int]int)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mask = ") {
			bitmask := strings.ReplaceAll(line, "mask = ", "")
			ones, zeroes = readBitmask(bitmask)
		} else {
			var index, value int
			n, err := fmt.Sscanf(scanner.Text(), "mem[%d] = %d", &index, &value)
			if err != nil {
				log.Fatalf("Cannot parse line, error: %s, line: %q", err.Error(), scanner.Text())
			} else if n != 2 {
				panic("Sscanf error")
			}

			maskedValue := (value | ones) &^ zeroes
			values[index] = maskedValue
		}
	}

	sum := 0
	for _, value := range values {
		sum += value
	}

	return sum
}

func main() {
	fmt.Println(readInput())
}
