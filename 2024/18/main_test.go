package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/18/"

type testCase struct {
	inputFile   string
	output      int
	description string
	end         lib.Coordinate
	limit       int
	part2       bool
	part2Result lib.Coordinate
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		{inputFile: "01.input", output: 22, description: "Puzzle case", end: lib.Coordinate{X: 6, Y: 6}, limit: 12},
		{inputFile: "puzzle.input", output: 304, description: "Puzzle case", end: lib.Coordinate{X: 70, Y: 70}, limit: 1024},
		{inputFile: "01.input", output: 22, description: "Puzzle case", end: lib.Coordinate{X: 6, Y: 6}, limit: 12, part2: true, part2Result: lib.Coordinate{X: 6, Y: 1}},
		{inputFile: "puzzle.input", output: 304, description: "Puzzle case", end: lib.Coordinate{X: 70, Y: 70}, limit: 1024, part2: true, part2Result: lib.Coordinate{X: 50, Y: 28}},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Measure start time
			startTime := time.Now()

			fmt.Println("Running test case", filepath.Join(dir, tc.inputFile))
			input, err := ParseInput(filepath.Join(dir, tc.inputFile), tc.end, tc.limit, tc.part2)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			output, _ := Solve(input, tc.part2)

			solution := tc.output

			if solution != output {
				t.Fatalf("Expected %v, got %v", solution, output)
			}

			if tc.part2 {
				i := 1
				for {
					// fmt.Println("Running test case", filepath.Join(dir, tc.inputFile), tc.limit+i)
					input, err := ParseInput(filepath.Join(dir, tc.inputFile), tc.end, tc.limit+i, tc.part2)
					if err != nil {
						t.Fatalf("Error: %v", err)
					}

					output, _ := Solve(input, tc.part2)

					if output == -1 {

						if input.Last != tc.part2Result {
							t.Fatalf("Expected %v, got %v", tc.part2Result, input.Last)
						}

						break
					}

					i++
				}
			}

			// Measure end time
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			// Log time and memory usage
			t.Logf("Time taken: %v", duration)
		})
	}
}
