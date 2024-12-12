package aoc2024_11

import (
	"math/big"
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/11/"

type testCase struct {
	inputFile   string
	outputFile  string
	description string
	part2       bool
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		// {inputFile: "01.input", outputFile: "01.output", description: "Test case 1"},
		// {inputFile: "puzzle.input", outputFile: "puzzle.output", description: "Puzzle case"},
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

			solutionBigInt, _ := new(big.Int).SetString(solution, 10)

			if solutionBigInt.Cmp(output) != 0 {
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
