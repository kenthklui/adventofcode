package main

import (
	"fmt"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

const LOW, HIGH = false, true

type pulse struct {
	source, target module
	signal         bool
}

func (p pulse) String() string {
	signalStr := "low"
	if p.signal {
		signalStr = "high"
	}
	if p.source == nil {
		return fmt.Sprintf("button -%s-> %s", signalStr, p.target.name())
	} else {
		return fmt.Sprintf("%s -%s-> %s", p.source.name(), signalStr, p.target.name())
	}
}

func makePulses(source module, targets []module, signal bool) []pulse {
	pulses := make([]pulse, 0, len(targets))
	for _, t := range targets {
		pulses = append(pulses, pulse{source, t, signal})
	}
	return pulses
}

type module interface {
	name() string
	handle(p pulse) []pulse
	addSource(source module)
	addTarget(target module)
}

type flipFlop struct {
	label   string
	state   bool
	targets []module
}

func (ff *flipFlop) name() string { return ff.label }
func (ff *flipFlop) handle(p pulse) []pulse {
	if p.signal {
		return nil
	} else {
		ff.state = !ff.state
		return makePulses(ff, ff.targets, ff.state)
	}
}
func (ff *flipFlop) addSource(source module) {}
func (ff *flipFlop) addTarget(target module) { ff.targets = append(ff.targets, target) }

type conjunction struct {
	label   string
	recents map[module]bool
	targets []module
}

func (c *conjunction) name() string { return c.label }
func (c *conjunction) handle(p pulse) []pulse {
	c.recents[p.source] = p.signal
	for _, h := range c.recents {
		if !h {
			return makePulses(c, c.targets, HIGH)
		}
	}
	return makePulses(c, c.targets, LOW)
}
func (c *conjunction) addSource(source module) { c.recents[source] = LOW }
func (c *conjunction) addTarget(target module) { c.targets = append(c.targets, target) }

type broadcast struct {
	targets []module
}

func (bc *broadcast) name() string            { return "broadcaster" }
func (bc *broadcast) handle(p pulse) []pulse  { return makePulses(bc, bc.targets, p.signal) }
func (bc *broadcast) addSource(source module) {}
func (bc *broadcast) addTarget(target module) { bc.targets = append(bc.targets, target) }

type output struct{}

func (o *output) name() string            { return "output" }
func (o *output) handle(p pulse) []pulse  { return nil }
func (o *output) addSource(source module) {}
func (o *output) addTarget(target module) {}

type machine struct {
	modules     map[string]module
	pulseQueue  []pulse
	lows, highs int
}

func (m *machine) button() {
	if bc, ok := m.modules["broadcaster"]; ok {
		buttonPulse := pulse{nil, bc, LOW}
		m.pulseQueue = append(m.pulseQueue, buttonPulse)
	} else {
		panic("Broadcaster module not found")
	}
	m.start()
}

func (m *machine) start() {
	for len(m.pulseQueue) > 0 {
		p := m.pulseQueue[0]
		m.pulseQueue = m.pulseQueue[1:]
		// fmt.Println(p)

		if p.signal {
			m.highs++
		} else {
			m.lows++
		}

		pulses := p.target.handle(p)
		if pulses != nil {
			for _, np := range pulses {
				m.pulseQueue = append(m.pulseQueue, np)
			}
		}
	}
}

func parseMachine(input []string) *machine {
	m := machine{
		modules:    make(map[string]module),
		pulseQueue: make([]pulse, 0, len(input)*2),
		lows:       0,
		highs:      0,
	}

	targetStrs := make(map[string][]string)
	for _, line := range input {
		name, targets, _ := strings.Cut(line, " -> ")
		targetNames := strings.Split(targets, ", ")
		emptyTargets := make([]module, 0, len(targetNames))
		switch name[0] {
		case 'b': // broadcaster
			bc := broadcast{emptyTargets}
			m.modules[name] = &bc
		case '%':
			name = name[1:]
			ff := flipFlop{name, LOW, emptyTargets}
			m.modules[name] = &ff
		case '&':
			name = name[1:]
			c := conjunction{name, make(map[module]bool), emptyTargets}
			m.modules[name] = &c
		}
		targetStrs[name] = targetNames
	}

	for name, targetNames := range targetStrs {
		source := m.modules[name]
		for _, targetName := range targetNames {
			target, ok := m.modules[targetName]
			if !ok {
				target = &output{}
				m.modules[targetName] = target // Why are modules missing?
			}
			target.addSource(source)
			source.addTarget(target)
		}
	}

	return &m
}

const buttonPresses = 1000

func main() {
	input := util.StdinReadlines()
	m := parseMachine(input)
	for i := 0; i < buttonPresses; i++ {
		m.button()
	}
	fmt.Println(m.lows * m.highs)
}
