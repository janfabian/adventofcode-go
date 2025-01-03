package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"strings"
	"sync"
)

type Onsen struct {
	PatternMap map[string]struct{}
	MaxN       int
	Towels     []string
}

func (o *Onsen) Patterns(towel string) int {
	counter := make([]int, len(towel))

	for i := 0; i < len(towel); i++ {
		if i != 0 && counter[i-1] == 0 {
			continue
		}

		for j := 0; j < min(o.MaxN, len(towel)-i); j++ {
			if _, ok := o.PatternMap[towel[i:i+j+1]]; ok {

				if i == 0 {
					counter[i+j]++
				} else {
					counter[i+j] += counter[i-1]
				}
			}
		}
	}

	return counter[len(towel)-1]
}

func Solve(onsen *Onsen, part2 bool) int {

	var wg sync.WaitGroup
	wg.Add(len(onsen.Towels))

	i := 0
	p := 0
	var mu sync.Mutex

	for _, towel := range onsen.Towels {
		lib.AddTask(func() {
			defer wg.Done()

			patterns := onsen.Patterns(towel)

			if patterns > 0 {
				mu.Lock()
				i++
				p += patterns
				mu.Unlock()
			}

		})
	}

	wg.Wait()

	if part2 {
		return p
	} else {
		return i
	}
}

func ParseInput(input string, part2 bool) (*Onsen, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	onsen := &Onsen{
		PatternMap: make(map[string]struct{}),
	}

	scanner.Scan()
	line := scanner.Text()
	scanner.Scan()

	maxN := 0
	for _, pattern := range strings.Split(line, ", ") {
		onsen.PatternMap[pattern] = struct{}{}
		if len(pattern) > maxN {
			maxN = len(pattern)
		}
	}

	onsen.MaxN = maxN

	for scanner.Scan() {
		line := scanner.Text()

		onsen.Towels = append(onsen.Towels, line)
	}

	return onsen, nil
}
