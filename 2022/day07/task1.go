package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dir struct {
	name           string
	parent         *dir
	subdirectories map[string]*dir
	files          map[string]*file
}

func NewDir(name string, parent *dir) *dir {
	return &dir{
		name:           name,
		parent:         parent,
		subdirectories: make(map[string]*dir),
		files:          make(map[string]*file),
	}
}

func (d *dir) size() int {
	sum := 0
	for _, subdir := range d.subdirectories {
		sum += subdir.size()
	}
	for _, file := range d.files {
		sum += file.size
	}
	return sum
}

type file struct {
	name   string
	parent *dir
	size   int
}

func NewFile(name string, parent *dir, size int) *file {
	return &file{
		name:   name,
		parent: parent,
		size:   size,
	}
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

func sumBelow(currDir *dir, sizeThreshold int) int {
	sum := 0
	size := currDir.size()
	if size <= sizeThreshold {
		sum += size
	}
	for _, subdir := range currDir.subdirectories {
		sum += sumBelow(subdir, sizeThreshold)
	}
	return sum
}

func parseInput(input []string) *dir {
	root := NewDir("/", nil)
	currDir := root

	for _, line := range input {
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "$":
			switch tokens[1] {
			case "cd":
				switch tokens[2] {
				case "/":
					currDir = root
				case "..":
					currDir = currDir.parent
				default:
					currDir, _ = currDir.subdirectories[tokens[2]]
				}
			}
		case "dir":
			currDir.subdirectories[tokens[1]] = NewDir(tokens[1], currDir)
		default:
			size, _ := strconv.Atoi(tokens[0])
			currDir.files[tokens[1]] = NewFile(tokens[1], currDir, size)
		}
	}

	return root
}

func main() {
	input := readInput()
	root := parseInput(input)

	fmt.Println(sumBelow(root, 100000))
}
