package solver

import (
	"fmt"
	"log"
	"math"

	ma "jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/stack"
)

type FloodFill struct {
	flood [][]int
	cells [][]ma.Wall

	or ma.Orientation
	mo mo.Mover

	Position
}

func NewFloodFill(or ma.Orientation, mover mo.Mover) FloodFill {
	flood := make([][]int, Height)
	cells := make([][]ma.Wall, Height)
	for i := 0; i < Height; i++ {
		flood[i] = make([]int, Width)
		cells[i] = make([]ma.Wall, Width)
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			flood[i][j] = abs(i-getNearest(i, FinishXFrom, FinishXTo)) + abs(j-getNearest(j, FinishYFrom, FinishYTo))
			cells[i][j] = ma.Unknown
		}
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
		if f.getFlood(f.Position) == 0 {
			log.Println("!!end")
			break
		}

		f.floodFill()

	}
}

func (f *FloodFill) getAndUpdateWalls() {
	state := f.mo.CellState()
	f.updateWallsIfNeeded(f.Position, state.Wall)
	f.updateNeighboursWalls(f.Position, state.Wall)
}

func (f *FloodFill) updateWallsIfNeeded(pos Position, wall ma.Wall) {
	if !validPosition(pos.x, pos.y) {
		return
	}
	f.cells[pos.x][pos.y] = wall
}

func (f *FloodFill) updateNeighboursWalls(pos Position, wall ma.Wall) {
	x, y := pos.x, pos.y
	if wall.Contains(ma.L) {
		f.updateWallsIfNeeded(Position{x, y - 1}, ma.R)
	}
	if wall.Contains(ma.U) {
		f.updateWallsIfNeeded(Position{x + 1, y}, ma.D)
	}
	if wall.Contains(ma.R) {
		f.updateWallsIfNeeded(Position{x, y + 1}, ma.L)
	}
	if wall.Contains(ma.D) {
		f.updateWallsIfNeeded(Position{x - 1, y}, ma.U)
	}
}

func (f *FloodFill) floodFill() {
	st := stack.Stack[Position]{}
	st.Push(Position{f.x, f.y})
	for !st.Empty() {
		topPos := st.Pop()
		minPos := f.getMinOpenNeighbour(topPos)

		if f.getFlood(topPos)-1 == f.getFlood(minPos) {
			continue
		}

		f.setFlood(topPos, f.getFlood(minPos)+1)
		for _, n := range f.getOpenNeighboursNotFinish(topPos) {
			st.Push(n)
		}
	}
}

func (f *FloodFill) setFlood(pos Position, val int) {
	f.flood[pos.x][pos.y] = val
}

func (f *FloodFill) getFlood(pos Position) int {
	return f.flood[pos.x][pos.y]
}

func (f *FloodFill) getCell(pos Position) ma.Wall {
	return f.cells[pos.x][pos.y]
}

func (f *FloodFill) getMinOpenNeighbour(pos Position) (res Position) {
	x, y := pos.x, pos.y
	mn := math.MaxInt
	for _, n := range getNeighboursNotFinish(x, y) {
		if !f.isOpen(x, y, n.x, n.y) {
			continue
		}
		if f.flood[n.x][n.y] < mn {
			mn = f.flood[n.x][n.y]
			res = n
		}
	}
	if mn == math.MaxInt {
		panic("look like no neighbours")
	}
	return res
}

func (f *FloodFill) getOpenNeighboursNotFinish(pos Position) (res []Position) {
	for _, n := range getNeighboursNotFinish(pos.x, pos.y) {
		if !f.isOpen(pos.x, pos.y, n.x, n.y) {
			continue
		}
		res = append(res, n)
	}
	return res
}

func (f *FloodFill) isOpen(x1, y1, x2, y2 int) bool {
	if abs(x1-x2)+abs(y1-y2) == 0 || abs(x1-x2)+abs(y1-y2) > 1 {
		panic("diagonal move or not neighbour")
	}
	switch {
	case x1 < x2:
		return f.cells[x1][y1]&ma.U == 0
	case x1 > x2:
		return f.cells[x1][y1]&ma.D == 0
	case y1 > y2:
		return f.cells[x1][y1]&ma.L == 0
	case y1 < y2:
		return f.cells[x1][y1]&ma.R == 0
	default:
		panic("wtf??")
	}
}

func (f *FloodFill) printFlood() {
	for i := Height - 1; i >= 0; i-- {
		for j := 0; j < Width; j++ {
			fmt.Print(f.flood[i][j], " ")
		}
		fmt.Println()
	}
	fmt.Println("---------")
}
