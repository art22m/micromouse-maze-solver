package solver

import (
	"fmt"

	"jackson/internal/maze"
)

const (
	width  = 16
	height = 16

	finishXFrom = 7
	finishXTo   = 8
	finishYFrom = 7
	finishYTo   = 8

	//width  = 8
	//height = 8
	//
	//finishXFrom = 3
	//finishXTo   = 4
	//finishYFrom = 3
	//finishYTo   = 4
)

type Position struct {
	x, y int
}

func NewPosition(x, y int) Position {
	return Position{x: x, y: y}
}

func (p Position) Shift(dir maze.Direction) Position {
	switch dir {
	case maze.Up:
		return Position{p.x + 1, p.y}
	case maze.Right:
		return Position{p.x, p.y + 1}
	case maze.Down:
		return Position{p.x - 1, p.y}
	case maze.Left:
		return Position{p.x, p.y - 1}
	}
	return p
}

func (p Position) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}

type PositionWithDirection struct {
	Position
	maze.Direction
}

func (p PositionWithDirection) String() string {
	return fmt.Sprintf("dir=%v, pos=%v", p.Direction.String(), p.Position.String())
}
