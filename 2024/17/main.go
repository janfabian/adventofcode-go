package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type Program struct {
	Instructions []*Instruction
	Pointer      int
	Registers    []*big.Int
}

func (p *Program) PrintInstructions(delimiter string) string {
	output := ""

	for i := 0; i < len(p.Instructions); i++ {
		if len(output) > 0 {
			output += delimiter
		}

		output += fmt.Sprintf("%d%s%d", p.Instructions[i].Opcode, delimiter, p.Instructions[i].Operand)
	}

	return output
}

func (p *Program) Copy() *Program {
	newProgram := &Program{
		Pointer:   p.Pointer,
		Registers: make([]*big.Int, len(p.Registers)),
	}

	for i := 0; i < len(p.Registers); i++ {
		newProgram.Registers[i] = new(big.Int).Set(p.Registers[i])
	}

	newProgram.Instructions = make([]*Instruction, len(p.Instructions))

	for i := 0; i < len(p.Instructions); i++ {
		newProgram.Instructions[i] = &Instruction{
			Opcode:  p.Instructions[i].Opcode,
			Operand: p.Instructions[i].Operand,
		}
	}

	return newProgram
}

func (p *Program) Cmp(other *Program) bool {
	if len(p.Registers) != len(other.Registers) {
		return false
	}

	for i := 0; i < len(p.Registers); i++ {
		if p.Registers[i].Cmp(other.Registers[i]) != 0 {
			return false
		}
	}

	if p.Pointer != other.Pointer {
		return false
	}

	if len(p.Instructions) != len(other.Instructions) {
		return false
	}

	for i := 0; i < len(p.Instructions); i++ {
		if p.Instructions[i].Opcode != other.Instructions[i].Opcode || p.Instructions[i].Operand != other.Instructions[i].Operand {
			return false
		}
	}

	return true

}

func (p *Program) Run(delimiter string) string {
	output := ""

	for p.Pointer < len(p.Instructions) {
		instruction := p.Instructions[p.Pointer]

		switch instruction.Opcode {
		case 0:
			p.adv(instruction.Operand)
		case 1:
			p.bxl(instruction.Operand)
		case 2:
			p.bst(instruction.Operand)
		case 3:
			if p.jnz(instruction.Operand) {
				continue
			}
		case 4:
			p.bxc()
		case 5:
			if len(output) > 0 {
				output += delimiter
			}

			output += p.out(instruction.Operand)
		case 6:
			p.bdv(instruction.Operand)
		case 7:
			p.cdv(instruction.Operand)
		}

		p.Pointer++
	}

	return output
}

func (p *Program) adv_calc(operand int8) *big.Int {
	num := p.registerA()
	den := new(big.Int).Exp(big.NewInt(2), p.comboOp(operand), nil)

	return new(big.Int).Div(num, den)
}

func (p *Program) adv(operand int8) {
	A := p.registerA()
	A.Set(p.adv_calc(operand))
}

func (p *Program) bxl(operand int8) {
	B := p.registerB()
	B.Xor(B, big.NewInt(int64(operand)))
}

func (p *Program) bst(operand int8) {
	B := p.registerB()
	B.Mod(p.comboOp(operand), big.NewInt(8))
}

func (p *Program) jnz(operand int8) bool {
	if p.registerA().Sign() != 0 {
		p.Pointer = int(operand) / 2
		return true
	}

	return false
}

func (p *Program) bxc() {
	C := p.registerC()
	B := p.registerB()

	B.Xor(B, C)
}

func (p *Program) out(operand int8) string {
	return new(big.Int).Mod(p.comboOp(operand), big.NewInt(8)).String()
}

func (p *Program) bdv(operand int8) {
	B := p.registerB()
	B.Set(p.adv_calc(operand))
}

func (p *Program) cdv(operand int8) {
	C := p.registerC()
	C.Set(p.adv_calc(operand))
}

func (p *Program) registerA() *big.Int {
	return p.Registers[0]
}

func (p *Program) registerB() *big.Int {
	return p.Registers[1]
}

func (p *Program) registerC() *big.Int {
	return p.Registers[2]
}

func (p *Program) comboOp(op int8) *big.Int {
	switch op {
	case 0:
		return big.NewInt(0)
	case 1:
		return big.NewInt(1)
	case 2:
		return big.NewInt(2)
	case 3:
		return big.NewInt(3)
	case 4:
		return new(big.Int).Set(p.registerA())
	case 5:
		return new(big.Int).Set(p.registerB())
	case 6:
		return new(big.Int).Set(p.registerC())
	}

	panic("invalid opcode")
}

type Instruction struct {
	Opcode  int8
	Operand int8
}

func createBigIntWithBinaryPrefixPostfix(binaryPrefix, binaryPostfix string) (*big.Int, error) {
	binaryString := binaryPrefix + binaryPostfix

	B := new(big.Int)
	_, ok := B.SetString(binaryString, 2)
	if !ok {
		return nil, fmt.Errorf("error converting binary string to big.Int")
	}

	return B, nil
}

func Solve(program *Program, part2 bool) (string, error) {
	if !part2 {
		return program.Run(","), nil
	}

	wanted := program.PrintInstructions("")

	n := 6
	postfixes := []string{}

	for i := 1; i <= n; i++ {
		maxJ := int(math.Pow(float64(2), float64(i)))
		for j := 0; j < maxJ; j++ {
			binary := fmt.Sprintf("%0*b", i, j)
			postfixes = append(postfixes, binary)
		}
	}

	prefixes := []string{""}
	tried := map[string]struct{}{}
	_ = prefixes

	var A *big.Int
	for {

		if len(prefixes) == 0 {
			return "", fmt.Errorf("no more prefixes to try")
		}

		prefix := prefixes[0]
		prefixes = prefixes[1:]

		for _, postfix := range postfixes {
			A, _ = createBigIntWithBinaryPrefixPostfix(prefix, postfix)
			programCopy := program.Copy()
			programCopy.Registers[0].Set(A)

			output := programCopy.Run("")

			if strings.HasSuffix(wanted, output) {

				if output == wanted {
					return A.Text(10), nil
				}

				// fmt.Printf("i: %s, (%s), output: %s\n", A.Text(10), A.Text(2), output)
				nextPrefix := A.Text(2)
				if len(nextPrefix) > 3 {
					nextPrefix = nextPrefix[:len(nextPrefix)-3]
				}

				if _, ok := tried[nextPrefix]; !ok {
					prefixes = append(prefixes, nextPrefix)
					tried[nextPrefix] = struct{}{}
				}
			}
		}

	}

}

func ParseInput(input string, part2 bool) (*Program, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	program := &Program{
		Registers: make([]*big.Int, 3),
	}

	registerRegex := regexp.MustCompile(`(\d+)$`)

	for i := 0; i < 3; i++ {
		scanner.Scan()
		line := scanner.Text()
		matches := registerRegex.FindStringSubmatch(line)

		if len(matches) != 2 {
			return nil, fmt.Errorf("invalid register line: %s", line)
		}

		register, success := new(big.Int).SetString(matches[1], 10)

		if !success {
			return nil, fmt.Errorf("invalid register value: %s", matches[1])
		}

		program.Registers[i] = register
	}

	// empty line
	scanner.Scan()
	scanner.Scan()
	line := strings.TrimPrefix(scanner.Text(), "Program: ")

	parts := strings.Split(line, ",")

	if len(parts)%2 != 0 {
		return nil, fmt.Errorf("invalid program: %s", line)
	}

	for i := 0; i < len(parts); i += 2 {
		opcode, err := strconv.Atoi(parts[i])

		if err != nil {
			return nil, fmt.Errorf("invalid opcode: %s", parts[i])
		}

		operand, err := strconv.Atoi(parts[i+1])

		if err != nil {
			return nil, fmt.Errorf("invalid operand: %s", parts[i+1])
		}

		program.Instructions = append(program.Instructions, &Instruction{
			Opcode:  int8(opcode),
			Operand: int8(operand),
		})
	}

	return program, nil
}
