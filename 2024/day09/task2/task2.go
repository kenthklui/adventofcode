package main

import (
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type chunk struct {
	id, pos, size int
}

func parse(line string) ([]chunk, []chunk) {
	files, space := []chunk{}, []chunk{}
	totalSize := 0
	for i, c := range line {
		size := int(c - '0')
		if i%2 == 0 {
			files = append(files, chunk{i / 2, totalSize, size})
		} else {
			space = append(space, chunk{0, totalSize, size})
		}
		totalSize += size
	}

	return files, space
}

func checksum(files []chunk) int {
	checksum := 0
	for _, file := range files {
		// sum(n, n+1, ..., n+k) = sum(1, ..., n+k) - sum(1, ..., n-1)
		fileEnd := file.pos + file.size - 1
		checksum += (fileEnd*(fileEnd+1) - (file.pos-1)*file.pos) / 2 * file.id
	}
	return checksum
}

func solve(input []string) (output string) {
	files, space := parse(input[0])
	for i := len(files) - 1; i > 0; i-- {
		for j := range space {
			if files[i].pos < space[j].pos {
				break
			}
			if files[i].size <= space[j].size {
				files[i].pos = space[j].pos
				space[j].size -= files[i].size
				space[j].pos += files[i].size
			}
		}
	}
	return strconv.Itoa(checksum(files))
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
