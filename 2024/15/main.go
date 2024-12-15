package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"strings"
)

type Warehouse struct {
	Robot   *Robot
	Objects []Object
	Area    [][]Object
}

func (w *Warehouse) SumBoxes() int {
	sum := 0
	for _, object := range w.Objects {
		if object.IsBox() {
			xPos := object.GetPosition()[0].X
			yPos := object.GetPosition()[0].Y

			sum += 100*yPos + xPos
		}
	}

	return sum
}

func (w *Warehouse) Print() {
	grid := lib.EmptyGrid(len(w.Area[0]), len(w.Area))

	for Y, row := range w.Area {
		for X, cell := range row {
			if cell == nil {
				grid[Y][X] = '.'
			} else {
				grid[Y][X] = cell.GetSymbol(&Coordinate{X, Y})
			}
		}
	}

	lib.PrintGridTable(grid)
}

func (w *Warehouse) MoveObject(o Object, v *Vector, simulate bool) bool {
	if !o.Movable() {
		return false
	}

	if v.dX == 0 && v.dY == 0 {
		return false
	}

	for _, coord := range o.GetEdgeInDir(v) {
		newX := coord.X + v.dX
		newY := coord.Y + v.dY

		if newX < 0 || newY < 0 || newX >= len(w.Area[0]) || newY >= len(w.Area) {
			return false
		}

		nextObject := w.Area[newY][newX]

		if nextObject != nil && !w.MoveObject(nextObject, v, true) {
			return false
		}
	}

	if simulate {
		return true
	}

	for _, coord := range o.GetPosition() {
		newX := coord.X + v.dX
		newY := coord.Y + v.dY

		nextObject := w.Area[newY][newX]
		if nextObject != nil {
			if nextObject != o {
				w.MoveObject(nextObject, v, false)
			}
		}

	}

	for _, coord := range o.GetPosition() {
		w.Area[coord.Y][coord.X] = nil
	}

	o.Move(v)

	for _, coord := range o.GetPosition() {
		w.Area[coord.Y][coord.X] = o
	}

	return true
}

type Object interface {
	Move(v *Vector)
	Movable() bool
	GetPosition() []*Coordinate
	GetSymbol(coord *Coordinate) rune
	GetEdgeInDir(dir *Vector) []*Coordinate
	IsBox() bool
}

type BaseObject struct {
	Position []*Coordinate
}

func (o *BaseObject) GetEdgeInDir(dir *Vector) []*Coordinate {
	switch *dir {
	case *left:
		return []*Coordinate{o.Position[0]}
	case *right:
		return []*Coordinate{o.Position[len(o.Position)-1]}
	}

	return o.Position
}

func (o *BaseObject) GetPosition() []*Coordinate {
	return o.Position
}

func (o *BaseObject) IsBox() bool {
	return false
}

func (o *BaseObject) Move(v *Vector) {
	for _, coord := range o.Position {
		coord.X += v.dX
		coord.Y += v.dY
	}
}

func (o *BaseObject) Movable() bool {
	return true
}

type Robot struct {
	*BaseObject
	Cmds []*Vector
}

func (r *Robot) GetSymbol(coord *Coordinate) rune {
	return '@'
}

type Box struct {
	*BaseObject
}

func (b *Box) GetSymbol(coord *Coordinate) rune {
	if len(b.Position) == 1 {
		return 'O'
	}

	if b.Position[0].X == coord.X && b.Position[0].Y == coord.Y {
		return '['
	} else {
		return ']'
	}
}

func (b *Box) IsBox() bool {
	return true
}

type Fence struct {
	*BaseObject
}

func (f *Fence) GetSymbol(coord *Coordinate) rune {
	return '#'
}

func (f *Fence) Move(v *Vector) {
	// no movement
}

func (f *Fence) Movable() bool {
	return false
}

type Coordinate struct {
	X, Y int
}

type Vector struct {
	dX, dY int
}

var left = &Vector{-1, 0}
var right = &Vector{1, 0}
var up = &Vector{0, -1}
var down = &Vector{0, 1}

func Solve(warehouse *Warehouse, part2 bool) int {
	// warehouse.Print()

	for _, cmd := range warehouse.Robot.Cmds {
		// fmt.Println("Move", cmd)
		warehouse.MoveObject(warehouse.Robot, cmd, false)
		// warehouse.Print()
	}

	return warehouse.SumBoxes()
}

func ParseInput(input string, part2 bool) (*Warehouse, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	warehouse := &Warehouse{
		Objects: make([]Object, 0),
		Area:    make([][]Object, 0),
	}

	Y := 0
	for {
		X := 0
		if !scanner.Scan() {
			return nil, fmt.Errorf("error reading file: %v", scanner.Err())
		}
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		symbols := strings.Split(line, "")
		warehouse.Area = append(warehouse.Area, make([]Object, 0))

		for _, symbol := range symbols {
			coords := []*Coordinate{{
				X: X,
				Y: Y,
			}}

			if part2 {
				switch symbol {
				case "#", "O", ".":
					X++
					coords = append(coords, &Coordinate{
						X: X,
						Y: Y,
					})
				}
			}

			base := &BaseObject{
				Position: coords,
			}

			var object Object

			switch symbol {
			case "@":
				robot := &Robot{
					BaseObject: base,
				}
				warehouse.Robot = robot
				object = robot
			case "#":
				object = &Fence{
					BaseObject: base,
				}
			case ".":
				// free space aka nil object
			case "O":
				object = &Box{
					BaseObject: base,
				}
			}

			if object != nil {
				warehouse.Objects = append(warehouse.Objects, object)
			}

			for range coords {
				warehouse.Area[Y] = append(warehouse.Area[Y], object)
			}

			if part2 {
				if symbol == "@" {
					X++
					warehouse.Area[Y] = append(warehouse.Area[Y], nil)
				}
			}

			X++
		}

		Y++

	}

	if warehouse.Robot == nil {
		return nil, fmt.Errorf("robot not found")
	}

	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		cmds := strings.Split(line, "")

		for _, cmd := range cmds {
			var vector *Vector

			switch cmd {
			case "<":
				vector = left
			case ">":
				vector = right
			case "^":
				vector = up
			case "v":
				vector = down
			}

			if vector != nil {
				warehouse.Robot.Cmds = append(warehouse.Robot.Cmds, vector)
			}
		}

	}

	return warehouse, nil
}
