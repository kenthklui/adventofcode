package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintStrings(strArray []string) {
	for _, str := range strArray {
		fmt.Println(str)
	}
}

func PrintScreen(screen [][]bool) {
	var b strings.Builder
	for _, row := range screen {
		for _, p := range row {
			if p {
				b.WriteRune('█')
			} else {
				b.WriteRune('░')
			}
		}
		b.WriteRune('\n')
	}
	fmt.Printf(b.String())
}

func StdinReadlines() (input []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for input = make([]string, 0); scanner.Scan(); {
		input = append(input, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}
