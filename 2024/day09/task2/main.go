package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type chunk struct {
	pos, size int
}

type chunkHeap []chunk

func (c chunkHeap) Len() int           { return len(c) }
func (c chunkHeap) Less(i, j int) bool { return c[i].pos < c[j].pos }
func (c chunkHeap) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c *chunkHeap) Push(x interface{}) {
	*c = append(*c, x.(chunk))
}
func (c *chunkHeap) Pop() interface{} {
	old := *c
	n := len(old)
	x := old[n-1]
	*c = old[:n-1]
	return x
}

func parse(line string) ([]chunk, []chunkHeap) {
	files, spaces := []chunk{}, make([]chunkHeap, 10)
	totalSize := 0
	for i, c := range line {
		size := int(c - '0')
		if i%2 == 0 {
			files = append(files, chunk{totalSize, size})
		} else {
			spaces[size] = append(spaces[size], chunk{totalSize, size})
		}
		totalSize += size
	}

	for i := range spaces {
		heap.Init(&spaces[i])
	}
	return files, spaces
}

func checksum(files []chunk) int {
	checksum := 0
	for id, file := range files {
		// sum(n, n+1, ..., n+k) = sum(1, ..., n+k) - sum(1, ..., n-1)
		fileEnd := file.pos + file.size - 1
		checksum += (fileEnd*(fileEnd+1) - (file.pos-1)*file.pos) / 2 * id
	}
	return checksum
}

func solve(input []string) (output string) {
	files, spaces := parse(input[0])
	for i := len(files) - 1; i > 0; i-- {
		earliest := files[i]
		for j := files[i].size; j < len(spaces); j++ {
			if len(spaces[j]) > 0 && earliest.pos > spaces[j][0].pos {
				earliest = spaces[j][0]
			}
		}
		if earliest.pos < files[i].pos {
			files[i].pos = earliest.pos

			heap.Pop(&spaces[earliest.size])
			earliest.size -= files[i].size
			earliest.pos += files[i].size
			heap.Push(&spaces[earliest.size], earliest)
		}
	}
	return strconv.Itoa(checksum(files))
}

func main() {
	input := util.StdinReadlines()
	solution := solve(input)
	fmt.Println(solution)
}
