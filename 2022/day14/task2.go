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

func (rm rockMap) addSand(floorY int) int {
	sandX := 500
	sandY := 0

	sand := 0
	for {
		yList, ok := rm[sandX]
		if !ok { // hit the floor
			rm.addRock(sandX, floorY-1)
			sand++
			continue
		}
		yIndex := sort.SearchInts(yList, sandY)
		if yIndex >= len(yList) { // hit the floor
			rm.addRock(sandX, floorY-1)
			sand++
			continue
		}
		landOnY := yList[yIndex]

		if sandX == 500 && landOnY == 0 { // source was already filled, done
			break
		}

		leftList, ok := rm[sandX-1]
		if !ok { // hit the floor
			rm.addRock(sandX-1, floorY-1)
			sand++
			continue
		}
		leftIndex := sort.SearchInts(leftList, landOnY)
		if leftIndex >= len(leftList) { // hit the floor
			rm.addRock(sandX-1, floorY-1)
			sand++
			continue
		} else if leftList[leftIndex] != landOnY { // left not occupied
			sandX--
			sandY = landOnY
			continue
		}

		rightList, ok := rm[sandX+1]
		if !ok { // hit the floor
			rm.addRock(sandX+1, floorY-1)
			sand++
			continue
		}
		rightIndex := sort.SearchInts(rightList, landOnY)
		if rightIndex >= len(rightList) { // hit the floor
			rm.addRock(sandX+1, floorY-1)
			sand++
			continue
		} else if rightList[rightIndex] != landOnY { // right not occupied
			sandX++
			sandY = landOnY
			continue
		}

		rm.addRock(sandX, landOnY-1)
		sand++
		sandX, sandY = 500, 0
	}

	return sand
}

func parseInput(input []string) (rockMap, int) {
	rm := make(rockMap)
	maxY := 0
	for _, line := range input {
		rocks := strings.Split(line, " -> ")
		point := strings.Split(rocks[0], ",")
		x1, _ := strconv.Atoi(point[0])
		y1, _ := strconv.Atoi(point[1])
		if y1 > maxY {
			maxY = y1
		}

		for _, rock := range rocks[1:] {
			point = strings.Split(rock, ",")
			x2, _ := strconv.Atoi(point[0])
			y2, _ := strconv.Atoi(point[1])
			if y2 > maxY {
				maxY = y2
			}

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
	return rm, maxY + 2
}

func main() {
	input := readInput()
	rm, floorY := parseInput(input)

	fmt.Println(rm.addSand(floorY))
}
