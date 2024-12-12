package aoc2024_11

import (
	"math/big"
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/12/"

type testCase struct {
	inputFile   string
	output      string
	description string
	part2       bool
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		{inputFile: "01.input", output: "140", description: "Test case 1"},
		{inputFile: "02.input", output: "772", description: "Test case 2"},
		{inputFile: "03.input", output: "1930", description: "Test case 3"},
		{inputFile: "puzzle.input", output: "1434856", description: "Puzzle case"},
		{inputFile: "01.input", output: "80", description: "Test case 1", part2: true},
		{inputFile: "02.input", output: "436", description: "Test case 2", part2: true},
		{inputFile: "03.input", output: "1206", description: "Test case 3", part2: true},
		{inputFile: "04.input", output: "236", description: "Test case 4", part2: true},
		{inputFile: "05.input", output: "368", description: "Test case 5", part2: true},
		{inputFile: "06.input", output: "64", description: "Test case 6", part2: true},
		{inputFile: "puzzle.input", output: "891106", description: "Puzzle case", part2: true},
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

			solution := tc.output
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
