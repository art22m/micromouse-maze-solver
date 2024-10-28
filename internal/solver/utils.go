package solver

import (
	"fmt"
	"log"
	"sort"

	ma "jackson/internal/maze"
)

func abs(v int) int {
	if v > 0 {
		return v
	}
	return -v
}

func (f *FloodFill) validPosition(pos Position) bool {
	return 0 <= pos.x && pos.x < height && 0 <= pos.y && pos.y < width
}

func (f *FloodFill) isFinish(pos Position) bool {
	return f.finishFrom.x <= pos.x && pos.x <= f.finishTo.x &&
		f.finishFrom.y <= pos.y && pos.y <= f.finishTo.y
}

func (f *FloodFill) getNeighboursWithDirection(pos Position) (res []PositionWithDirection) {
	if f.validPosition(pos.Shift(ma.Down)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Down), ma.Down})
	}
	if f.validPosition(pos.Shift(ma.Up)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Up), ma.Up})
	}
	if f.validPosition(pos.Shift(ma.Left)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Left), ma.Left})
	}
	if f.validPosition(pos.Shift(ma.Right)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Right), ma.Right})
	}
	return res
}

func (f *FloodFill) getNeighbours(pos Position) (res []Position) {
	if f.validPosition(pos.Shift(ma.Down)) {
		res = append(res, pos.Shift(ma.Down))
	}
	if f.validPosition(pos.Shift(ma.Up)) {
		res = append(res, pos.Shift(ma.Up))
	}
	if f.validPosition(pos.Shift(ma.Left)) {
		res = append(res, pos.Shift(ma.Left))
	}
	if f.validPosition(pos.Shift(ma.Right)) {
		res = append(res, pos.Shift(ma.Right))
	}
	return res
}

func (f *FloodFill) calculateDirection(pos Position) func() PositionWithDirection {
	switch {
	case f.pos.Shift(ma.Down).Equal(pos):
		return func() PositionWithDirection { return PositionWithDirection{pos, ma.Down} }
	case f.pos.Shift(ma.Up).Equal(pos):
		return func() PositionWithDirection { return PositionWithDirection{pos, ma.Up} }
	case f.pos.Shift(ma.Left).Equal(pos):
		return func() PositionWithDirection { return PositionWithDirection{pos, ma.Left} }
	case f.pos.Shift(ma.Right).Equal(pos):
		return func() PositionWithDirection { return PositionWithDirection{pos, ma.Right} }
	default:
		panic("diagonal")
	}
}

func (f *FloodFill) getOpenNeighbourWithSmallestFlood(pos Position) Position {
	ns := f.getOpenNeighbours(pos)
	if len(ns) == 0 {
		panic("no open neighbours")
	}
	sort.Slice(ns, func(i, j int) bool {
		return f.getFlood(ns[i]) < f.getFlood(ns[j])
	})
	return ns[0]
}

func (f *FloodFill) getOpenNeighbours(pos Position) (res []Position) {
	for _, n := range f.getNeighbours(pos) {
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

func (f *FloodFill) updateWallsIfNeeded(pos Position, wall ma.Wall) {
	if !f.validPosition(pos) {
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

func (f *FloodFill) setFlood(pos Position, val int) {
	f.flood[pos.x][pos.y] = val
}

func (f *FloodFill) getFlood(pos Position) int {
	return f.flood[pos.x][pos.y]
}

func (f *FloodFill) getCell(pos Position) ma.Wall {
	return f.cells[pos.x][pos.y]
}

func (f *FloodFill) setVisited() {
	f.visited[f.pos.x][f.pos.y] = true
}

func (f *FloodFill) printFlood() {
	log.Println("----flood-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-4v", f.flood[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}

func (f *FloodFill) printWalls() {
	log.Println("----walls-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-4v", f.cells[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}

func (f *FloodFill) printVisited() {
	log.Println("----visited-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			if f.visited[i][j] {
				fmt.Print("x")
			} else {
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
	log.Println("-------------")
}
