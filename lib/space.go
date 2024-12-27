package lib

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) Add(v Vector) *Coordinate {
	return &Coordinate{X: c.X + v.DX, Y: c.Y + v.DY}
}

func (c *Coordinate) Around() []*Coordinate {
	var result []*Coordinate
	for _, dir := range Directions {
		result = append(result, c.Add(dir))
	}
	return result
}

type Vector struct {
	DX int
	DY int
}

var Left = Vector{DX: -1, DY: 0}
var Right = Vector{DX: 1, DY: 0}
var Up = Vector{DX: 0, DY: -1}
var Down = Vector{DX: 0, DY: 1}

var Directions = []Vector{Left, Right, Up, Down}
