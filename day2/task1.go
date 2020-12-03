package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func validateLines() int {
	valid := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var low, high int
		var char, password string
		n, err := fmt.Sscanf(scanner.Text(), "%d-%d %1s: %s", &low, &high, &char, &password)
		failOnErr(err)

		if n != 4 {
			log.Fatal("failed to parse 4 items")
		}

		count := strings.Count(password, char)
		if count <= high && count >= low {
			valid++
		}
	}
	failOnErr(scanner.Err())

	return valid
}

func main() {
	fmt.Println(validateLines())
}
