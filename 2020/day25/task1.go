package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var mod int = 20201227
var subject int = 7

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

func readPublicKeys(input []string) (int, int, int, int) {
	cardKey, _ := strconv.Atoi(input[0])
	doorKey, _ := strconv.Atoi(input[1])

	var cardLoop, doorLoop int
	for value, loop := 1, 1; loop < 100000000; loop++ {
		value *= subject
		value %= mod

		if cardLoop == 0 && value == cardKey {
			cardLoop = loop
		}
		if doorLoop == 0 && value == doorKey {
			doorLoop = loop
		}

		if cardLoop != 0 && doorLoop != 0 {
			return cardKey, cardLoop, doorKey, doorLoop
		}
	}

	panic("Not found")
}

func buildEncryptionKey(cardKey, cardLoop, doorKey, doorLoop int) int {
	value := 1

	for i := 0; i < int(doorLoop); i++ {
		value *= cardKey
		value %= mod
	}

	return value
}

func main() {
	input := readInput()
	cardKey, cardLoop, doorKey, doorLoop := readPublicKeys(input)
	encryptionKey := buildEncryptionKey(cardKey, cardLoop, doorKey, doorLoop)
	fmt.Println(encryptionKey)
}
