#!/usr/bin/env python
from z3 import *
from parse import parse

import sys

def parseInput():
    lines = []
    for line in sys.stdin:
        result = parse('{}, {}, {} @ {}, {}, {}', line.strip())
        lines.append([int(s) for s in result])
    return lines

def solveLine(lines):
    pos_x, pos_y, pos_z = Int('x'), Int('y'), Int('z')
    dir_x, dir_y, dir_z = Int('dx'), Int('dy'), Int('dz')
    times = IntVector('t', 5)
    constraints = []
    s = Solver()
    s.add(And(dir_x < 1000, dir_x > -1000))
    s.add(And(dir_y < 1000, dir_y > -1000))
    s.add(And(dir_z < 1000, dir_z > -1000))
    for i, line in enumerate(lines[:5]):
        s.add(times[i] >= 0)
        s.add(And(pos_x + times[i] * dir_x == line[0] + times[i] * line[3]))
        s.add(And(pos_y + times[i] * dir_y == line[1] + times[i] * line[4]))
        s.add(And(pos_z + times[i] * dir_z == line[2] + times[i] * line[5]))
    s.check()
    m = s.model()
    data = [m[pos_x], m[pos_y], m[pos_z], m[dir_x], m[dir_y], m[dir_z]]
    return [i.as_long() for i in data]

if __name__ == "__main__":
    lines = parseInput()
    solution = solveLine(lines)
    # print(solution)
    print(sum(solution[:3]))
