package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kenthklui/adventofcode/util"
)

type race struct {
	time, record int
}

func (r race) waysToWin() int {
	var i, j int
	for i = 1; i < r.time; i++ {
		distance := (r.time - i) * i
		if distance > r.record {
			break
		}
	}
	for j = r.time - 1; j > 0; j-- {
		distance := (r.time - j) * j
		if distance > r.record {
			break
		}
	}

	return j - i + 1
}

func parseRaces(timeStr, recordStr string) []race {
	intRegex := regexp.MustCompile(`-?\d+`)

	timeMatches := intRegex.FindAllString(timeStr, -1)
	recordMatches := intRegex.FindAllString(recordStr, -1)

	races := make([]race, len(timeMatches))
	for i := range races {
		if time, err := strconv.Atoi(timeMatches[i]); err == nil {
			races[i].time = time
		} else {
			panic(err)
		}
		if record, err := strconv.Atoi(recordMatches[i]); err == nil {
			races[i].record = record
		} else {
			panic(err)
		}
	}
	return races
}

func main() {
	input := util.StdinReadlines()
	races := parseRaces(input[0], input[1])
	product := 1
	for _, r := range races {
		product *= r.waysToWin()
	}
	fmt.Println(product)
}
