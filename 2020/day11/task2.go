package main

import (
	"bufio"
	"fmt"
	"os"
)

type seat struct {
	Occupied, next bool
}

type layout [][]*seat

func (l layout) NeighborsOccupied(row, column int) int {
	occupied := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			for r, c := row+i, column+j; r >= 0 && c >= 0 && r < len(l) && c < len(l[0]); {
				if l[r][c] != nil {
					if l[r][c].Occupied {
						occupied++
					}
					break
				}

				r += i
				c += j
			}
		}
	}

	return occupied
}

func (l layout) SetNext() {
	for row, rowSeats := range l {
		for column, s := range rowSeats {
			if s == nil {
				continue
			}

			occupied := l.NeighborsOccupied(row, column)
			if occupied == 0 {
				s.next = true
			} else if occupied >= 5 {
				s.next = false
			} else {
				s.next = s.Occupied
			}
		}
	}
}

// returns true if changes will happen
func (l layout) Move() bool {
	changed := false
	for _, rowSeats := range l {
		for _, s := range rowSeats {
			if s != nil && s.next != s.Occupied {
				s.Occupied = s.next
				changed = true
			}
		}
	}

	return changed
}

func (l layout) String() string {
	str := ""
	for _, row := range l {
		for _, s := range row {
			if s == nil {
				str += "."
			} else if s.Occupied {
				str += "#"
			} else {
				str += "L"
			}
		}
		str += "\n"
	}

	return str
}

func (l layout) Occupied() int {
	occupied := 0
	for _, row := range l {
		for _, s := range row {
			if s != nil && s.Occupied {
				occupied++
			}
		}
	}

	return occupied
}

func readLayout() layout {
	l := make(layout, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		row := scanner.Text()
		rowLength := len(row)

		rowSeats := make([]*seat, rowLength)
		for i, r := range row {
			switch r {
			case 'L':
				rowSeats[i] = &seat{false, false}
			case '.':
				rowSeats[i] = nil
			}
		}
		l = append(l, rowSeats)
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return l
}

func main() {
	l := readLayout()

	changed := true
	for changed {
		// fmt.Println(l)
		l.SetNext()
		changed = l.Move()
	}

	fmt.Println(l.Occupied())
}
