package solver

import "jackson/internal/maze"

const (
	width  = 16
	height = 16

	finishXFrom = 7
	finishXTo   = 8

	finishYFrom = 7
	finishYTo   = 8
)

type Position struct {
	x, y int
}

func NewPosition(x, y int) Position {
	return Position{x: x, y: y}
}

type PositionWithDirection struct {
	Position
	maze.Direction
}
