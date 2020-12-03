package main

import (
	"bufio"
	"fmt"
	"os"
)

func countTrees() int {
	treeMap := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		treeMap = append(treeMap, scanner.Text())
	}

	x := 0
	width := len(treeMap[0])
	trees := 0
	for y := 0; y < len(treeMap); y++ {
		if treeMap[y][x] == '#' {
			trees++
		}

		x = (x + 3) % width
	}

	return trees
}

func main() {
	fmt.Printf("%d\n", countTrees())
}
