package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type child struct {
	Quantity int
	Name     string
}

// bag type <key> contains <child.quantity> of <child.name>
type childMap map[string][]child

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseRules() childMap {
	ruleRegex := regexp.MustCompile("((([1-9]+ )?[a-z]+ [a-z]+) bags?)")
	cm := make(childMap)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		matches := ruleRegex.FindAllString(scanner.Text(), -1)

		parentName := strings.TrimSuffix(matches[0], " bags")

		if matches[1] == "no other bags" {
			cm[parentName] = nil
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

			c := child{quantity, childName}
			if childList, ok := cm[parentName]; ok {
				cm[parentName] = append(childList, c)
			} else {
				cm[parentName] = []child{c}
			}
		}
	}
	failOnErr(scanner.Err())

	return cm
}

func countChildren(parentName string, cm childMap) int {
	count := 0
	if children, ok := cm[parentName]; ok {
		if children != nil {
			for _, child := range children {
				count += child.Quantity * (1 + countChildren(child.Name, cm))
			}
		}
	} else {
		log.Fatalf("Failed to find %q in child map\n", parentName)
	}

	return count
}

func main() {
	rules := parseRules()

	fmt.Println(countChildren("shiny gold", rules))
}
