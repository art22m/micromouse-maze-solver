package main

import (
	"fmt"

	"jackson/internal/maze"
	"jackson/internal/solver"
)

func main() {
	fmt.Println("district")
	ff := solver.NewFloodFill(maze.Up, nil)
	ff.Solve()
}
