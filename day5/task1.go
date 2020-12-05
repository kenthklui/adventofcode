package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func maxSeat() int64 {
	var max int64

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		text = strings.ReplaceAll(text, "F", "0")
		text = strings.ReplaceAll(text, "B", "1")
		text = strings.ReplaceAll(text, "L", "0")
		text = strings.ReplaceAll(text, "R", "1")

		if i, err := strconv.ParseInt(text, 2, 64); err == nil {
			if i > max {
				max = i
			}
		} else {
			log.Fatalf("Parse error: %s\n", scanner.Text())
		}
	}

	return max
}

func main() {
	fmt.Println(maxSeat())
}
