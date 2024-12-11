package aoc2024_10

import (
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/10/"

type testCase struct {
	inputFile   string
	outputFile  string
	description string
	part2       bool
}

func TestPart1(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		{inputFile: "01.input", outputFile: "01.output", description: "Test case 1"},
		{inputFile: "puzzle.input", outputFile: "puzzle.output", description: "Puzzle case"},
		{inputFile: "01_part2.input", outputFile: "01_part2.output", description: "Test case 1", part2: true},
		{inputFile: "puzzle.input", outputFile: "puzzle_part2.output", description: "Puzzle case", part2: true},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Measure start time
			startTime := time.Now()

			input, err := ParseInput(filepath.Join(dir, tc.inputFile))
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			output := Solve(input, tc.part2)

			solution, err := ParseOutput(filepath.Join(dir, tc.outputFile))
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			if output != solution {
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
