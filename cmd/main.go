package main

import (
	"jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/solver"
)

const (
	sensorsIP = "localhost:8080"
	motorsIP  = "localhost:8080"
	robotID   = "1"
)

func main() {
	mover := mo.NewDummyMover(sensorsIP, motorsIP, robotID)

	config := solver.FloodFillConfig{
		StartDirection:  maze.Up,
		StartPosition:   solver.NewPosition(0, 0),
		MoveForwardOnly: false,
		Mover:           mover,
	}

	ff := solver.NewFloodFill(config)
	ff.Solve()
}
