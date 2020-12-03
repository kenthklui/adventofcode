package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		var first, second int
		var char, password string

		n, err := fmt.Sscanf(scanner.Text(), "%d-%d %1s: %s", &first, &second, &char, &password)
		failOnErr(err)

		if n != 4 {
			log.Fatal("failed to parse 4 items")
		}

		if (password[first-1] == char[0]) != (password[second-1] == char[0]) {
			valid++
		}
	}
	failOnErr(scanner.Err())

	return valid
}

func main() {
	fmt.Println(validateLines())
}
