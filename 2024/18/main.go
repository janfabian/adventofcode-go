package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"strconv"
	"strings"
)

type Memory struct {
	Space [][]bool
	End   lib.Coordinate
	Last  lib.Coordinate
}

func (memory *Memory) Print() {
	g := lib.EmptyGrid(len(memory.Space[0]), len(memory.Space))

	for Y, row := range memory.Space {
		for X, cell := range row {
			if cell {
				g[Y][X] = '#'
			}
		}
	}

	g[memory.End.Y][memory.End.X] = 'X'

	lib.PrintGridTable(g)
}

func (memory *Memory) FindPath() int {
	start := lib.Coordinate{X: 0, Y: 0}

	toVisit := []lib.Coordinate{start}
	visited := map[lib.Coordinate]int{start: 0}

	var current lib.Coordinate

	for len(toVisit) > 0 {
		current = toVisit[0]
		toVisit = toVisit[1:]

		if current == memory.End {
			break
		}

		for _, next := range current.Around() {
			if next.X < 0 || next.X > memory.End.X || next.Y < 0 || next.Y > memory.End.Y {
				continue
			}

			if memory.Space[next.Y][next.X] {
				continue
			}

			if l, ok := visited[*next]; ok {
				if l <= visited[current]+1 {
					continue
				}
			}

			toVisit = append(toVisit, *next)
			visited[*next] = visited[current] + 1
		}
	}

	if current != memory.End {
		return -1
	}

	return visited[current]
}

func Solve(memory *Memory, part2 bool) (int, error) {
	// memory.Print()

	return memory.FindPath(), nil
}

func ParseInput(input string, end lib.Coordinate, limit int, part2 bool) (*Memory, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	memory := &Memory{
		End: end,
	}

	memory.Space = make([][]bool, end.Y+1)
	for i := range memory.Space {
		memory.Space[i] = make([]bool, end.X+1)
	}

	i := 0
	for scanner.Scan() {
		if i >= limit {
			break
		}

		line := scanner.Text()

		coords := strings.Split(line, ",")
		X, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing x coordinate: %v", err)
		}

		Y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing y coordinate: %v", err)
		}

		memory.Space[Y][X] = true
		memory.Last = lib.Coordinate{X: X, Y: Y}
		i++
	}

	return memory, nil
}
