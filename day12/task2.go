package main

import (
	"bufio"
	"fmt"
	"os"
)

func readDirections() int {
	waypointX := 10
	waypointY := 1

	var x, y int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var action string
		var value int
		n, err := fmt.Sscanf(scanner.Text(), "%1s%d", &action, &value)
		if err != nil || n != 2 {
			panic("Failed Sscanf")
		}

		switch action {
		case "N":
			waypointY += value
		case "S":
			waypointY -= value
		case "E":
			waypointX += value
		case "W":
			waypointX -= value
		case "L":
			switch value {
			case 90:
				waypointX, waypointY = -waypointY, waypointX
			case 180:
				waypointX, waypointY = -waypointX, -waypointY
			case 270:
				waypointX, waypointY = waypointY, -waypointX
			default:
				panic("Odd L value")
			}
		case "R":
			switch value {
			case 90:
				waypointX, waypointY = waypointY, -waypointX
			case 180:
				waypointX, waypointY = -waypointX, -waypointY
			case 270:
				waypointX, waypointY = -waypointY, waypointX
			default:
				panic("Odd R value")
			}
		case "F":
			x += value * waypointX
			y += value * waypointY
		default:
			panic("Odd instruction")
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return x + y
}

func main() {
	fmt.Println(readDirections())
}
