package main

import (
	"jackson/internal/maze"
	"jackson/internal/solver"
)

func main() {
	ff := solver.NewFloodFill(maze.Up, solver.NewPosition(0, 0), nil)
	ff.Solve()
}
