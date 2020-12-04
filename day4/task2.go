package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type passport struct {
	Byr, Iyr, Eyr, Hgt, Hcl, Ecl, Pid, Cid string
}

func (p passport) valid() bool {
	byr, err := strconv.Atoi(p.Byr)
	if err != nil || byr < 1920 || byr > 2002 {
		return false
	}

	iyr, err := strconv.Atoi(p.Iyr)
	if err != nil || iyr < 2010 || iyr > 2020 {
		return false
	}

	eyr, err := strconv.Atoi(p.Eyr)
	if err != nil || eyr < 2020 || eyr > 2030 {
		return false
	}

	var hgtInt int
	if _, err := fmt.Sscanf(p.Hgt, "%dcm", &hgtInt); err == nil {
		if hgtInt < 150 || hgtInt > 193 {
			return false
		}
	} else if _, err := fmt.Sscanf(p.Hgt, "%din", &hgtInt); err == nil {
		if hgtInt < 59 || hgtInt > 76 {
			return false
		}
	} else {
		return false
	}

	if matched, err := regexp.MatchString("^#[0-9a-f]{6}$", p.Hcl); !matched || err != nil {
		return false
	}

	validEcl := map[string]int{"amb": 1, "blu": 1, "brn": 1, "gry": 1, "grn": 1, "hzl": 1, "oth": 1}
	if _, ok := validEcl[p.Ecl]; !ok {
		return false
	}

	if matched, err := regexp.MatchString("^[0-9]{9}$", p.Pid); !matched || err != nil {
		return false
	}

	return true
}

func countValidPassports() int {
	valid := 0

	var currentPassport *passport
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentPassport != nil {
				if currentPassport.valid() {
					valid++
				}

				currentPassport = nil
			}
			continue
		}

		if currentPassport == nil {
			currentPassport = &passport{}
		}

		for _, entry := range strings.Split(line, " ") {
			var field, value string
			if _, err := fmt.Sscanf(entry, "%3s:%s", &field, &value); err == nil {
				switch field {
				case "byr":
					currentPassport.Byr = value
				case "iyr":
					currentPassport.Iyr = value
				case "eyr":
					currentPassport.Eyr = value
				case "hgt":
					currentPassport.Hgt = value
				case "hcl":
					currentPassport.Hcl = value
				case "ecl":
					currentPassport.Ecl = value
				case "pid":
					currentPassport.Pid = value
				case "cid":
					currentPassport.Cid = value
				default:
					log.Fatalf("Invalid field: %s", field)
				}
			} else {
				log.Fatalf("Parse error: %s", err.Error())
			}
		}
	}
	if currentPassport != nil {
		if currentPassport.valid() {
			valid++
		}
	}

	return valid
}

func main() {
	fmt.Println(countValidPassports())
}
