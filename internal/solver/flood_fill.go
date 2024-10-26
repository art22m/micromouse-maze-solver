package solver

import (
	"fmt"
	"log"

	mz "jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/stack"
)

type FloodFill struct {
	flood [][]int
	cells [][]mz.Wall

	or mz.Orientation
	mo mo.Mover

	Position
}

func NewFloodFill(or mz.Orientation, mover mo.Mover) FloodFill {
	flood := make([][]int, Height)
	cells := make([][]mz.Wall, Height)
	for i := 0; i < Height; i++ {
		flood[i] = make([]int, Width)
		cells[i] = make([]mz.Wall, Width)
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			flood[i][j] = abs(i-getNearest(i, FinishXFrom, FinishXTo)) + abs(j-getNearest(j, FinishYFrom, FinishYTo))
			cells[i][j] = mz.Unknown
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			fmt.Print(flood[i][j], " ")
		}
		fmt.Println()
	}

	return FloodFill{
		flood: flood,
		cells: cells,
		mo:    mover,
		or:    or,
	}
}

func (f *FloodFill) Solve() {
	for {
		state := f.mo.CellState()
		f.updateWalls(state)
		if f.getFlood() == 0 {
			log.Println("end")
			break
		}

	}
}

func (f *FloodFill) updateWalls(state mo.Cell) {
	f.cells[f.x][f.y] = state.Wall
}

func (f *FloodFill) floodFill() {
	st := stack.Stack[Position]{}
	st.Push(Position{f.x, f.y})
	for !st.Empty() {

	}
}

func (f *FloodFill) getFlood() int {
	return f.flood[f.x][f.y]
}

func (f *FloodFill) getCell() mz.Wall {
	return f.cells[f.x][f.y]
}

//func (f *FloodFill) minOpenNeighbour() Position {
//
//}
