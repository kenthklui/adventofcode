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

func parseRaces(timeStrFull, recordStrFull string) []race {
	intRegex := regexp.MustCompile(`-?\d+`)

	timeMatches := intRegex.FindAllString(timeStrFull, -1)
	recordMatches := intRegex.FindAllString(recordStrFull, -1)
	timeStr, recordStr := "", ""

	for i := range timeMatches {
		timeStr += timeMatches[i]
		recordStr += recordMatches[i]
	}

	var r race
	if time, err := strconv.Atoi(timeStr); err == nil {
		r.time = time
	} else {
		panic(err)
	}
	if record, err := strconv.Atoi(recordStr); err == nil {
		r.record = record
	} else {
		panic(err)
	}
	return []race{r}
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
