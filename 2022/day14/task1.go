package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

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

type rockMap map[int][]int

func (rm rockMap) addRock(x, y int) {
	if yList, ok := rm[x]; ok {
		i := sort.SearchInts(yList, y)
		if i == len(yList) {
			rm[x] = append(yList, y)
		} else if yList[i] == y {
			return
		} else {
			yList = append(yList[:i+1], yList[i:]...)
			yList[i] = y
			rm[x] = yList
		}
	} else {
		rm[x] = []int{y}
	}
}

func (rm rockMap) addSand() int {
	sandX := 500
	sandY := 0

	sand := 0
	for {
		yList, ok := rm[sandX]
		if !ok { // to the abyss you go
			break
		}
		yIndex := sort.SearchInts(yList, sandY)
		if yIndex >= len(yList) { // to the abyss you go
			break
		}
		landOnY := yList[yIndex]

		leftList, ok := rm[sandX-1]
		if !ok { // to the abyss you go
			break
		}
		leftIndex := sort.SearchInts(leftList, landOnY)
		if leftIndex >= len(leftList) { // to the abyss you go
			break
		} else if leftList[leftIndex] != landOnY { // left not occupied
			sandX--
			sandY = landOnY
			continue
		}

		rightList, ok := rm[sandX+1]
		if !ok { // to the abyss you go
			break
		}
		rightIndex := sort.SearchInts(rightList, landOnY)
		if rightIndex >= len(rightList) { // to the abyss you go
			break
		} else if rightList[rightIndex] != landOnY { // right not occupied
			sandX++
			sandY = landOnY
			continue
		}

		// both left and right occupied
		rm.addRock(sandX, landOnY-1)
		sand++
		sandX, sandY = 500, 0
	}

	return sand
}

func parseInput(input []string) rockMap {
	rm := make(rockMap)
	for _, line := range input {
		rocks := strings.Split(line, " -> ")
		point := strings.Split(rocks[0], ",")
		x1, _ := strconv.Atoi(point[0])
		y1, _ := strconv.Atoi(point[1])

		for _, rock := range rocks[1:] {
			point = strings.Split(rock, ",")
			x2, _ := strconv.Atoi(point[0])
			y2, _ := strconv.Atoi(point[1])

			if x1 < x2 {
				for x := x1; x <= x2; x++ {
					rm.addRock(x, y1)
				}
			} else if x1 > x2 {
				for x := x2; x <= x1; x++ {
					rm.addRock(x, y1)
				}
			} else if y1 < y2 {
				for y := y1; y <= y2; y++ {
					rm.addRock(x1, y)
				}
			} else if y1 > y2 {
				for y := y2; y <= y1; y++ {
					rm.addRock(x1, y)
				}
			} else {
				panic("Two rocks are the same point")
			}
			x1 = x2
			y1 = y2
		}
	}
	return rm
}

func main() {
	input := readInput()
	rm := parseInput(input)

	fmt.Println(rm.addSand())
}
