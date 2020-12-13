package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func nextBus() int {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	time, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	var bestBus int
	minWait := time
	scanner.Scan()
	for _, idStr := range strings.Split(scanner.Text(), ",") {
		if idStr == "x" {
			continue
		}

		if busID, err := strconv.Atoi(idStr); err == nil {
			busWait := busID - time%busID
			if minWait > busWait {
				minWait = busWait
				bestBus = busID
			}
		} else {
			panic(err)
		}
	}

	return minWait * bestBus
}

func main() {
	fmt.Println(nextBus())
}
