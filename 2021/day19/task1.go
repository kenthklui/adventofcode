package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Let's go by right hand rule
type vec3 struct {
	x, y, z int

	// for memoizing rotations
	orientations vec3s
}
type vec3s []*vec3

func NewVec3(x, y, z int) *vec3                  { return &vec3{x, y, z, nil} }
func (v *vec3) String() string                   { return fmt.Sprintf("[%d,%d,%d]", v.x, v.y, v.z) }
func (v *vec3) equals(u *vec3) bool              { return v.x == u.x && v.y == u.y && v.z == u.z }
func (v *vec3) negate() *vec3                    { return NewVec3(-v.x, -v.y, -v.z) }
func (v *vec3) translate(u *vec3) *vec3          { return NewVec3(v.x+u.x, v.y+u.y, v.z+u.z) }
func (v *vec3) rotate(rot int) *vec3             { return v.rotations()[rot] }
func (v *vec3) transform(rot int, t *vec3) *vec3 { return v.rotate(rot).translate(t) }

// https://stackoverflow.com/a/16467849
func (v *vec3) turn() *vec3 { return NewVec3(-v.y, v.x, v.z) }
func (v *vec3) roll() *vec3 { return NewVec3(-v.z, v.y, v.x) }
func (v *vec3) rotations() vec3s {
	if v.orientations == nil {
		v.orientations = make(vec3s, 0, 24)

		vec := v
		for cycle := 0; cycle < 2; cycle++ {
			for rolls := 0; rolls < 3; rolls++ {
				vec = vec.roll()
				v.orientations = append(v.orientations, vec)
				for turns := 0; turns < 3; turns++ {
					vec = vec.turn()
					v.orientations = append(v.orientations, vec)
				}
			}
			vec = vec.roll().turn().roll()
		}
	}

	return v.orientations
}

func vec3Less(u, v *vec3) bool {
	if u.x == v.x {
		if u.y == v.y {
			return u.z < v.z
		} else {
			return u.y < v.y
		}
	} else {
		return u.x < v.x
	}
}

// For sorting
func (v3s vec3s) Len() int           { return len(v3s) }
func (v3s vec3s) Swap(i, j int)      { v3s[i], v3s[j] = v3s[j], v3s[i] }
func (v3s vec3s) Less(i, j int) bool { return vec3Less(v3s[i], v3s[j]) }

type scanner struct {
	beacons vec3s

	translation *vec3
	rotation    int
	normalized  vec3s
}

func NewScanner() *scanner { return &scanner{beacons: make(vec3s, 0)} }
func (s *scanner) String() string {
	if s.translation == nil {
		return "[]"
	} else {
		return fmt.Sprintf("[%s : Rotation %d]", s.translation, s.rotation)
	}
}
func (s *scanner) addBeacon(x, y, z int) {
	s.beacons = append(s.beacons, NewVec3(x, y, z))
}

func (s *scanner) rotatedBeacons(rot int) vec3s {
	rotated := make(vec3s, 0, len(s.beacons))
	for _, b := range s.beacons {
		rotated = append(rotated, b.rotate(rot))
	}
	// fmt.Printf("\nRotation %d:\n%s\nto\n%s\n", rot, s.beacons, rotated)
	return rotated
}

func (s *scanner) transformedBeacons(rotation int, translation *vec3) vec3s {
	transformed := make(vec3s, 0, len(s.beacons))
	for _, b := range s.beacons {
		transformed = append(transformed, b.transform(rotation, translation))
	}
	return transformed
}

func (s *scanner) normalizedBeacons() vec3s {
	if s.translation == nil {
		panic("Cannot produce normalized beacons yet")
	}

	if s.normalized == nil {
		s.normalized = s.transformedBeacons(s.rotation, s.translation)
	}

	return s.normalized
}

func (s *scanner) findBeacon(tar *vec3) int {
	for i, b := range s.normalizedBeacons() {
		if b.equals(tar) {
			return i
		}
	}

	return -1
}

func (s *scanner) checkInRange(beacon *vec3) bool {
	dist := s.translation.translate(beacon.negate())
	threshold := 1000
	threshSq := threshold * threshold

	if dist.x*dist.x > threshSq || dist.y*dist.y > threshSq || dist.z*dist.z > threshSq {
		return false
	}

	return true
}

func matchScanners(s1, s2 *scanner) bool {
	if s1.translation == nil || s2.translation != nil {
		return false
	}

	for rotation := 0; rotation < 24; rotation++ {
		rotated := s2.rotatedBeacons(rotation)

		for i := range s1.normalizedBeacons() {
			for j, b2 := range rotated {
				s2pos := locateScanner(s1, i, b2)

				matchedBeacons := 1
				for k, tb := range rotated {
					if j == k {
						continue
					}

					tb = tb.translate(s2pos)
					if s1.checkInRange(tb) && (s1.findBeacon(tb) != -1) {
						matchedBeacons++
					}
				}

				if matchedBeacons >= 12 {
					s2.translation = s2pos
					s2.rotation = rotation
					return true
				}
			}
		}
	}

	return false
}

func locateScanner(s1 *scanner, beaconIndex int, b2 *vec3) *vec3 {
	b1 := s1.normalizedBeacons()[beaconIndex]
	s2pos := b1.translate(b2.negate())

	return s2pos
}

func mapScanners(scanners []*scanner) {
	solvedScanners := 1

	for solvedScanners < len(scanners) {
		for i, s1 := range scanners {
			for j, s2 := range scanners {
				if i == j {
					continue
				}

				if matchScanners(s1, s2) {
					solvedScanners++
				}
			}
		}
	}
}

func listBeacons(scanners []*scanner) vec3s {
	all := make(vec3s, 0)
	for _, s := range scanners {
		all = append(all, s.normalizedBeacons()...)
	}
	sort.Sort(all)

	deduped := make(vec3s, 1)
	deduped[0] = all[0]
	for _, v := range all[1:] {
		if deduped[len(deduped)-1].equals(v) {
			continue
		} else {
			deduped = append(deduped, v)
		}
	}

	return deduped
}

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

func parseInput(input []string) []*scanner {
	scanners := make([]*scanner, 0)

	var s *scanner
	var x, y, z int

	for _, line := range input {
		if len(line) == 0 {
			scanners = append(scanners, s)
			s = nil
		} else if strings.Contains(line, "---") {
			s = NewScanner()
		} else {
			n, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
			if err != nil {
				panic(err)
			} else if n != 3 {
				panic("Failed to parse 3 coordinate entries")
			} else {
				s.addBeacon(x, y, z)
			}
		}
	}
	if s != nil {
		scanners = append(scanners, s)
	}

	// Set origin and default alignment using first scanner
	scanners[0].translation = NewVec3(0, 0, 0)
	scanners[0].rotation = 11

	return scanners
}

func main() {
	input := readInput()
	scanners := parseInput(input)

	mapScanners(scanners)
	beaconList := listBeacons(scanners)

	fmt.Println(len(beaconList))
}
