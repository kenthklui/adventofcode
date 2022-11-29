package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func algorithm(s string) []int {
	algo := make([]int, 512)

	for i, r := range s {
		switch r {
		case '.':
			algo[i] = 0
		case '#':
			algo[i] = 1
		default:
			panic("Invalid algorithm character")
		}
	}

	return algo
}

type image struct {
	width, height int
	pixels        []int
}

func NewImage(width, height int) image {
	return image{
		width:  width,
		height: height,
		pixels: make([]int, width*height),
	}
}

func (i image) print() {
	var b strings.Builder

	index := 0
	for h := 0; h < i.height; h++ {
		for w := 0; w < i.width; w++ {
			switch i.pixels[index] {
			case 1:
				b.WriteString("#")
			case 0:
				b.WriteString(".")
			default:
				panic("Invalid pixel")
			}
			index++
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")

	fmt.Printf(b.String())
}

// I'm sure there's a faster way to do this with FFTs
func (i image) getPixelValue(x, y, step int, algo []int) int {
	maskOffset := 1 // I bet this gets bigger in task 2
	pixelValue := 0
	for maskY := y - maskOffset; maskY <= y+maskOffset; maskY++ {
		for maskX := x - maskOffset; maskX <= x+maskOffset; maskX++ {
			pixelValue *= 2

			var light int
			if maskX < 0 || maskY < 0 || maskX >= i.width || maskY >= i.height {
				// AHHH IT'S A TRAP
				// algo[0] == 1 means the color outside the main canvas can also flip on odd steps
				if algo[0] == 1 {
					light = step % 2
				} else {
					light = 0
				}
			} else {
				index := maskY*i.width + maskX
				light = i.pixels[index]
			}

			pixelValue += light
		}
	}

	return pixelValue
}

func (i image) countLights() int {
	count := 0
	for _, p := range i.pixels {
		count += p
	}

	return count
}

func (i image) enhance(algo []int, step int) image {
	newImage := NewImage(i.width+2, i.height+2)

	for y := 0; y < newImage.height; y++ {
		for x := 0; x < newImage.width; x++ {
			oldX, oldY := x-1, y-1
			pixelValue := i.getPixelValue(oldX, oldY, step, algo)

			index := y*newImage.width + x
			newImage.pixels[index] = algo[pixelValue]
		}
	}

	return newImage
}

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

func parseInput(input []string) ([]int, image) {
	algo := algorithm(input[0])

	width := len(input[2])
	height := len(input) - 2
	image := NewImage(width, height)

	index := 0
	for _, line := range input[2:] {
		for _, r := range line {
			switch r {
			case '.':
				image.pixels[index] = 0
			case '#':
				image.pixels[index] = 1
			default:
				panic("Invalid image character")
			}
			index++
		}
	}

	return algo, image
}

func main() {
	input := readInput()
	algo, image := parseInput(input)
	// image.print()

	for i := 0; i < 2; i++ {
		image = image.enhance(algo, i)
	}

	fmt.Println(image.countLights())
}
