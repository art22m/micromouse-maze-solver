package solver

import "jackson/internal/maze"

const (
	Width  = 6
	Height = 6

	StartX = 0
	StartY = 0

	FinishXFrom = 2
	FinishXTo   = 3

	FinishYFrom = 2
	FinishYTo   = 3
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
