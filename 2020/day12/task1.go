package main

import (
	"bufio"
	"fmt"
	"os"
)

func readDirections() int {
	heading := 90
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
			y += value
		case "S":
			y -= value
		case "E":
			x += value
		case "W":
			x -= value
		case "L":
			heading = (heading + 360 - value) % 360
		case "R":
			heading = (heading + value) % 360
		case "F":
			switch heading {
			case 0:
				y += value
			case 90:
				x += value
			case 180:
				y -= value
			case 270:
				x -= value
			default:
				panic("Odd heading")
			}
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
