package solver

import "jackson/internal/maze"

const (
	Width  = 6
	Height = 6

	FinishXFrom = 3
	FinishXTo   = 4

	FinishYFrom = 3
	FinishYTo   = 4
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
