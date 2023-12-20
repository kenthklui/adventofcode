package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/kenthklui/adventofcode/util"
)

const LOW, HIGH = false, true
const broadcasterName = "broadcaster"

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
	handle(p pulse, buttonPresses int) []pulse
	addSource(source module)
	addTarget(target module)
	upstream() moduleList
	stateString() string
}

type moduleList []module

func (ml moduleList) Len() int           { return len(ml) }
func (ml moduleList) Less(i, j int) bool { return ml[i].name() < ml[j].name() }
func (ml moduleList) Swap(i, j int)      { ml[i], ml[j] = ml[j], ml[i] }

func (ml moduleList) stateString() string {
	var b strings.Builder
	for _, m := range ml {
		b.WriteString(m.stateString())
	}
	return b.String()
}

type void struct{}

var nul void

func sourceList(sources []module) moduleList {
	dedup := make(map[module]void)
	queue := slices.Clone(sources)
	for len(queue) > 0 {
		q := queue[0]
		queue = queue[1:]

		if q.name() == broadcasterName {
			continue
		}
		if _, ok := dedup[q]; ok {
			continue
		}

		dedup[q] = nul
		queue = append(queue, q.upstream()...)
	}
	ml := make(moduleList, 0, len(dedup))
	for d := range dedup {
		ml = append(ml, d)
	}
	sort.Sort(ml)
	ml = slices.CompactFunc(ml, func(m1, m2 module) bool { return m1.name() == m2.name() })
	return ml
}

type flipFlop struct {
	label            string
	state            bool
	sources, targets []module
}

func (ff *flipFlop) name() string { return ff.label }
func (ff *flipFlop) handle(p pulse, buttonPresses int) []pulse {
	if p.signal {
		return nil
	} else {
		ff.state = !ff.state
		return makePulses(ff, ff.targets, ff.state)
	}
}
func (ff *flipFlop) addSource(source module) { ff.sources = append(ff.sources, source) }
func (ff *flipFlop) addTarget(target module) { ff.targets = append(ff.targets, target) }
func (ff *flipFlop) upstream() moduleList    { return ff.sources }
func (ff *flipFlop) stateString() string {
	if ff.state {
		return "1"
	} else {
		return "0"
	}
}

type conjunction struct {
	label            string
	recents          map[module]bool
	sources, targets []module
}

func (c *conjunction) name() string { return c.label }
func (c *conjunction) handle(p pulse, buttonPresses int) []pulse {
	c.recents[p.source] = p.signal
	for _, h := range c.recents {
		if !h {
			return makePulses(c, c.targets, HIGH)
		}
	}
	return makePulses(c, c.targets, LOW)
}
func (c *conjunction) addSource(source module) {
	c.sources = append(c.sources, source)
	c.recents[source] = LOW
}
func (c *conjunction) addTarget(target module) { c.targets = append(c.targets, target) }
func (c *conjunction) upstream() moduleList    { return c.sources }
func (c *conjunction) stateString() string {
	byteLen := (len(c.sources)-1)/8 + 1
	bytes := make([]byte, byteLen)
	for i, s := range c.sources {
		if c.recents[s] {
			bytes[i/8] |= 1 << (i % 8)
		}
	}
	return string(bytes)
}

type broadcast struct {
	targets []module
}

func (bc *broadcast) name() string { return broadcasterName }
func (bc *broadcast) handle(p pulse, buttonPresses int) []pulse {
	return makePulses(bc, bc.targets, p.signal)
}
func (bc *broadcast) addSource(source module) {}
func (bc *broadcast) addTarget(target module) { bc.targets = append(bc.targets, target) }
func (bc *broadcast) upstream() moduleList    { return make(moduleList, 0) }
func (bc *broadcast) stateString() string     { return "" }

type output struct {
	label   string
	sources []module
	lows    int
}

func (o *output) name() string { return o.label }
func (o *output) handle(p pulse, buttonPresses int) []pulse {
	if p.signal == LOW {
		o.lows++
	}
	return nil
}
func (o *output) addSource(source module) {
	o.sources = append(o.sources, source)
}
func (o *output) addTarget(target module) {}
func (o *output) upstream() moduleList    { return o.sources }
func (o *output) stateString() string     { return "" }

func (o *output) reset() { o.lows = 0 }

type machine struct {
	modules       map[string]module
	pulseQueue    []pulse
	out           *output
	buttonPresses int
}

func (m *machine) button() {
	m.buttonPresses++
	if bc, ok := m.modules[broadcasterName]; ok {
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

		pulses := p.target.handle(p, m.buttonPresses)
		if pulses != nil {
			for _, np := range pulses {
				m.pulseQueue = append(m.pulseQueue, np)
			}
		}
	}
}

func parseMachine(input []string) *machine {
	m := machine{
		modules:       make(map[string]module),
		pulseQueue:    make([]pulse, 0, len(input)*2),
		buttonPresses: 0,
	}

	targetStrs := make(map[string][]string)
	for _, line := range input {
		name, targets, _ := strings.Cut(line, " -> ")
		targetNames := strings.Split(targets, ", ")
		emptyTargets := make([]module, 0, len(targetNames))
		switch name[0] {
		case 'b': // broadcaster
			bc := broadcast{
				targets: emptyTargets,
			}
			m.modules[name] = &bc
		case '%':
			name = name[1:]
			ff := flipFlop{
				label:   name,
				state:   LOW,
				targets: emptyTargets,
			}
			m.modules[name] = &ff
		case '&':
			name = name[1:]
			c := conjunction{
				label:   name,
				recents: make(map[module]bool),
				sources: make([]module, 0),
				targets: emptyTargets,
			}
			m.modules[name] = &c
		}
		targetStrs[name] = targetNames
	}

	for name, targetNames := range targetStrs {
		source := m.modules[name]
		for _, targetName := range targetNames {
			target, ok := m.modules[targetName]
			if !ok {
				if m.out != nil {
					panic("Already have an output module")
				}
				m.out = &output{
					label:   targetName,
					sources: make([]module, 0),
				}
				target = m.out
				m.modules[targetName] = target
			}
			target.addSource(source)
			source.addTarget(target)
		}
	}

	return &m
}

func (m *machine) detectCycles() []int {
	cycleGroups := make([]moduleList, 0)
	for _, s1 := range m.out.sources {
		c1 := s1.(*conjunction)
		for _, s2 := range c1.sources {
			c2 := s2.(*conjunction)
			cycleGroups = append(cycleGroups, sourceList(c2.sources))
		}
	}

	cycles := make([]int, len(cycleGroups))
	prevStateStrs := make([]string, len(cycles))
	stateStrMaps := make([]map[string]string, len(cycles))
	for i, cg := range cycleGroups {
		stateStrMaps[i] = make(map[string]string)
		prevStateStrs[i] = cg.stateString()
	}

	cyclesDetected := 0
	for cyclesDetected < len(cycles) {
		m.button()
		for i, cg := range cycleGroups {
			if cycles[i] != 0 {
				continue
			}
			ss := cg.stateString()
			stateStrMaps[i][prevStateStrs[i]] = ss
			if _, ok := stateStrMaps[i][ss]; ok {
				cycleLength := 1
				for str := ss; stateStrMaps[i][str] != ss; str = stateStrMaps[i][str] {
					cycleLength++
				}
				cycles[i] = cycleLength
				cyclesDetected++
			} else {
				prevStateStrs[i] = ss
			}
		}
	}
	sort.Ints(cycles)
	return cycles
}

func main() {
	input := util.StdinReadlines()
	m := parseMachine(input)
	product := 1
	for _, c := range m.detectCycles() {
		product *= c
	}
	fmt.Println(product)
}
