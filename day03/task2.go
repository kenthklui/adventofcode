package main

import (
	"bufio"
	"fmt"
	"os"
)

func getTreeMap() []string {
	treeMap := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		treeMap = append(treeMap, scanner.Text())
	}

	return treeMap
}

func countTrees(treeMap []string, xDelta, yDelta int) int {

	x := 0
	width := len(treeMap[0])
	trees := 0
	for y := 0; y < len(treeMap); y += yDelta {
		if treeMap[y][x] == '#' {
			trees++
		}

		x = (x + xDelta) % width
	}

	return trees
}

func main() {
	treeProduct := 1
	treeMap := getTreeMap()

	treeProduct *= countTrees(treeMap, 1, 1)
	treeProduct *= countTrees(treeMap, 3, 1)
	treeProduct *= countTrees(treeMap, 5, 1)
	treeProduct *= countTrees(treeMap, 7, 1)
	treeProduct *= countTrees(treeMap, 1, 2)

	fmt.Printf("%d\n", treeProduct)
}
