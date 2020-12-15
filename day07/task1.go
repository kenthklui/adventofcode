package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// bag type <key> can be placed in bag type <value>
type parentMap map[string][]string

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseRules() parentMap {
	ruleRegex := regexp.MustCompile("((([1-9]+ )?[a-z]+ [a-z]+) bags?)")
	pm := make(parentMap)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		matches := ruleRegex.FindAllString(scanner.Text(), -1)

		parentName := strings.TrimSuffix(matches[0], " bags")

		if matches[1] == "no other bags" {
			// Nobody has this bag as a parent
			continue
		}

		for _, childStr := range matches[1:] {
			var quantity int
			var firstName, secondName, bags string
			n, err := fmt.Sscanf(childStr, "%d %s %s %s", &quantity, &firstName, &secondName, &bags)
			failOnErr(err)
			if n != 4 {
				log.Fatal("Failed to parse line: %q", childStr)
			}

			childName := fmt.Sprintf("%s %s", firstName, secondName)
			if parentList, ok := pm[childName]; ok {
				pm[childName] = append(parentList, parentName)
			} else {
				pm[childName] = []string{parentName}
			}
		}
	}
	failOnErr(scanner.Err())

	return pm
}

func findParents(childName string, pm parentMap, isParent map[string]bool) {
	if parents, ok := pm[childName]; ok {
		if parents == nil {
			return
		}

		for _, parent := range parents {
			isParent[parent] = true
			findParents(parent, pm, isParent)
		}
	} else {
		return
	}
}

func main() {
	rules := parseRules()

	isParent := make(map[string]bool)
	findParents("shiny gold", rules, isParent)

	fmt.Println(len(isParent))
}
