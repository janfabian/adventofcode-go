package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"math"
	"strings"
)

type PriceList struct {
	PricePerMove int
	PricePerTurn int
}

type Robot struct {
	Position  *Coordinate
	Facing    *Vector
	PriceList *PriceList
}

func NewRobot(position Coordinate) *Robot {
	initRobotDir := *right

	return &Robot{
		Position: &position,
		Facing:   &initRobotDir,
		PriceList: &PriceList{
			PricePerMove: 1,
			PricePerTurn: 1000,
		},
	}
}

func (r *Robot) Clone() *Robot {
	return &Robot{
		Position: &Coordinate{
			X: r.Position.X,
			Y: r.Position.Y,
		},
		Facing: &Vector{
			dX: r.Facing.dX,
			dY: r.Facing.dY,
		},
		PriceList: r.PriceList,
	}
}

func (r *Robot) Turn(direction *Vector) int {
	if *direction == *r.Facing {
		return 0
	}

	turns := 1

	// opposite direction
	if direction.IsOpposite(r.Facing) {
		turns = 2
	}

	r.Facing = direction
	return turns * r.PriceList.PricePerTurn
}

func (r *Robot) Move(direction *Vector) int {
	turnPrice := r.Turn(direction)

	r.Position.X += direction.dX
	r.Position.Y += direction.dY

	return turnPrice + r.PriceList.PricePerMove
}

type Maze struct {
	Fence [][]bool
	Start *Coordinate
	End   *Coordinate
}

func (m *Maze) CanMove(position *Coordinate, direction *Vector) bool {
	newX := position.X + direction.dX
	newY := position.Y + direction.dY

	if newX < 0 || newX >= len(m.Fence[0]) || newY < 0 || newY >= len(m.Fence) {
		return false
	}

	if m.Fence[newY][newX] {
		return false
	}

	return true
}

func (m *Maze) PossibleMoves(position *Coordinate) []*Vector {
	possibleMoves := make([]*Vector, 0)

	for _, direction := range directions {
		if !m.CanMove(position, direction) {
			continue
		}

		possibleMoves = append(possibleMoves, direction)
	}

	return possibleMoves
}

func (m *Maze) Print(path *MazePath) {
	grid := lib.EmptyGrid(len(m.Fence[0]), len(m.Fence))

	for Y := 0; Y < len(m.Fence); Y++ {
		for X := 0; X < len(m.Fence[0]); X++ {
			if m.Fence[Y][X] {
				grid[Y][X] = '#'
			}
		}
	}

	grid[m.Start.Y][m.Start.X] = 'S'
	grid[m.End.Y][m.End.X] = 'E'

	robot := NewRobot(*m.Start)

	for _, move := range path.Path {
		robot.Move(move)

		if *move == *left {
			grid[robot.Position.Y][robot.Position.X] = '<'
		} else if *move == *right {
			grid[robot.Position.Y][robot.Position.X] = '>'
		} else if *move == *up {
			grid[robot.Position.Y][robot.Position.X] = '^'
		} else if *move == *down {
			grid[robot.Position.Y][robot.Position.X] = 'v'
		}

	}

	lib.PrintGridTable(grid)
}

type MazePath struct {
	Price int
	Path  []*Vector
	Robot *Robot
}

func (m *MazePath) Clone() *MazePath {
	newPath := &MazePath{
		Price: m.Price,
		Robot: m.Robot.Clone(),
		Path:  make([]*Vector, len(m.Path)),
	}

	copy(newPath.Path, m.Path)

	return newPath
}

func (m *Maze) FindPath(part2 bool) int {
	robot := NewRobot(*m.Start)

	allPaths := []*MazePath{{
		Price: 0,
		Path:  make([]*Vector, 0),
		Robot: robot,
	}}

	minPricePerCoord := make([][]int, len(m.Fence))
	for Y := 0; Y < len(m.Fence); Y++ {
		minPricePerCoord[Y] = make([]int, len(m.Fence[0]))
	}
	minPricePerCoord[m.Start.Y][m.Start.X] = -1

	cheapestPath := &MazePath{Price: math.MaxInt}
	allCheapPaths := make([]*MazePath, 0)

	for {
		if len(allPaths) == 0 {
			break
		}

		currentPath := allPaths[0]
		allPaths = allPaths[1:]

		// m.Print(currentPath)
		// fmt.Println("Price", currentPath.Price)

		if currentPath.Robot.Position.X == m.End.X && currentPath.Robot.Position.Y == m.End.Y {
			if currentPath.Price < cheapestPath.Price {
				cheapestPath = currentPath

				if len(allCheapPaths) > 0 {
					allCheapPaths = []*MazePath{}
				}
			}

			if currentPath.Price == cheapestPath.Price {
				allCheapPaths = append(allCheapPaths, currentPath)
			}

			continue
		}

		possibleMoves := m.PossibleMoves(currentPath.Robot.Position)

		for _, move := range possibleMoves {
			newPath := currentPath.Clone()
			newPath.Price += newPath.Robot.Move(move)
			newPath.Path = append(newPath.Path, move)

			nextX := newPath.Robot.Position.X
			nextY := newPath.Robot.Position.Y

			if minPricePerCoord[nextY][nextX] == 0 || (newPath.Price-newPath.Robot.PriceList.PricePerTurn) <= minPricePerCoord[nextY][nextX] {
				minPricePerCoord[nextY][nextX] = newPath.Price
				allPaths = append(allPaths, newPath)
			}
		}

	}

	if part2 {
		// include start
		total := 1
		visited := make([][]bool, len(m.Fence))
		for Y := 0; Y < len(m.Fence); Y++ {
			visited[Y] = make([]bool, len(m.Fence[0]))
		}

		for _, path := range allCheapPaths {

			robot := NewRobot(*m.Start)
			for _, move := range path.Path {
				robot.Move(move)
				if !visited[robot.Position.Y][robot.Position.X] {
					visited[robot.Position.Y][robot.Position.X] = true
					total++
				}
			}
		}

		return total
	}

	return cheapestPath.Price
}

type Coordinate struct {
	X, Y int
}

type Vector struct {
	dX, dY int
}

func (v *Vector) IsOpposite(other *Vector) bool {
	return v.dX*other.dY-v.dY*other.dX == 0
}

var left = &Vector{-1, 0}
var right = &Vector{1, 0}
var up = &Vector{0, -1}
var down = &Vector{0, 1}

var directions = []*Vector{left, right, up, down}

func Solve(maze *Maze, part2 bool) int {
	return maze.FindPath(part2)
}

func ParseInput(input string, part2 bool) (*Maze, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	maze := &Maze{
		Fence: make([][]bool, 0),
	}

	Y := 0
	for {
		X := 0
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		symbols := strings.Split(line, "")
		maze.Fence = append(maze.Fence, make([]bool, len(symbols)))

		for _, symbol := range symbols {
			if symbol == "S" {
				maze.Start = &Coordinate{
					X: X,
					Y: Y,
				}

				X++
				continue
			} else if symbol == "E" {
				maze.End = &Coordinate{
					X: X,
					Y: Y,
				}

				X++
				continue
			}

			maze.Fence[Y][X] = symbol == "#"

			X++
		}

		Y++

	}

	if maze.Start == nil {
		return nil, fmt.Errorf("start not found")
	}

	return maze, nil
}
