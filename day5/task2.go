package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func findMissingSeat(min int64) int64 {
	occupied := make(map[int64]bool)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		text = strings.ReplaceAll(text, "F", "0")
		text = strings.ReplaceAll(text, "B", "1")
		text = strings.ReplaceAll(text, "L", "0")
		text = strings.ReplaceAll(text, "R", "1")

		if i, err := strconv.ParseInt(text, 2, 64); err == nil {
			occupied[i] = true
		} else {
			log.Fatalf("Parse error: %s", scanner.Text())
		}
	}

	for i := min; i < 1e6; i++ {
		if _, ok := occupied[i]; !ok {
			return i
		}
	}

	log.Fatal("Couldn't find the seat")
	return -1
}

func main() {
	fmt.Println(findMissingSeat(16))
}
