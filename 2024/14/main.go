package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Robot struct {
	Position  Coordinate
	Direction Vector
}

type Coordinate struct {
	X, Y int
}

type Vector struct {
	dX, dY int
}

func printGridTable(grid [][]rune) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)

	for _, row := range grid {
		for _, cell := range row {
			fmt.Fprintf(w, "%c\t", cell)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
}

func emptyGrid(X, Y int) [][]rune {

	grid := [][]rune{}

	for i := 0; i < Y; i++ {
		row := make([]rune, X)
		for j := 0; j < X; j++ {
			row[j] = ' '
		}
		grid = append(grid, row)
	}

	return grid

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func avg(input []int) int {
	sum := 0
	for _, i := range input {
		sum += i
	}
	return sum / len(input)
}

func calculateDistances(input []*Robot) int {
	dist := 0
	for i := 0; i < len(input); i++ {
		for j := i + 1; j < len(input); j++ {
			dist += abs(input[i].Position.X-input[j].Position.X) + abs(input[i].Position.Y-input[j].Position.Y)
		}
	}

	return dist
}

func Solve(input []*Robot, part2 bool) int {
	X := 101
	Y := 103
	S := 100

	if part2 {
		S = 1
	}

	halfX := X / 2
	halfY := Y / 2

	quadrants := make([]int, 4)
	r := 0.7

	if part2 {
		all_dists := make([]int, 0, 10000)
		counter := 0
		for {
			grid := emptyGrid(X, Y)

			for _, robot := range input {
				robot.Position.X = (robot.Position.X + S*robot.Direction.dX + S*X) % X
				robot.Position.Y = (robot.Position.Y + S*robot.Direction.dY + S*Y) % Y

				grid[robot.Position.Y][robot.Position.X] = '#'
			}
			counter++
			// Print the grid table

			dist := calculateDistances(input)
			all_dists = append(all_dists, dist)

			if float64(dist) < float64(avg(all_dists))*r {
				printGridTable(grid)
				fmt.Println("Distances:", dist)
				fmt.Println("Counter:", counter)
				return counter
			}

		}
	}

	for _, robot := range input {
		robot.Position.X = (robot.Position.X + S*robot.Direction.dX + S*X) % X
		robot.Position.Y = (robot.Position.Y + S*robot.Direction.dY + S*Y) % Y

		if robot.Position.X < halfX {
			if robot.Position.Y < halfY {
				quadrants[0]++
			} else if robot.Position.Y > halfY {
				quadrants[2]++
			}
		} else if robot.Position.X > halfX {
			if robot.Position.Y < halfY {
				quadrants[1]++
			} else if robot.Position.Y > halfY {
				quadrants[3]++
			}
		}
	}

	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func ParseInput(input string, part2 bool) ([]*Robot, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := []*Robot{}

	lineRe := regexp.MustCompile(`=(.*)\ .*=(.*)$`)

	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		match := lineRe.FindAllStringSubmatch(line, 2)

		p_s, v_s := strings.Split(match[0][1], ","), strings.Split(match[0][2], ",")
		p1, _ := strconv.Atoi(p_s[0])
		p2, _ := strconv.Atoi(p_s[1])
		v1, _ := strconv.Atoi(v_s[0])
		v2, _ := strconv.Atoi(v_s[1])
		result = append(result, &Robot{Coordinate{p1, p2}, Vector{v1, v2}})

	}

	return result, nil
}
