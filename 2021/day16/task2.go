package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readBinary(s string) int {
	value := 0

	for _, r := range s {
		value *= 2
		if r == '1' {
			value += 1
		}
	}

	return value
}

type packet interface {
	version() int
	typeId() int
	value() int // I bet this will be useful later for operators...
	versionSum() int
}

type literal struct {
	ver int
	val int
}

func (l *literal) version() int    { return l.ver }
func (l *literal) typeId() int     { return 4 }
func (l *literal) value() int      { return l.val }
func (l *literal) versionSum() int { return l.ver }
func (l *literal) String() string {
	return fmt.Sprintf("{%d %d %d}", l.version(), l.typeId(), l.value())
}

func parseLiteralPayload(version int, s string) (*literal, int) {
	parsedChars := 0
	groupIndex := 0
	value := 0
	for {
		parsedChars += 5
		value *= 16

		group := s[groupIndex+1 : groupIndex+5]
		value += readBinary(group)

		prefix := s[groupIndex]
		if prefix == '0' {
			break
		}

		groupIndex += 5
	}

	return &literal{version, value}, parsedChars
}

type operator struct {
	ver          int
	typ          int
	lengthTypeId int
	subPackets   []packet
}

func (o *operator) version() int { return o.ver }
func (o *operator) typeId() int  { return o.typ }

func (o *operator) value() int {
	switch o.typ {
	case 0:
		sum := 0
		for _, p := range o.subPackets {
			sum += p.value()
		}
		return sum
	case 1:
		product := 1
		for _, p := range o.subPackets {
			product *= p.value()
		}
		return product
	case 2:
		min := 1000000000000000
		for _, p := range o.subPackets {
			val := p.value()
			if min > val {
				min = val
			}
		}
		return min
	case 3:
		max := -1
		for _, p := range o.subPackets {
			val := p.value()
			if max < val {
				max = val
			}
		}
		return max
	case 5:
		if o.subPackets[0].value() > o.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	case 6:
		if o.subPackets[0].value() < o.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	case 7:
		if o.subPackets[0].value() == o.subPackets[1].value() {
			return 1
		} else {
			return 0
		}
	default:
		panic("Unknown type ID")
	}
}

func (o *operator) versionSum() int {
	sum := o.ver
	for _, p := range o.subPackets {
		sum += p.versionSum()
	}
	return sum
}
func (o *operator) String() string {
	return fmt.Sprintf("{%d %d %d %s}", o.version(), o.typeId(), o.lengthTypeId, o.subPackets)
}

func parsePacketsByBitCount(bitCount int, s string) []packet {
	parsedChars := 0
	packets := make([]packet, 0)

	for parsedChars < bitCount {
		packet, packetLength := parsePacket(s[parsedChars:])
		packets = append(packets, packet)
		parsedChars += packetLength
	}

	return packets
}

func parsePacketsByPacketCount(packetCount int, s string) ([]packet, int) {
	parsedChars := 0
	packets := make([]packet, packetCount)
	var packetLength int

	for i := range packets {
		packets[i], packetLength = parsePacket(s[parsedChars:])
		parsedChars += packetLength
	}

	return packets, parsedChars
}

func parseOperatorPayload(version, typeId int, s string) (*operator, int) {
	lengthTypeId := int(s[0] - '0')

	var parsedChars int
	var subPackets []packet

	switch lengthTypeId {
	case 0:
		bitCount := readBinary(s[1:16])
		subPackets = parsePacketsByBitCount(bitCount, s[16:])
		parsedChars = 16 + bitCount
	case 1:
		packetCount := readBinary(s[1:12])

		subPackets, parsedChars = parsePacketsByPacketCount(packetCount, s[12:])
		parsedChars += 12
	}

	o := operator{
		ver:          version,
		typ:          typeId,
		lengthTypeId: lengthTypeId,
		subPackets:   subPackets,
	}

	return &o, parsedChars
}

func parsePacket(s string) (packet, int) {
	version, typeId := readBinary(s[0:3]), readBinary(s[3:6])

	var m packet
	var parsedChars int
	switch typeId {
	case 4:
		m, parsedChars = parseLiteralPayload(version, s[6:])
	default:
		m, parsedChars = parseOperatorPayload(version, typeId, s[6:])
	}

	parsedChars += 6 // version and typeId strings

	return m, parsedChars
}

// When you write everything assuming binary string input but forgot it's actually hexadecimal...
func hexadecimalToBinaryString(s string) string {
	var b strings.Builder

	for _, r := range s {
		val := int(r - '0')
		if val > 15 {
			val = 10 + int(r-'A')
		}

		b.WriteString(fmt.Sprintf("%.4b", val))
	}

	return b.String()
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

func parseInput(input []string) []packet {
	packets := make([]packet, len(input))
	for i, line := range input {
		binaryString := hexadecimalToBinaryString(line)
		packets[i], _ = parsePacket(binaryString)
	}

	return packets
}

func main() {
	input := readInput()
	packets := parseInput(input)

	for _, p := range packets {
		fmt.Println(p.value())
	}

}
