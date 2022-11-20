package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readInput() []string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return lines
}

type food struct {
	Ingredients map[string]string
}
type foods []*food

func (f *food) String() string {
	str := ""
	for i := range f.Ingredients {
		str += " " + i
	}
	return str[1:]
}

type allergen struct {
	Name       string
	Foods      foods
	Candidates map[string]string
}
type allergens map[string]*allergen

func (a *allergen) String() string {
	return fmt.Sprintf("[Foods: %q, Candidates: %q]", a.Foods, a.Candidates)
}

func parseAllergens(input []string) (foods, allergens) {
	af := make(foods, 0)
	as := make(allergens)

	for _, line := range input {
		ingredients := make(map[string]string)

		tokens := strings.Split(line, " ")
		var i int
		var t string
		for i, t = range tokens {
			if t[0] == '(' {
				break
			}

			ingredients[t] = ""
		}
		f := &food{ingredients}
		af = append(af, f)

		for _, t = range tokens[i+1:] { // +1 to skip "(contains"
			allergenName := strings.TrimSuffix(t, ")")
			allergenName = strings.TrimSuffix(allergenName, ",")

			a, ok := as[allergenName]
			if ok {
				a.Foods = append(a.Foods, f)
			} else {
				as[allergenName] = &allergen{Name: allergenName, Foods: []*food{f}}
			}
		}
	}

	return af, as
}

func solveAllergens(allFoods foods, as allergens) map[string]string {
	allergenMap := make(map[string]string)

	unsolved := make(map[string]string)
	for name, a := range as {
		unsolved[name] = ""
		candidates := make(map[string]string)

		for i := range a.Foods[0].Ingredients {
			allPresent := true
			for _, f := range a.Foods[1:] {
				if _, ok := f.Ingredients[i]; !ok {
					allPresent = false
					break
				}
			}

			if allPresent {
				candidates[i] = ""
			}
		}

		as[name].Candidates = candidates
	}

	for len(allergenMap) < len(as) {
		newlySolved := make(map[string]string)
		for name := range unsolved {
			a, _ := as[name]
			if len(a.Candidates) != 1 {
				continue
			}

			var soln string
			for i := range a.Candidates {
				soln = i
				allergenMap[soln] = name

			}
			newlySolved[name] = ""

			for _, a := range as {
				if _, ok := a.Candidates[soln]; ok {
					delete(a.Candidates, soln)
				}
			}
		}

		if len(newlySolved) == 0 {
			panic("Stalled")
		}
	}

	return allergenMap
}

func countClean(af foods, allergenMap map[string]string) int {
	var count int
	for _, f := range af {
		for i := range f.Ingredients {
			if _, ok := allergenMap[i]; !ok {
				count++
			}
		}
	}

	return count
}

func main() {
	input := readInput()
	af, as := parseAllergens(input)
	am := solveAllergens(af, as)

	fmt.Println(countClean(af, am))
}
