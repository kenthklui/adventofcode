package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type obj interface {
	name() string
	size() int
	isDir() bool
	parent() *dir
}

type dir struct {
	nameStr   string
	parentDir *dir
	children  map[string]obj
}

func NewDir(nameStr string, parentDir *dir) *dir {
	return &dir{
		nameStr:   nameStr,
		parentDir: parentDir,
		children:  make(map[string]obj),
	}
}

func (d *dir) name() string { return d.nameStr }
func (d *dir) size() int {
	sum := 0
	for _, o := range d.children {
		sum += o.size()
	}
	return sum
}
func (d *dir) isDir() bool  { return true }
func (d *dir) parent() *dir { return d.parentDir }

type file struct {
	nameStr   string
	parentDir *dir
	bytes     int
}

func NewFile(nameStr string, parentDir *dir, size int) *file {
	return &file{
		nameStr:   nameStr,
		parentDir: parentDir,
		bytes:     size,
	}
}

func (f *file) name() string { return f.nameStr }
func (f *file) size() int    { return f.bytes }
func (f *file) isDir() bool  { return false }
func (f *file) parent() *dir { return f.parentDir }

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

func findClosestAbove(currDir *dir, sizeThreshold int) int {
	closest := -1

	currSize := currDir.size()
	if currSize >= sizeThreshold {
		closest = currSize
	}
	for _, o := range currDir.children {
		if o.isDir() {
			d := o.(*dir)
			childClosest := findClosestAbove(d, sizeThreshold)
			if childClosest > 0 && childClosest < closest {
				closest = childClosest
			}
		}
	}
	return closest
}

func parseInput(input []string) *dir {
	root := NewDir("/", nil)
	currDir := root

	lineNum := 0
	for lineNum < len(input) {
		line := input[lineNum]
		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "$":
			switch tokens[1] {

			case "cd":
				switch tokens[2] {
				case "/":
					currDir = root
				case "..":
					currDir = currDir.parent()
				default:
					if child, ok := currDir.children[tokens[2]]; !ok {
						panic("Child not found")
					} else if child.isDir() {
						currDir = child.(*dir)
					} else {
						panic("cd into non directory")
					}
				}
				lineNum++

			case "ls":
				lineNum++
				for lineNum < len(input) {
					line := input[lineNum]
					tokens := strings.Split(line, " ")

					if tokens[0] == "$" {
						break
					}

					switch tokens[0] {
					case "dir":
						if _, ok := currDir.children[tokens[1]]; !ok {
							currDir.children[tokens[1]] = NewDir(tokens[1], currDir)
						}
					default: // Number, will be filesize
						size, err := strconv.Atoi(tokens[0])
						if err != nil {
							panic(err)
						}
						if _, ok := currDir.children[tokens[1]]; !ok {
							currDir.children[tokens[1]] = NewFile(tokens[1], currDir, size)
						}
					}
					lineNum++
				}
			}
		default:
			panic("Not a command")
		}
	}

	return root
}

func main() {
	input := readInput()
	root := parseInput(input)

	totalSize := 70000000
	neededSpace := 30000000

	freeSpace := totalSize - root.size()
	toFree := neededSpace - freeSpace
	fmt.Println(findClosestAbove(root, toFree))
}
