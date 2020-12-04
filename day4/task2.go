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
	if byr, err := strconv.Atoi(p.Byr); err == nil {
		if byr < 1920 || byr > 2002 {
			log.Printf("Invalid byr: %q\n", p.Byr)
			return false
		}
	} else {
		log.Printf("Invalid byr: %q\n", p.Byr)
		return false
	}

	if iyr, err := strconv.Atoi(p.Iyr); err == nil {
		if iyr < 2010 || iyr > 2020 {
			log.Printf("Invalid iyr: %q\n", p.Iyr)
			return false
		}
	} else {
		log.Printf("Invalid iyr: %q\n", p.Iyr)
		return false
	}

	if eyr, err := strconv.Atoi(p.Eyr); err == nil {
		if eyr < 2020 || eyr > 2030 {
			log.Printf("Invalid eyr: %q\n", p.Eyr)
			return false
		}
	} else {
		log.Printf("Invalid eyr: %q\n", p.Eyr)
		return false
	}

	var hgtInt int
	if n, err := fmt.Sscanf(p.Hgt, "%dcm", &hgtInt); err == nil {
		if n != 1 {
			log.Println("hgt cm Sscanf parse failure")
			return false
		}

		if hgtInt < 150 || hgtInt > 193 {
			log.Printf("Invalid hgt cm: %q\n", p.Hgt)
			return false
		}
	} else if _, err := fmt.Sscanf(p.Hgt, "%din", &hgtInt); err == nil {
		if n != 1 {
			log.Println("hgt in Sscanf parse failure")
			return false
		}

		if hgtInt < 59 || hgtInt > 76 {
			log.Printf("Invalid hgt in: %q\n", p.Hgt)
			return false
		}
	} else {
		log.Printf("Invalid hgt: %q\n", p.Hgt)
		return false
	}

	if matched, err := regexp.MatchString("^#[0-9a-f]{6}$", p.Hcl); !matched || err != nil {
		log.Printf("Invalid hcl: %q\n", p.Hcl)
		return false
	}

	validEcl := map[string]bool{
		"amb": true,
		"blu": true,
		"brn": true,
		"gry": true,
		"grn": true,
		"hzl": true,
		"oth": true,
	}
	if _, ok := validEcl[p.Ecl]; !ok {
		log.Printf("Invalid ecl: %q\n", p.Ecl)
		return false
	}

	if matched, err := regexp.MatchString("^[0-9]{9}$", p.Pid); !matched || err != nil {
		log.Printf("Invalid pid: %q\n", p.Pid)
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
