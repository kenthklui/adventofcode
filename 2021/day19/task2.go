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

func NewVec3(x, y, z int) *vec3      { return &vec3{x, y, z, nil} }
func (v *vec3) String() string       { return fmt.Sprintf("[%d,%d,%d]", v.x, v.y, v.z) }
func (v *vec3) equals(u *vec3) bool  { return v.x == u.x && v.y == u.y && v.z == u.z }
func (v *vec3) add(u *vec3) *vec3    { return NewVec3(v.x+u.x, v.y+u.y, v.z+u.z) }
func (v *vec3) sub(u *vec3) *vec3    { return NewVec3(v.x-u.x, v.y-u.y, v.z-u.z) }
func (v *vec3) rotate(rot int) *vec3 { return v.rotations()[rot] }

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

func manhattan(u, v *vec3) int {
	w := u.sub(v)
	m := 0
	if w.x > 0 {
		m += w.x
	} else {
		m -= w.x
	}
	if w.y > 0 {
		m += w.y
	} else {
		m -= w.y
	}
	if w.z > 0 {
		m += w.z
	} else {
		m -= w.z
	}

	return m
}

// For sorting
func (v3s vec3s) Len() int           { return len(v3s) }
func (v3s vec3s) Swap(i, j int)      { v3s[i], v3s[j] = v3s[j], v3s[i] }
func (v3s vec3s) Less(i, j int) bool { return vec3Less(v3s[i], v3s[j]) }

type scanner struct {
	beacons    vec3s
	normalized vec3s

	index    int
	position *vec3
	rotation int
}

func NewScanner(index int) *scanner { return &scanner{index: index, beacons: make(vec3s, 0)} }
func (s *scanner) String() string {
	if s.position == nil {
		return "[]"
	} else {
		return fmt.Sprintf("[%s : Rotation %d]", s.position, s.rotation)
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
	return rotated
}

func (s *scanner) normalizedBeacons() vec3s {
	if s.position == nil {
		panic("Cannot produce normalized beacons yet")
	}

	if s.normalized == nil {
		s.normalized = make(vec3s, 0, len(s.beacons))
		for _, b := range s.beacons {
			bTransform := b.rotate(s.rotation).add(s.position)
			s.normalized = append(s.normalized, bTransform)
		}
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
	dist := s.position.sub(beacon)
	threshold := 1000
	threshSq := threshold * threshold

	if dist.x*dist.x > threshSq || dist.y*dist.y > threshSq || dist.z*dist.z > threshSq {
		return false
	}

	return true
}

func matchScanners(s1, s2 *scanner) bool {
	if s1.position == nil {
		panic(fmt.Errorf("Matching unsolved %d against %d", s1.index, s2.index))
	} else if s2.position != nil {
		panic(fmt.Errorf("Matching %d against solved %d", s1.index, s2.index))
	}

	matchThreshold := 12

	for rotation := 0; rotation < 24; rotation++ {
		s2rotated := s2.rotatedBeacons(rotation)

		for _, b1 := range s1.normalizedBeacons() {
			for j, b2 := range s2rotated {
				s2pos := b1.sub(b2)

				matchedBeacons := 1
				for k, tb := range s2rotated {
					if j == k {
						continue
					}

					remainingBeacons := len(s2rotated) - k
					if remainingBeacons+matchedBeacons < matchThreshold {
						// Insufficient remaining beacons to reach threshold; failed match
						break
					}

					tb = tb.add(s2pos)
					if s1.checkInRange(tb) {
						if s1.findBeacon(tb) != -1 {
							matchedBeacons++
						} else {
							// Failure to detect beacon in range implies failed match
							break
						}
					}
				}

				if matchedBeacons >= matchThreshold {
					s2.position = s2pos
					s2.rotation = rotation
					return true
				}
			}
		}
	}

	return false
}

func mapScanners(scanners []*scanner) {
	solvedCount := 1

	for _, solver := range scanners {
		success := 0

		unsolved := scanners[solvedCount:]
		for j, toSolve := range unsolved {
			if matchScanners(solver, toSolve) {
				// Swap solved scanners to front
				unsolved[j], unsolved[success] = unsolved[success], unsolved[j]
				success++
			}
		}

		solvedCount += success
	}
}

func listBeacons(scanners []*scanner) vec3s {
	deduped := make(vec3s, 0)
	for _, s := range scanners {
		deduped = append(deduped, s.normalizedBeacons()...)
	}
	sort.Sort(deduped)

	count := 1
	for _, v := range deduped[1:] {
		if deduped[count-1].equals(v) {
			continue
		} else {
			deduped[count] = v
			count++
		}
	}

	return deduped[:count]
}

func maxManhattan(scanners []*scanner) int {
	max := 0
	for i, s1 := range scanners {
		for j, s2 := range scanners {
			if i == j {
				continue
			}

			m := manhattan(s1.position, s2.position)
			if m > max {
				max = m
			}
		}
	}

	return max
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
			s = NewScanner(len(scanners))
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
	scanners[0].position = NewVec3(0, 0, 0)
	scanners[0].rotation = 11

	return scanners
}

func main() {
	input := readInput()
	scanners := parseInput(input)

	mapScanners(scanners)
	fmt.Println(maxManhattan(scanners))
}
