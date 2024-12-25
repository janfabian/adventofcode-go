package aoc2024

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/22/"

type testCase struct {
	inputFile   string
	output      string
	description string
	part2       bool
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		// {inputFile: "00.input", output: "", description: "Puzzle case"},
		{inputFile: "01.input", output: "37327623", description: "Puzzle case"},
		// {inputFile: "03.input", output: "23", description: "Puzzle case", part2: true},
		{inputFile: "puzzle.input", output: "16894083306", description: "Puzzle case"},
		{inputFile: "puzzle.input", output: "1925", description: "Puzzle case", part2: true},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Measure start time
			startTime := time.Now()

			fmt.Println("Running test case", filepath.Join(dir, tc.inputFile))
			input, err := ParseInput(filepath.Join(dir, tc.inputFile), tc.part2)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			output, _ := Solve(input, tc.part2)

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
