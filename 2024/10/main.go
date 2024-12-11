package aoc2024_10

import (
	"adventofcode/lib"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Coordinate struct {
	X, Y int
}

type NextTrailStop struct {
	coord      Coordinate
	prevHeight int
}

type Trailhead struct {
	start   Coordinate
	visited map[Coordinate]struct{}
	next    []NextTrailStop
	input   [][]int
	result  int
	part2   bool
}

func FindTrailheads(input [][]int, part2 bool) []*Trailhead {
	trailHeads := []*Trailhead{}

	size := len(input) * len(input[0])

	for Y, line := range input {
		for X, val := range line {
			if val == 0 {
				trailHeads = append(trailHeads, &Trailhead{
					start:   Coordinate{X, Y},
					visited: make(map[Coordinate]struct{}, size),
					next:    make([]NextTrailStop, 0, size),
					input:   input,
					part2:   part2,
				})
			}
		}
	}

	return trailHeads

}

func (t *Trailhead) Search() {
	t.next = append(t.next, NextTrailStop{coord: t.start, prevHeight: -1})

	for {
		if len(t.next) == 0 {
			break
		}

		t.Traverse(t.next[0])
		t.next = t.next[1:]
	}
}

func (t *Trailhead) Traverse(trailStop NextTrailStop) {
	coord := trailStop.coord

	if !t.part2 {
		if _, exists := t.visited[coord]; exists {
			return
		}
	}

	height := t.input[coord.Y][coord.X]
	if height != trailStop.prevHeight+1 {
		return
	}

	if !t.part2 {
		t.visited[coord] = struct{}{}
	}

	if height == 9 {
		t.result++
		return
	}

	if coord.X > 0 {
		t.next = append(t.next, NextTrailStop{coord: Coordinate{coord.X - 1, coord.Y}, prevHeight: height})
	}

	if coord.X < len(t.input[0])-1 {
		t.next = append(t.next, NextTrailStop{coord: Coordinate{coord.X + 1, coord.Y}, prevHeight: height})
	}

	if coord.Y > 0 {
		t.next = append(t.next, NextTrailStop{coord: Coordinate{coord.X, coord.Y - 1}, prevHeight: height})
	}

	if coord.Y < len(t.input)-1 {
		t.next = append(t.next, NextTrailStop{coord: Coordinate{coord.X, coord.Y + 1}, prevHeight: height})
	}

}

func Solve(input [][]int, part2 bool) int {
	trailHeads := FindTrailheads(input, part2)

	var wg sync.WaitGroup
	poolSize := runtime.NumCPU() * 100
	pool := make(chan struct{}, poolSize)

	for _, trailHead := range trailHeads {
		wg.Add(1)
		pool <- struct{}{}

		go func(trailHead *Trailhead) {
			defer wg.Done()
			defer func() { <-pool }()

			trailHead.Search()
		}(trailHead)
	}

	wg.Wait()
	close(pool)

	result := 0
	for _, trailHead := range trailHeads {
		result += trailHead.result
	}

	return result
}

func ParseInput(input string) ([][]int, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "")

		lineInt := []int{}

		for _, part := range parts {
			intval, err := strconv.Atoi(part)

			if err != nil {
				return nil, fmt.Errorf("error converting to int: %v", err)
			}

			lineInt = append(lineInt, intval)
		}

		result = append(result, lineInt)
	}

	return result, nil
}

func ParseOutput(input string) (int, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return 0, fmt.Errorf("error reading file: %v", err)
	}

	scanner.Scan()
	line := scanner.Text()
	intval, err := strconv.Atoi(line)

	if err != nil {
		return 0, fmt.Errorf("error converting to int: %v", err)
	}

	return intval, nil

}
