#!/usr/bin/env ruby

Workflow = Struct.new(:name, :rules, :dest)
Rule = Struct.new(:index, :op, :threshold, :dest)

def apply(rule, input)
  divert = input.map(&:dup)
  if rule.op == '<'
    if input[rule.index].last < rule.threshold # Everything matches
      input = nil
    elsif input[rule.index].first >= rule.threshold # Nothing matches
      divert = nil
    else # Split
      input[rule.index][0] = rule.threshold
      divert[rule.index][1] = rule.threshold - 1
    end
  elsif rule.op == '>'
    if input[rule.index].first > rule.threshold # Everything matches
      input = nil
    elsif input[rule.index].last <= rule.threshold # Nothing matches
      divert = nil
    else # Split
      input[rule.index][1] = rule.threshold
      divert[rule.index][0] = rule.threshold + 1
    end
  end

  [input, divert]
end

def process(workflows, name, input)
  accepted = 0
  workflows[name].rules.each do |rule|
    break if input.nil?
    input, divert = apply(rule, input)
    accepted += countAccepted(workflows, rule.dest, divert)
  end
  accepted += countAccepted(workflows, workflows[name].dest, input)

  accepted
end

def countAccepted(workflows, dest, input)
  return 0 if input.nil?
  case dest
  when 'A'
    input.reduce(1) { |p, x| p *= (x[1] - x[0] + 1) } # Product of range sizes
  when 'R'
    0
  else
    process(workflows, dest, input)
  end
end

def parseWorkflow(line)
  lineMatch = /(\w+){(.+)}/.match(line)
  rules = lineMatch[2].split(',').map do |str|
    if ruleMatch = /(\w)([<>])(\d+):(\w+)/.match(str)
      index = 'xmas'.index(ruleMatch[1])
      Rule.new(index, ruleMatch[2], ruleMatch[3].to_i, ruleMatch[4])
    else
      str
    end
  end
  dest = rules.pop()
  Workflow.new(lineMatch[1], rules, dest)
end

def parse()
  lines = STDIN.readlines.map(&:chomp)
  lines = lines.first(lines.index("")) # Remove everything after the break
  flows = lines.map { |l| parseWorkflow(l) }
  Hash[flows.map { |f| [f.name, f] }]
end

workflows = parse()
input = [[1, 4000], [1, 4000], [1, 4000], [1, 4000]]
puts countAccepted(workflows, "in", input)
