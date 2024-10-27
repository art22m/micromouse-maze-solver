package solver

import (
	"math"

	ma "jackson/internal/maze"
)

func abs(v int) int {
	if v > 0 {
		return v
	}
	return -v
}

func getNearest(x, from, to int) int {
	switch {
	case x < from:
		return from
	case to < x:
		return to
	default:
		return x
	}
}

func getNeighboursWithDirection(pos Position) (res []PositionWithDirection) {
	if validPosition(pos.Shift(ma.Down)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Down), ma.Down})
	}
	if validPosition(pos.Shift(ma.Up)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Up), ma.Up})
	}
	if validPosition(pos.Shift(ma.Left)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Left), ma.Left})
	}
	if validPosition(pos.Shift(ma.Right)) {
		res = append(res, PositionWithDirection{pos.Shift(ma.Right), ma.Right})
	}
	return res
}

func (f *FloodFill) getNeighboursNotFinish(pos Position) (res []Position) {
	if f.checkPositionNotFinish(pos.Shift(ma.Down)) {
		res = append(res, pos.Shift(ma.Down))
	}
	if f.checkPositionNotFinish(pos.Shift(ma.Up)) {
		res = append(res, pos.Shift(ma.Up))
	}
	if f.checkPositionNotFinish(pos.Shift(ma.Left)) {
		res = append(res, pos.Shift(ma.Left))
	}
	if f.checkPositionNotFinish(pos.Shift(ma.Right)) {
		res = append(res, pos.Shift(ma.Right))
	}
	return res
}

func validPosition(pos Position) bool {
	return 0 <= pos.x && pos.x < height && 0 <= pos.y && pos.y < width
}

func (f *FloodFill) checkPositionNotFinish(pos Position) bool {
	if f.finishXFrom <= pos.x && pos.x <= f.finishXTo && f.finishYFrom <= pos.y && pos.y <= f.finishYTo {
		return false
	}
	return validPosition(pos)
}

func (f *FloodFill) getMinOpenNeighbourNotFinish(pos Position) (res Position) {
	mn := math.MaxInt
	for _, n := range f.getNeighboursNotFinish(pos) {
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
	for _, n := range f.getNeighboursNotFinish(pos) {
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
