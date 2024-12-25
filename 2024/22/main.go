package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"math/big"
	"strconv"
	"sync"
)

type Market struct {
	Buyers []*Buyer
}

func (m *Market) FindMaxBananas() ([4]int8, int) {
	tried := map[[4]int8]struct{}{}
	max := 0
	var maxSeq [4]int8
	var wg sync.WaitGroup
	var mu sync.RWMutex
	for _, b := range m.Buyers {
		wg.Add(1)
		lib.AddTask(func() {
			defer wg.Done()
			for seq := range b.Bananas {
				mu.RLock()
				if _, ok := tried[seq]; ok {
					mu.RUnlock()
					continue
				}
				mu.RUnlock()

				acc := 0

				for _, b2 := range m.Buyers {
					if val, ok := b2.Bananas[seq]; ok {
						acc += int(val)
					}
				}

				mu.Lock()
				if acc > max {
					max = acc
					maxSeq = seq
				}
				tried[seq] = struct{}{}
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	return maxSeq, max
}

func (m *Market) Iterate(n int) {
	var wg sync.WaitGroup
	for _, b := range m.Buyers {
		wg.Add(1)
		lib.AddTask(
			func() {

				var prev int8
				diffs := [4]int8{}

				for i := 0; i < n; i++ {
					prev = int8(new(big.Int).Mod(b.Amount, big.NewInt(10)).Uint64())
					b.Next()
					next := int8(new(big.Int).Mod(b.Amount, big.NewInt(10)).Uint64())
					diff := next - prev

					diffs[0] = diffs[1]
					diffs[1] = diffs[2]
					diffs[2] = diffs[3]
					diffs[3] = diff

					if i < 3 {
						continue
					}

					if _, ok := b.Bananas[diffs]; !ok {
						b.Bananas[diffs] = next
					}

				}
				wg.Done()
				// }()
			})
	}

	wg.Wait()
}

type Buyer struct {
	Amount  *big.Int
	Mod     *big.Int
	Bananas map[[4]int8]int8
}

func (b *Buyer) Next() {
	// fmt.Println("=============")
	// fmt.Println("next", b.Amount.String(), b.Amount.Text(2))
	m := new(big.Int).Mul(b.Amount, big.NewInt(64))
	// fmt.Println("multiply by 64", m.String(), m.Text(2))
	b.Mix(m)
	// fmt.Println("mix", b.Amount.String(), b.Amount.Text(2))
	b.Prune()
	// fmt.Println("prune", b.Amount.String(), b.Amount.Text(2))
	d := new(big.Int).Div(b.Amount, big.NewInt(32))
	// fmt.Println("divide by 32", d.String(), d.Text(2))
	b.Mix(d)
	// fmt.Println("mix", b.Amount.String(), b.Amount.Text(2))
	b.Prune()
	// fmt.Println("prune", b.Amount.String(), b.Amount.Text(2))
	m = new(big.Int).Mul(b.Amount, big.NewInt(2048))
	// fmt.Println("multiply by 2048", m.String(), m.Text(2))
	b.Mix(m)
	// fmt.Println("mix", b.Amount.String(), b.Amount.Text(2))
	b.Prune()
	// fmt.Println("prune", b.Amount.String(), b.Amount.Text(2))
}

func (b *Buyer) Mix(n *big.Int) {
	b.Amount.Xor(b.Amount, n)
}

func (b *Buyer) Prune() {
	b.Amount.Mod(b.Amount, b.Mod)
}

func Solve(program *Market, part2 bool) (string, error) {
	program.Iterate(2000)

	result := big.NewInt(0)

	for _, b := range program.Buyers {
		result.Add(result, b.Amount)
	}

	if !part2 {
		return result.String(), nil
	}

	s, max := program.FindMaxBananas()

	fmt.Println("max", max, s)

	return strconv.Itoa(max), nil

}

func ParseInput(input string, part2 bool) (*Market, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	market := &Market{
		Buyers: []*Buyer{},
	}

	for scanner.Scan() {
		line := scanner.Text()

		amount, success := new(big.Int).SetString(line, 10)

		if !success {
			return nil, fmt.Errorf("invalid buyer: %s", line)
		}

		market.Buyers = append(market.Buyers, &Buyer{
			Amount:  amount,
			Mod:     big.NewInt(16777216),
			Bananas: map[[4]int8]int8{},
		})
	}

	return market, nil
}
