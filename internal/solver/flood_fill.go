package solver

import (
	"fmt"
	"log"
	"math"
	"sort"

	ma "jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/stack"
)

type FloodFill struct {
	flood [][]int
	cells [][]ma.Wall

	dir ma.Direction
	pos Position

	mo mo.Mover
}

func NewFloodFill(dir ma.Direction, pos Position, mover mo.Mover) *FloodFill {
	flood := make([][]int, height)
	cells := make([][]ma.Wall, height)
	for i := 0; i < height; i++ {
		flood[i] = make([]int, width)
		cells[i] = make([]ma.Wall, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			flood[i][j] = abs(i-getNearest(i, finishXFrom, finishXTo)) + abs(j-getNearest(j, finishYFrom, finishYTo))
			cells[i][j] = 0
		}
	}

	return &FloodFill{
		flood: flood,
		cells: cells,
		mo:    mover,
		pos:   pos,
		dir:   dir,
	}
}

func (f *FloodFill) Solve() {
	it := 0
	for {
		if f.getFlood(f.pos) == 0 {
			log.Println("reached finish")
			break
		}

		it++
		log.Printf("iteration #%d", it)

		f.updateWalls()
		f.floodFill()
		f.move()
	}
}

func (f *FloodFill) move() {
	nextPos := f.getNextPosition()
	log.Printf("next pos=%v\n", nextPos.y)

	newDir := f.rotateIfNeeded(nextPos)
	log.Printf("prev dir=%v, new dir=%v\n", f.dir.String(), newDir.String())
	f.dir = newDir

	f.mo.Forward(1)
	log.Printf("prev pos=%v, new pos=%v\n", f.pos.String(), nextPos.Position.String())
	f.pos = nextPos.Position
}

func (f *FloodFill) rotateIfNeeded(nextPos PositionWithDirection) ma.Direction {
	switch {
	case f.dir.TurnsCount(nextPos.Direction) == 0:
		fmt.Println("no rotate")
		return nextPos.Direction
	case f.dir.TurnsCount(nextPos.Direction) == 2:
		fmt.Println("rotate 180")
		f.mo.Rotate()
		return nextPos.Direction
	default:
		switch f.dir {
		case ma.Left:
			if nextPos.Direction == ma.Up {
				f.mo.Right()
			} else {
				f.mo.Left()
			}
			return nextPos.Direction
		case ma.Right:
			if nextPos.Direction == ma.Up {
				f.mo.Left()
			} else {
				f.mo.Right()
			}
			return nextPos.Direction
		case ma.Down:
			if nextPos.Direction == ma.Left {
				f.mo.Right()
			} else {
				f.mo.Left()
			}
			return nextPos.Direction
		case ma.Up:
			if nextPos.Direction == ma.Left {
				f.mo.Left()
			} else {
				f.mo.Right()
			}
			return nextPos.Direction
		}
		panic("invalid diff turn")
	}
}

func (f *FloodFill) getNextPosition() PositionWithDirection {
	res := make([]PositionWithDirection, 0, 4)
	for _, n := range getNeighboursWithDirection(f.pos) {
		if !f.isOpen(f.pos, n.Position) {
			continue
		}
		res = append(res, n)
	}

	if len(res) == 0 {
		panic("no next position")
	}

	sort.Slice(res, func(i, j int) bool {
		return f.getFlood(res[i].Position) < f.getFlood(res[j].Position) ||
			(f.getFlood(res[i].Position) == f.getFlood(res[j].Position) &&
				f.dir.TurnsCount(res[i].Direction) < f.dir.TurnsCount(res[j].Direction))
	})

	return res[0]
}

func (f *FloodFill) updateWalls() {
	state := f.mo.CellState(f.dir)
	log.Printf("got state: wall=%v\n", state.Wall.String())

	f.updateWallsIfNeeded(f.pos, state.Wall)
	f.updateNeighboursWallsIfNeeded(f.pos, state.Wall)
}

func (f *FloodFill) updateWallsIfNeeded(pos Position, wall ma.Wall) {
	if !validPosition(pos) {
		return
	}
	f.cells[pos.x][pos.y] |= wall
}

func (f *FloodFill) updateNeighboursWallsIfNeeded(pos Position, wall ma.Wall) {
	if wall.Contains(ma.L) {
		f.updateWallsIfNeeded(pos.Shift(ma.Left), ma.R)
	}
	if wall.Contains(ma.U) {
		f.updateWallsIfNeeded(pos.Shift(ma.Up), ma.D)
	}
	if wall.Contains(ma.R) {
		f.updateWallsIfNeeded(pos.Shift(ma.Right), ma.L)
	}
	if wall.Contains(ma.D) {
		f.updateWallsIfNeeded(pos.Shift(ma.Down), ma.U)
	}
}

func (f *FloodFill) floodFill() {
	st := stack.Stack[Position]{}
	st.Push(f.pos)
	for !st.Empty() {
		topPos := st.Pop()
		minPos := f.getMinOpenNeighbourNotFinish(topPos)

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

func (f *FloodFill) getMinOpenNeighbourNotFinish(pos Position) (res Position) {
	mn := math.MaxInt
	for _, n := range getNeighboursNotFinish(pos) {
		if !f.isOpen(pos, n) {
			continue
		}
		if f.getFlood(n) < mn {
			mn = f.getFlood(n)
			res = n
		}
	}
	if mn == math.MaxInt {
		panic("look like no neighbours")
	}
	return res
}

func (f *FloodFill) getOpenNeighboursNotFinish(pos Position) (res []Position) {
	for _, n := range getNeighboursNotFinish(pos) {
		if !f.isOpen(pos, n) {
			continue
		}
		res = append(res, n)
	}
	return res
}

func (f *FloodFill) isOpen(from Position, to Position) bool {
	x1, y1 := from.x, from.y
	x2, y2 := to.x, to.y
	if abs(x1-x2)+abs(y1-y2) == 0 || abs(x1-x2)+abs(y1-y2) > 1 {
		panic("diagonal move or not neighbour")
	}
	switch {
	case x1 < x2:
		return f.getCell(from)&ma.U == 0
	case x1 > x2:
		return f.getCell(from)&ma.D == 0
	case y1 > y2:
		return f.getCell(from)&ma.L == 0
	case y1 < y2:
		return f.getCell(from)&ma.R == 0
	default:
		panic("wtf??")
	}
}

func (f *FloodFill) printFlood() {
	log.Println("----flood-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-3v", f.flood[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}

func (f *FloodFill) printWalls() {
	log.Println("----walls-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-3v", f.cells[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}
