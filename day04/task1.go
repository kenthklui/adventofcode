package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type passport map[string]string

func (p passport) valid() bool {
	requiredFields := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}
	for _, field := range requiredFields {
		if value, ok := p[field]; !ok || value == "" {
			return false
		}
	}

	return true
}

func NewPassport(str string) passport {
	var p passport = make(map[string]string)

	str = strings.Trim(str, "\n")
	str = strings.ReplaceAll(str, "\n", " ")
	for _, entry := range strings.Split(str, " ") {
		var field, value string
		if _, err := fmt.Sscanf(entry, "%3s:%s", &field, &value); err == nil {
			p[field] = value
		} else {
			log.Fatalf("Parse error: %s", err.Error())
		}
	}

	return p
}

func countValidPassports() int {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	sectionStrings := strings.Split(string(b), "\n\n")

	valid := 0
	for _, str := range sectionStrings {
		passport := NewPassport(str)
		if passport.valid() {
			valid++
		}
	}

	return valid
}

func main() {
	fmt.Println(countValidPassports())
}
