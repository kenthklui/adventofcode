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

func (f field) String() string {
	return fmt.Sprintf("{%s: [%d, %d], [%d, %d]}", f.Name, f.Min1, f.Max1, f.Min2, f.Max2)
}

func (f field) Validate(value int) bool {
	// fmt.Printf("Validating %d against %q...", value, f)
	valid := ((value >= f.Min1 && value <= f.Max1) || (value >= f.Min2 && value <= f.Max2))
	// fmt.Printf("...%t\n", valid)
	return valid
}

type fieldList map[string]field

type ticket struct {
	Values []int
}

func (t ticket) String() string {
	strs := make([]string, len(t.Values))
	i := 0
	for _, v := range t.Values {
		strs[i] = strconv.Itoa(v)
		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ","))
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

type options map[int]bool
type candidates map[string]options

func (o options) String() string {
	strs := make([]string, len(o))
	i := 0
	for j := range o {
		strs[i] = strconv.Itoa(j)
		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ","))
}

func newCandidates(fieldCount int, fl fieldList) candidates {
	c := make(candidates)
	for name := range fl {
		o := make(options)
		for i := 0; i < fieldCount; i++ {
			o[i] = true
		}

		c[name] = o
	}

	// fmt.Println(c)
	return c
}

func (c candidates) Eliminate(name string, i int, fieldMap []string) {
	if _, ok := c[name][i]; ok {
		// fmt.Printf("Eliminated %d from %s\n", i, name)
		delete(c[name], i)

		if len(c[name]) == 1 {
			for solution := range c[name] { // should only be one
				// fmt.Printf("Only %d left for %s\n", solution, name)
				for otherName := range c {
					if otherName != name {
						c.Eliminate(otherName, solution, fieldMap)
					}
				}

				fieldMap[solution] = name
			}
		}
	}
	// fmt.Println(c)
}

func determineFieldMap(nearbyTickets []ticket, fl fieldList) []string {
	c := newCandidates(len(nearbyTickets[0].Values), fl)

	fieldMap := make([]string, len(nearbyTickets[0].Values))
	for _, t := range nearbyTickets {
		// fmt.Printf("Parsing ticket: %q\n", t)
		for i, v := range t.Values {
			for name, f := range fl {
				if !f.Validate(v) {
					c.Eliminate(name, i, fieldMap)
				}
			}
		}
	}

	// fmt.Println(c)
	// fmt.Println(fieldMap)

	return fieldMap
}

func departureProduct(t ticket, fieldMap []string) int {
	product := 1
	for i, name := range fieldMap {
		if strings.Index(name, "departure") == 0 {
			product *= t.Values[i]
		}
	}

	return product
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

func validateNearbyTickets(nearbyTickets []ticket, fl fieldList) []ticket {
	validTickets := make([]ticket, 0)
	for _, t := range nearbyTickets {
		fitsAllFields := true
		for _, v := range t.Values {
			fitsAllFields = fitsAllFields && fitsAnyField(fl, v)
		}
		if fitsAllFields {
			validTickets = append(validTickets, t)
		}
	}

	return validTickets
}

func main() {
	input := readInput()
	fl, myTicket, nearbyTickets := splitSections(input)
	validTickets := validateNearbyTickets(nearbyTickets, fl)
	fieldMap := determineFieldMap(validTickets, fl)
	product := departureProduct(myTicket, fieldMap)

	fmt.Println(product)
}
