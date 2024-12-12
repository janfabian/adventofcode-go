package aoc2024_11

import (
	"adventofcode/lib"
	"fmt"
	"math/big"
	"strings"
	"sync"
)

func bigIntFromStr(s string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(s, 10)

	if !ok {
		return nil, fmt.Errorf("error converting to bigint: %v", s)
	}

	return n, nil
}

type Stone struct {
	stone *big.Int
	prev  *Stone
	next  *Stone
}

func (s *Stone) Occurences() map[string]*big.Int {
	occurences := map[string]*big.Int{}

	for s != nil {
		if _, ok := occurences[s.stone.String()]; !ok {
			occurences[s.stone.String()] = big.NewInt(0)
		}

		occurences[s.stone.String()].Add(occurences[s.stone.String()], big.NewInt(1))

		s = s.next
	}

	return occurences
}

func (s *Stone) Sum() *big.Int {
	sum := big.NewInt(0)

	for s != nil {
		sum.Add(sum, big.NewInt(1))
		s = s.next
	}

	return sum
}

func (s *Stone) Loop() error {
	for s != nil {
		var err error
		s, err = s.Iterate()

		if err != nil {
			return fmt.Errorf("error iterating: %v", err)
		}
	}

	return nil
}

func (s *Stone) Iterate() (*Stone, error) {
	if s.stone.Cmp(big.NewInt(0)) == 0 {
		s.stone = big.NewInt(1)

		return s.next, nil
	}

	nstr := s.stone.String()
	if len(nstr)%2 == 0 {
		midpoint := len(nstr) / 2
		part1 := nstr[:midpoint]
		part2 := nstr[midpoint:]

		s1N, err := bigIntFromStr(part1)

		if err != nil {
			return nil, err
		}

		s2N, err := bigIntFromStr(part2)

		if err != nil {
			return nil, err
		}

		s.stone = s1N

		// prev -> s -> next
		// prev -> s -> s2 -> next
		s2 := &Stone{stone: s2N, next: s.next, prev: s}

		originalNext := s.next

		if originalNext != nil {
			originalNext.prev = s2
		}

		s.next = s2

		return originalNext, nil
	}

	s.stone.Mul(s.stone, big.NewInt(2024))

	return s.next, nil
}

type Subprocessing struct {
	start      *Stone
	iterations int
	rounds     int
	cache      map[string]map[string]*big.Int
	mu         *sync.RWMutex
	result     *big.Int
}

func (s *Subprocessing) AddToCache(key string, value map[string]*big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.cache[key]; !ok {
		s.cache[key] = value
	}

}

func (s *Subprocessing) GetFromCache(key string) (map[string]*big.Int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.cache[key]
	return value, ok
}

func (s *Subprocessing) Loop() error {
	subprocesses := make([]*Subprocess, 0, 20000)
	subprocesses = append(subprocesses, &Subprocess{val: s.start.stone, n: big.NewInt(1)})
	wg := sync.WaitGroup{}
	lock := &sync.Mutex{}
	for r := 0; r < s.rounds; r++ {
		subprocessesAfterRound := []*Subprocess{}

		fmt.Println("START: subprocess length", len(subprocesses), "round", r)
		for _, sp := range subprocesses {
			if cache, ok := s.GetFromCache(sp.val.String()); ok {
				for stoneStr, occOfStone := range cache {
					if r == s.rounds-1 {
						lock.Lock()
						s.result.Add(s.result, new(big.Int).Mul(occOfStone, sp.n))
						lock.Unlock()
						continue
					}

					stoneN, err := bigIntFromStr(stoneStr)

					if err != nil {
						return err
					}

					lock.Lock()
					subprocessesAfterRound = append(subprocessesAfterRound, &Subprocess{val: stoneN, n: new(big.Int).Mul(occOfStone, sp.n)})
					lock.Unlock()
				}

				continue
			}

			wg.Add(1)
			lib.AddTask(func() {
				defer wg.Done()
				stone := &Stone{stone: new(big.Int).Set(sp.val)}
				for i := 0; i < s.iterations; i++ {
					err := stone.Loop()
					if err != nil {
						fmt.Println("error looping", err)
						return
					}

				}
				occ := stone.Occurences()

				s.AddToCache(sp.val.String(), occ)

				for stoneStr, occOfStone := range occ {
					if r == s.rounds-1 {
						lock.Lock()
						s.result.Add(s.result, new(big.Int).Mul(occOfStone, sp.n))
						lock.Unlock()
						continue
					}

					stoneN, err := bigIntFromStr(stoneStr)

					if err != nil {
						fmt.Println(err)
						return
					}

					lock.Lock()
					subprocessesAfterRound = append(subprocessesAfterRound, &Subprocess{val: stoneN, n: new(big.Int).Mul(occOfStone, sp.n)})
					lock.Unlock()
				}

			})
		}

		wg.Wait()

		fmt.Println("FINISH: subprocess length", len(subprocesses), "round", r)
		subprocesses = subprocessesAfterRound
	}

	return nil
}

type Subprocess struct {
	val *big.Int
	n   *big.Int
}

func Solve(input []*big.Int, part2 bool) *big.Int {

	lock := &sync.RWMutex{}
	cache := make(map[string]map[string]*big.Int, 10000)
	result := big.NewInt(0)
	rounds := 1
	iterations := 25

	if part2 {
		rounds = 3
	}

	wg := sync.WaitGroup{}
	result_lock := &sync.Mutex{}

	for _, stone := range input {
		wg.Add(1)
		stone := &Stone{stone: stone}
		subprocessing := &Subprocessing{start: stone, iterations: iterations, rounds: rounds, cache: cache, mu: lock, result: big.NewInt(0)}
		go func() {
			defer wg.Done()
			subprocessing.Loop()

			result_lock.Lock()
			result.Add(result, subprocessing.result)
			result_lock.Unlock()
		}()
	}

	wg.Wait()

	return result

}

func ParseInput(input string) ([]*big.Int, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := []*big.Int{}

	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, " ")

	for _, part := range parts {
		stone, ok := new(big.Int).SetString(part, 10)

		if !ok {
			return nil, fmt.Errorf("error converting to bigint: %v", part)
		}

		result = append(result, stone)
	}

	return result, nil
}

func ParseOutput(input string) (string, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	scanner.Scan()
	line := scanner.Text()
	return line, nil

}
