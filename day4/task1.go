package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type passport struct {
	Byr, Iyr, Eyr, Hgt, Hcl, Ecl, Pid, Cid string
}

func (p passport) valid() bool {
	if p.Byr == "" || p.Iyr == "" || p.Eyr == "" || p.Hgt == "" ||
		p.Hcl == "" || p.Ecl == "" || p.Pid == "" {
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
