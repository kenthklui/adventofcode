package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

type iMap [3]int

type production struct {
	seeds []int
	// SeedToSoil, SoilToFert, FertToWater, WaterToLight, LightToTemp, TempToHumid, HumidToLoc
	Maps [7][]iMap
}

func convert(input int, ims []iMap) int {
	output := input
	for _, im := range ims {
		diff := input - im[1]
		if diff >= 0 && diff < im[2] {
			output += (im[0] - im[1])
		}
	}

	return output
}

func (p production) getLocations() []int {
	locations := make([]int, len(p.seeds))
	for i, num := range p.seeds {
		for _, ims := range p.Maps {
			num = convert(num, ims)
		}
		locations[i] = num
	}

	return locations
}

func parse(input []string) production {
	var err error

	_, after, _ := strings.Cut(input[0], ": ")
	seedStrs := strings.Split(after, " ")
	seeds := make([]int, len(seedStrs))
	for i, str := range seedStrs {
		if seeds[i], err = strconv.Atoi(str); err != nil {
			panic(err)
		}
	}

	var maps [7][]iMap
	lineNum := 3
	for i := 0; i < len(maps); i++ {
		maps[i] = make([]iMap, 0)
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
			maps[i] = append(maps[i], im)
		}
		lineNum += len(maps[i]) + 2
	}

	return production{seeds, maps}
}

func main() {
	input := util.StdinReadlines()
	p := parse(input)
	min := 1000000000
	for _, location := range p.getLocations() {
		if location < min {
			min = location
		}
	}
	fmt.Println(min)
}
