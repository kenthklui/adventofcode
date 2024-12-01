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

func ConvertScreen(screen [][]bool) []string {
	output := make([]string, len(screen))
	for i, row := range screen {
		var b strings.Builder
		for _, p := range row {
			if p {
				b.WriteRune('█')
			} else {
				b.WriteRune('░')
			}
		}
		output[i] = b.String()
	}
	return output
}

func PrintScreen(screen [][]bool) {
	PrintStrings(ConvertScreen(screen))
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
