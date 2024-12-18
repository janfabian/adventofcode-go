package aoc2024

import (
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
	"time"
)

const dir = "2024/17/"

type testCase struct {
	inputFile    string
	output       string
	description  string
	part2        bool
	programAfter *Program
}

func Test(t *testing.T) {
	// Define your test cases
	testCases := []testCase{
		{inputFile: "01.input", output: "4,6,3,5,6,3,5,2,1,0", description: "Puzzle case"},
		{inputFile: "02.input", output: "", description: "Puzzle case", programAfter: &Program{
			Registers: []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(9)},
			Instructions: []*Instruction{
				{Opcode: 2, Operand: 6},
			},
			Pointer: 1,
		}},
		{inputFile: "03.input", output: "0,1,2", description: "Puzzle case", programAfter: &Program{
			Registers: []*big.Int{big.NewInt(10), big.NewInt(0), big.NewInt(0)},
			Instructions: []*Instruction{
				{Opcode: 5, Operand: 0},
				{Opcode: 5, Operand: 1},
				{Opcode: 5, Operand: 4},
			},
			Pointer: 3,
		}},
		{inputFile: "04.input", output: "4,2,5,6,7,7,7,7,3,1,0"},
		{inputFile: "05.input", output: "", programAfter: &Program{
			Registers: []*big.Int{big.NewInt(0), big.NewInt(26), big.NewInt(0)},
			Instructions: []*Instruction{
				{Opcode: 1, Operand: 7},
			},
			Pointer: 1,
		}},
		{inputFile: "06.input", output: "", programAfter: &Program{
			Registers: []*big.Int{big.NewInt(0), big.NewInt(44354), big.NewInt(43690)},
			Instructions: []*Instruction{
				{Opcode: 4, Operand: 0},
			},
			Pointer: 1,
		}},
		{inputFile: "07.input", output: "", programAfter: &Program{
			Registers: []*big.Int{big.NewInt(1024), big.NewInt(512), big.NewInt(512)},
			Instructions: []*Instruction{
				{Opcode: 0, Operand: 1},
				{Opcode: 6, Operand: 1},
				{Opcode: 7, Operand: 1},
			},
			Pointer: 3,
		}},
		{inputFile: "puzzle.input", output: "6,7,5,2,1,3,5,1,7", description: "Puzzle case"},
		{inputFile: "08.input", output: "117440", description: "Puzzle case", part2: true},
		{inputFile: "puzzle.input", output: "216549846240877", description: "Puzzle case", part2: true},
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

			if tc.programAfter != nil {
				if !input.Cmp(tc.programAfter) {
					t.Fatalf("Programs are not equal")
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
