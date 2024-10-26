package main

import (
	"jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/solver"
)

const (
	sensorsIP = ""
	motorsIP  = ""
	robotID   = ""
)

func main() {
	mover := mo.NewDummyMover(sensorsIP, motorsIP, robotID)
	startPosition := solver.NewPosition(0, 0)
	baseDirection := maze.Up

	ff := solver.NewFloodFill(
		baseDirection,
		startPosition,
		mover,
	)

	ff.Solve()
}
