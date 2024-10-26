package solver

import (
	"fmt"

	"jackson/internal/mover"
)

type FloodFill struct {
	maze [][]int
	mover.Mover
}

func NewFloodFill(mover mover.Mover) FloodFill {
	maze := make([][]int, Height)
	for i := 0; i < Height; i++ {
		maze[i] = make([]int, Width)
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			maze[i][j] = abs(i-getNearest(i, FinishXFrom, FinishXTo)) + abs(j-getNearest(j, FinishYFrom, FinishYTo))
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			fmt.Print(maze[i][j], " ")
		}
		fmt.Println()
	}

	return FloodFill{
		maze:  maze,
		Mover: mover,
	}
}

func (f *FloodFill) Solve() {
}
