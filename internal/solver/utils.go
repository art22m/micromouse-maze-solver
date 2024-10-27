package solver

import (
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
