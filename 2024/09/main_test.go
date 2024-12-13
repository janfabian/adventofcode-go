package aoc2024

import (
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/09/"

type testCase struct {
	inputFile   string
	output      int
	description string
	part2       bool
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		{inputFile: "01.input", output: 1928, description: "Test case 1"},
		{inputFile: "02.input", output: 60, description: "Test case 2"},
		{inputFile: "03.input", output: 513, description: "Test case 3"},
		{inputFile: "puzzle.input", output: 6283404590840, description: "Puzzle case"},

		{inputFile: "01.input", output: 2858, description: "Test case 1", part2: true},
		{inputFile: "puzzle.input", output: 6304576012713, description: "Puzzle case", part2: true},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Measure start time
			startTime := time.Now()

			input, err := ParseInput(filepath.Join(dir, tc.inputFile), tc.part2)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			output := Solve(input, tc.part2)

			solution := tc.output

			if solution != output {
				t.Fatalf("Expected %v, got %v", solution, output)
			}

			// Measure end time
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			// Log time and memory usage
			t.Logf("Time taken: %v", duration)
		})
	}
}
