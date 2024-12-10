package util

import (
	"bufio"
	"fmt"
	"io"
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
	reader := bufio.NewReader(os.Stdin)
	fullLine := make([]byte, 0, reader.Size())
	for {
		if line, isPrefix, err := reader.ReadLine(); err == nil {
			fullLine = append(fullLine, line...)
			if !isPrefix {
				input = append(input, string(fullLine))
				fullLine = fullLine[:0]
			}
		} else if err == io.EOF {
			break
		} else {
			panic(err)
		}
	}

	return
}
