package solver

import "jackson/internal/maze"

const (
	Width  = 16
	Height = 16

	FinishXFrom = 7
	FinishXTo   = 8

	FinishYFrom = 7
	FinishYTo   = 8
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
