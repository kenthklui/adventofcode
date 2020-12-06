package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type passport map[string]string

func (p passport) valid() bool {
	byr, err := strconv.Atoi(p["byr"])
	if err != nil || byr < 1920 || byr > 2002 {
		return false
	}

	iyr, err := strconv.Atoi(p["iyr"])
	if err != nil || iyr < 2010 || iyr > 2020 {
		return false
	}

	eyr, err := strconv.Atoi(p["eyr"])
	if err != nil || eyr < 2020 || eyr > 2030 {
		return false
	}

	var hgtInt int
	if _, err := fmt.Sscanf(p["hgt"], "%dcm", &hgtInt); err == nil {
		if hgtInt < 150 || hgtInt > 193 {
			return false
		}
	} else if _, err := fmt.Sscanf(p["hgt"], "%din", &hgtInt); err == nil {
		if hgtInt < 59 || hgtInt > 76 {
			return false
		}
	} else {
		return false
	}

	if matched, err := regexp.MatchString("^#[0-9a-f]{6}$", p["hcl"]); !matched || err != nil {
		return false
	}

	validEcl := map[string]int{"amb": 1, "blu": 1, "brn": 1, "gry": 1, "grn": 1, "hzl": 1, "oth": 1}
	if _, ok := validEcl[p["ecl"]]; !ok {
		return false
	}

	if matched, err := regexp.MatchString("^[0-9]{9}$", p["pid"]); !matched || err != nil {
		return false
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
