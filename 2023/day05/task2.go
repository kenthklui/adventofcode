package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type seedRange [2]int
type seedRanges []seedRange

// Assume all seed ranges are distinct and non overlapping and therefore sortable
func (srs seedRanges) Len() int           { return len(srs) }
func (srs seedRanges) Less(i, j int) bool { return srs[i][0] < srs[j][0] }
func (srs seedRanges) Swap(i, j int)      { srs[i], srs[j] = srs[j], srs[i] }

type iMap [3]int
type converter []iMap

// Assume all map ranges are distinct and non overlapping and therefore sortable
func (c converter) Len() int           { return len(c) }
func (c converter) Less(i, j int) bool { return c[i][1] < c[j][1] }
func (c converter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c converter) convert(input int) (int, int) {
	output, increment := input, -1

	for _, m := range c {
		diff := input - m[1]
		if diff < 0 {
			increment = -diff
			break
		} else if diff >= m[2] {
			continue
		}

		output = input + m[0] - m[1]
		increment = m[2] - diff
		break
	}

	return output, increment
}

type production struct {
	ranges seedRanges
	// SeedToSoil, SoilToFert, FertToWater, WaterToLight, LightToTemp, TempToHumid, HumidToLoc
	converters [7]converter
}

func (p production) getMinLocation(limit int) int {
	minLocation := limit
	seedValue, seedRangeIndex := 0, 0

	for seedRangeIndex < len(p.ranges) {
		if seedValue < p.ranges[seedRangeIndex][0] {
			seedValue = p.ranges[seedRangeIndex][0]
		}

		num, minIncrement := seedValue, limit
		for _, c := range p.converters {
			next, increment := c.convert(num)

			num = next
			if increment != -1 && increment < minIncrement {
				minIncrement = increment
			}
		}

		if minLocation > num {
			minLocation = num
		}

		seedValue = seedValue + minIncrement
		if seedValue > p.ranges[seedRangeIndex][1] {
			seedRangeIndex++
		}
	}

	return minLocation
}

func parse(input []string) production {
	var err error

	_, after, _ := strings.Cut(input[0], ": ")
	seedStrs := strings.Split(after, " ")
	ranges := make(seedRanges, 0)
	for i := 0; i < len(seedStrs); i += 2 {
		var seedMin, rangeSize int
		if seedMin, err = strconv.Atoi(seedStrs[i]); err != nil {
			panic(err)
		}
		if rangeSize, err = strconv.Atoi(seedStrs[i+1]); err != nil {
			panic(err)
		}
		ranges = append(ranges, seedRange{seedMin, seedMin + rangeSize - 1})
	}
	sort.Sort(ranges)

	var converters [7]converter
	lineNum := 3
	for i := range converters {
		converters[i] = make(converter, 0)
		for _, line := range input[lineNum:] {
			if line == "" {
				break
			}

			var im iMap
			for j, str := range strings.Split(line, " ")[:3] {
				if im[j], err = strconv.Atoi(str); err != nil {
					panic(err)
				}
			}
			converters[i] = append(converters[i], im)
		}
		lineNum += len(converters[i]) + 2

		sort.Sort(converters[i])
	}

	return production{ranges, converters}
}

func main() {
	input := util.StdinReadlines()
	p := parse(input)
	minLocation := p.getMinLocation(1000000000)
	fmt.Println(minLocation)
}
