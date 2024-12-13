package aoc2024_11

import (
	"adventofcode/lib"
	"fmt"
	"math/big"
	"regexp"
)

type ClawMachine struct {
	A     Vector
	B     Vector
	Price Coordinate
	CostA *big.Int
	CostB *big.Int
}

func (c *ClawMachine) CalcMovements() (*Vector, error) {

	numerator := new(big.Int).Sub(
		new(big.Int).Mul(c.A.dX, c.Price.Y),
		new(big.Int).Mul(c.A.dY, c.Price.X),
	)

	denominator := new(big.Int).Sub(
		new(big.Int).Mul(c.A.dX, c.B.dY),
		new(big.Int).Mul(c.A.dY, c.B.dX),
	)

	quotient, remainder := new(big.Int), new(big.Int)
	quotient.QuoRem(numerator, denominator, remainder)

	if remainder.Sign() != 0 {
		return nil, fmt.Errorf("no solution")
	}

	b := quotient

	numerator = new(big.Int).Sub(
		c.Price.X,
		new(big.Int).Mul(b, c.B.dX),
	)
	denominator = c.A.dX

	quotient, remainder = new(big.Int), new(big.Int)
	quotient.QuoRem(numerator, denominator, remainder)

	if remainder.Sign() != 0 {
		return nil, fmt.Errorf("no solution")
	}

	a := quotient

	return &Vector{dX: a, dY: b}, nil
}

func (c *ClawMachine) CalcCost(v *Vector) *big.Int {
	return new(big.Int).Add(
		new(big.Int).Mul(c.CostA, v.dX),
		new(big.Int).Mul(c.CostB, v.dY),
	)
}

type Coordinate struct {
	X, Y *big.Int
}

type Vector struct {
	dX, dY *big.Int
}

func Solve(input []*ClawMachine, part2 bool) *big.Int {
	result := big.NewInt(0)

	for _, clawMachine := range input {
		clawResult, err := clawMachine.CalcMovements()

		if err != nil {
			continue
		}

		result.Add(result, clawMachine.CalcCost(clawResult))

	}

	return result
}

func ParseInput(input string, part2 bool) ([]*ClawMachine, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := []*ClawMachine{}

	// Define the regular expression pattern
	vectorRe := regexp.MustCompile(`\+(\d+)`)
	priceRe := regexp.MustCompile(`=(\d+)`)

	for {
		// A
		scanner.Scan()
		line := scanner.Text()
		match := vectorRe.FindAllStringSubmatch(line, 2)

		dX, ok := new(big.Int).SetString(match[0][1], 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		dY, ok := new(big.Int).SetString(match[1][1], 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		A := Vector{dX: dX, dY: dY}

		// B
		scanner.Scan()
		line = scanner.Text()
		match = vectorRe.FindAllStringSubmatch(line, 2)

		dX, ok = new(big.Int).SetString(match[0][1], 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		dY, ok = new(big.Int).SetString(match[1][1], 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		B := Vector{dX: dX, dY: dY}

		// Price
		scanner.Scan()
		line = scanner.Text()
		match = priceRe.FindAllStringSubmatch(line, 2)

		priceX := match[0][1]
		X, ok := new(big.Int).SetString(priceX, 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		if part2 {
			X.Add(X, big.NewInt(10000000000000))
		}

		priceY := match[1][1]
		Y, ok := new(big.Int).SetString(priceY, 10)
		if !ok {
			return nil, fmt.Errorf("error converting to big.Int: %v", ok)
		}

		if part2 {
			Y.Add(Y, big.NewInt(10000000000000))
		}

		Price := Coordinate{X, Y}

		result = append(result, &ClawMachine{A, B, Price, big.NewInt(3), big.NewInt(1)})

		if !scanner.Scan() {
			break
		}
	}

	return result, nil
}
