package aoc2024_11

import (
	"adventofcode/lib"
	"fmt"
	"math/big"
	"sort"
	"sync"
)

type Coordinate struct {
	X, Y int
}

type ByX []*Coordinate

func (a ByX) Len() int           { return len(a) }
func (a ByX) Less(i, j int) bool { return a[i].X < a[j].X }
func (a ByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByY []*Coordinate

func (a ByY) Len() int           { return len(a) }
func (a ByY) Less(i, j int) bool { return a[i].Y < a[j].Y }
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// bit mask
var left = 0b0001
var right = 0b0010
var up = 0b0100
var down = 0b1000

// var directions = []int{left, right, up, down}

func (c *Coordinate) Compare(other *Coordinate) int {
	if c.X == other.X-1 {
		return right
	}

	if c.X == other.X+1 {
		return left
	}

	if c.Y == other.Y-1 {
		return down
	}

	if c.Y == other.Y+1 {
		return up
	}

	fmt.Println("Error comparing coordinates")
	return 0
}

type Region struct {
	coords  []*Coordinate
	xCoords map[int][]*Coordinate
	yCoords map[int][]*Coordinate
	plant   rune
}

func NewRegion(coord *Coordinate, plant rune) *Region {
	return &Region{
		plant:   plant,
		coords:  []*Coordinate{coord},
		xCoords: map[int][]*Coordinate{coord.X: {coord}},
		yCoords: map[int][]*Coordinate{coord.Y: {coord}},
	}
}

func (r *Region) Add(coord *Coordinate) {
	r.coords = append(r.coords, coord)
	r.xCoords[coord.X] = append(r.xCoords[coord.X], coord)
	sort.Sort(ByY(r.xCoords[coord.X]))
	r.yCoords[coord.Y] = append(r.yCoords[coord.Y], coord)
	sort.Sort(ByX(r.yCoords[coord.Y]))
}

func (r *Region) Area() *big.Int {
	return big.NewInt(int64(len(r.coords)))
}

type Farm struct {
	plants [][]rune
}

func (f *Farm) GetWidth() int {
	return len(f.plants[0])
}

func (f *Farm) GetHeight() int {
	return len(f.plants)
}

func (f *Farm) CheckNeighbors(coord *Coordinate, plant rune) []Coordinate {
	region := []Coordinate{}
	if coord.X > 0 {
		plantLeft := f.plants[coord.Y][coord.X-1]

		if plantLeft == plant {
			region = append(region, Coordinate{coord.X - 1, coord.Y})
		}
	}

	if coord.X < len(f.plants[0])-1 {
		plantRight := f.plants[coord.Y][coord.X+1]

		if plantRight == plant {
			region = append(region, Coordinate{coord.X + 1, coord.Y})
		}
	}

	if coord.Y > 0 {
		plantUp := f.plants[coord.Y-1][coord.X]

		if plantUp == plant {
			region = append(region, Coordinate{coord.X, coord.Y - 1})
		}
	}

	if coord.Y < len(f.plants)-1 {
		plantDown := f.plants[coord.Y+1][coord.X]

		if plantDown == plant {
			region = append(region, Coordinate{coord.X, coord.Y + 1})
		}
	}

	return region
}

func (f *Farm) FindRegions() []*Region {
	regions := []*Region{}
	visited := make(map[Coordinate]struct{})

	for Y, line := range f.plants {
		for X, plant := range line {
			coord := Coordinate{X, Y}
			if _, ok := visited[coord]; ok {
				continue
			}
			visited[coord] = struct{}{}

			region := NewRegion(&coord, plant)

			regionCoords := []*Coordinate{&coord}

			for len(regionCoords) > 0 {

				regionCoord := regionCoords[0]
				regionCoords = regionCoords[1:]
				regionNeighbords := f.CheckNeighbors(regionCoord, plant)

				for _, neighbor := range regionNeighbords {
					if _, ok := visited[neighbor]; ok {
						continue
					}
					visited[neighbor] = struct{}{}
					regionCoords = append(regionCoords, &neighbor)
					region.Add(&neighbor)
				}

			}

			regions = append(regions, region)
		}
	}

	return regions
}

func (f *Farm) FencePrice(r *Region) *big.Int {
	perimeter := big.NewInt(0)

	for _, coord := range r.coords {
		n := f.CheckNeighbors(coord, r.plant)

		perimeter.Add(perimeter, big.NewInt(int64(4-len(n))))
	}

	return new(big.Int).Mul(perimeter, r.Area())
}

func (f *Farm) FencePriceDiscount(r *Region) *big.Int {
	perimeter := big.NewInt(0)

	for i := 0; i < f.GetWidth(); i++ {
		column := r.xCoords[i]

		if len(column) == 0 {
			continue
		}

		prevLeftDiscount := false

		for _, coord := range column {
			n := f.CheckNeighbors(coord, r.plant)

			neighborsMask := 0

			for _, neighbor := range n {
				neighborsMask |= coord.Compare(&neighbor)
			}

			fenceMask := ^neighborsMask
			hasDiscountableNeighbors := neighborsMask&(up) != 0

			if fenceMask&left == left {
				if !prevLeftDiscount || !hasDiscountableNeighbors {
					perimeter.Add(perimeter, big.NewInt(1))
					prevLeftDiscount = true
				}
			} else {
				prevLeftDiscount = false
			}
		}
	}

	// traverse right to left
	for i := f.GetWidth() - 1; i >= 0; i-- {
		column := r.xCoords[i]

		if len(column) == 0 {
			continue
		}

		prevRightDiscount := false

		for _, coord := range column {
			n := f.CheckNeighbors(coord, r.plant)

			neighborsMask := 0

			for _, neighbor := range n {
				neighborsMask |= coord.Compare(&neighbor)
			}

			fenceMask := ^neighborsMask
			hasDiscountableNeighbors := neighborsMask&(up) != 0

			if fenceMask&right == right {
				if !prevRightDiscount || !hasDiscountableNeighbors {
					perimeter.Add(perimeter, big.NewInt(1))
					prevRightDiscount = true
				}
			} else {
				prevRightDiscount = false
			}
		}
	}

	// traverse top to bottom
	for i := 0; i < f.GetHeight(); i++ {
		row := r.yCoords[i]

		if len(row) == 0 {
			continue
		}

		prevUpDiscount := false

		for _, coord := range row {
			n := f.CheckNeighbors(coord, r.plant)

			neighborsMask := 0

			for _, neighbor := range n {
				neighborsMask |= coord.Compare(&neighbor)
			}

			fenceMask := ^neighborsMask
			hasDiscountableNeighbors := neighborsMask&(left) != 0

			if fenceMask&up == up {
				if !prevUpDiscount || !hasDiscountableNeighbors {
					perimeter.Add(perimeter, big.NewInt(1))
					prevUpDiscount = true

				}
			} else {
				prevUpDiscount = false
			}

		}
	}

	// traverse bottom to top
	for i := f.GetHeight() - 1; i >= 0; i-- {
		row := r.yCoords[i]

		if len(row) == 0 {
			continue
		}

		prevDownDiscount := false

		for _, coord := range row {
			n := f.CheckNeighbors(coord, r.plant)

			neighborsMask := 0

			for _, neighbor := range n {
				neighborsMask |= coord.Compare(&neighbor)
			}

			fenceMask := ^neighborsMask
			hasDiscountableNeighbors := neighborsMask&(left) != 0

			if fenceMask&down == down {
				if !prevDownDiscount || !hasDiscountableNeighbors {
					perimeter.Add(perimeter, big.NewInt(1))
					prevDownDiscount = true
				}
			} else {
				prevDownDiscount = false
			}

		}
	}

	// fmt.Println("Plant", string(r.plant), "Perimeter", perimeter, "Area", r.Area())

	return new(big.Int).Mul(perimeter, r.Area())
}

func Solve(input [][]rune, part2 bool) *big.Int {
	farm := Farm{plants: input}

	regions := farm.FindRegions()

	totalPrice := big.NewInt(0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, r := range regions {
		wg.Add(1)
		go func(r *Region) {
			defer wg.Done()
			var fencePrice *big.Int

			if part2 {
				fencePrice = farm.FencePriceDiscount(r)
			} else {
				fencePrice = farm.FencePrice(r)
			}

			mu.Lock()
			totalPrice.Add(totalPrice, fencePrice)
			mu.Unlock()

		}(r)
	}

	wg.Wait()

	return totalPrice
}

func ParseInput(input string) ([][]rune, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := [][]rune{}

	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, []rune(line))
	}

	return result, nil
}
