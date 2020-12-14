package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readBitmask(bitmask string) (int64, []int64) {
	var ones int64
	floats := make([]int64, 0)

	for i, r := range bitmask {
		significance := 1 << (35 - i)
		switch r {
		case '1':
			ones += int64(significance)
		case '0':
			continue
		case 'X':
			floats = append(floats, int64(35-i))
		}
	}

	// fmt.Println(ones, floats)
	return ones, floats
}

func readInput() int64 {
	scanner := bufio.NewScanner(os.Stdin)

	var ones int64
	var floats []int64
	values := make(map[int64]int64)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mask = ") {
			bitmask := strings.ReplaceAll(line, "mask = ", "")
			ones, floats = readBitmask(bitmask)
		} else {
			var index, value int64
			n, err := fmt.Sscanf(scanner.Text(), "mem[%d] = %d", &index, &value)
			if err != nil {
				log.Fatalf("Cannot parse line, error: %s, line: %q", err.Error(), scanner.Text())
			} else if n != 2 {
				panic("Sscanf error")
			}

			overwrittenIndex := index | ones

			for i := 0; i < (1 << len(floats)); i++ {
				var floatOnes, floatZeroes int64
				for j, f := range floats {
					if (1<<j)&i == 0 {
						floatZeroes += (1 << f)
					} else {
						floatOnes += (1 << f)
					}
				}

				maskedIndex := (overwrittenIndex | floatOnes) &^ floatZeroes
				values[maskedIndex] = value
			}
		}
	}

	var sum int64
	for _, value := range values {
		sum += value
	}

	return sum
}

func main() {
	fmt.Println(readInput())
}
