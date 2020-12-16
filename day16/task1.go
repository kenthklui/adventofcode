package main

import (
	"bufio"
	"fmt"
	"os"
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

type field struct {
	Name                   string
	Min1, Max1, Min2, Max2 int
}

type fieldList map[string]field

type ticket struct {
	Values []int
}

func newTicket(line string) ticket {
	valueStrs := strings.Split(line, ",")
	values := make([]int, len(valueStrs))

	var err error
	for i, s := range valueStrs {
		if values[i], err = strconv.Atoi(s); err != nil {
			panic(err)
		}
	}

	return ticket{values}
}

func splitSections(input []string) (map[string]field, ticket, []ticket) {
	fl := make(fieldList)

	var i int
	var line string

	for i, line = range input {
		if line == "" {
			break
		}

		f := field{}
		tokens := strings.Split(line, ":")
		f.Name = tokens[0]

		n, err := fmt.Sscanf(tokens[1], "%d-%d or %d-%d", &f.Min1, &f.Max1, &f.Min2, &f.Max2)
		if err != nil || n != 4 {
			panic(err)
		}

		f.Name = strings.TrimSuffix(f.Name, ":")

		fl[f.Name] = f
	}

	i++ // skip empty line
	i++ // skip "your ticket:"
	myTicket := newTicket(input[i])
	i++

	i++ // skip empty line
	i++ // skip "nearby tickets:"
	nearbyTickets := make([]ticket, len(input)-i)
	for i, line = range input[i:] {
		nearbyTickets[i] = newTicket(line)
	}

	return fl, myTicket, nearbyTickets
}

func fitsAnyField(fl fieldList, input int) bool {
	for _, f := range fl {
		if input >= f.Min1 && input <= f.Max1 {
			return true
		}
		if input >= f.Min2 && input <= f.Max2 {
			return true
		}
	}

	return false
}

func validateNearbyTickets(nearbyTickets []ticket, fl fieldList) int {
	var sum int
	for _, t := range nearbyTickets {
		for _, v := range t.Values {
			if !fitsAnyField(fl, v) {
				sum += v
			}
		}
	}

	return sum
}

func main() {
	input := readInput()
	fl, _, nearbyTickets := splitSections(input)
	output := validateNearbyTickets(nearbyTickets, fl)
	fmt.Println(output)
}
