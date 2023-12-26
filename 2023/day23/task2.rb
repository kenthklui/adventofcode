#!/usr/bin/env ruby

Node = Struct.new(:location, :index, :adjacents)

class Area
  def initialize(lines)
    @start = [1, 0]
    @finish = [lines[0].size - 2, lines.size - 1]
    @width = lines[0].size
    @height = lines.size
    @lines = lines.map {|l| l.gsub(/[\^>v<]/, '.') }
  end

  def getNeighbors()
    @neighbors = { @start => [] }
    @lines.each_with_index do |row, y|
      next if y == 0
      row.split('').each_with_index do |c, x|
        next if x == 0
        next if c == '#'
        @neighbors[[x, y]] = []
        [[x-1, y], [x, y-1]].each do |new_x, new_y|
          next if @lines[new_y][new_x] == '#'
          @neighbors[[new_x, new_y]] << [x, y]
          @neighbors[[x, y]] << [new_x, new_y]
        end
      end
    end
    @neighbors
  end

  def buildGraph()
    startNode = Node.new(@start, 0, [])
    finishNode = Node.new(@finish, 1, [])
    @nodes = [startNode, finishNode]

    nodeFinder = { @start => startNode, @finish => finishNode }
    @neighbors.each do |source, dests|
      next if dests.size < 3
      newNode = Node.new(source, @nodes.size, [])
      @nodes << newNode
      nodeFinder[source] = newNode
    end

    @nodeDistances = []
    @nodes.size.times { |n| @nodeDistances << [0]*@nodes.size }
    @nodes.each_with_index do |node|
      @neighbors[node.location].each do |neighbor|
        nextFork, steps = findFork(node.location, neighbor, 1)
        next if nextFork.nil?
        next if nextFork == node.location

        forkNode = nodeFinder[nextFork]

        curValue = @nodeDistances[node.index][forkNode.index]
        if curValue == 0
          @nodeDistances[forkNode.index][node.index] = steps
          @nodeDistances[node.index][forkNode.index] = steps
          forkNode.adjacents << node
          node.adjacents << forkNode
        elsif steps > curValue
          @nodeDistances[forkNode.index][node.index] = steps
          @nodeDistances[node.index][forkNode.index] = steps
        end
      end
    end
  end

  def findFork(prev, cur, steps)
    if cur == @start || cur == @finish
      return cur, steps
    elsif @neighbors[cur].size > 2
      return cur, steps
    elsif @neighbors[cur].size == 2
      if @neighbors[cur][0] == prev
        return findFork(cur, @neighbors[cur][1], steps + 1)
      else
        return findFork(cur, @neighbors[cur][0], steps + 1)
      end
    else
      return nil, steps
    end
  end

  def longestHike()
    @visited = [false] * @nodes.size
    recurseHike(0, 1)
  end

  def recurseHike(fromIndex, toIndex)
    return 0 if fromIndex == toIndex

    @visited[fromIndex] = true
    distances = []
    unvisited = @nodes[fromIndex].adjacents.reject { |a| @visited[a.index] }
    unvisited.each do |ua|
      remain = recurseHike(ua.index, toIndex)
      if remain >= 0
        distances << @nodeDistances[fromIndex][ua.index] + remain
      end
    end
    @visited[fromIndex] = false

    return -1 if distances.empty?
    return distances.max
  end
end


def parse()
  lines = STDIN.readlines.map(&:chomp)
  return Area.new(lines)
end

area = parse()
area.getNeighbors()
area.buildGraph()
hike = area.longestHike()
puts hike
