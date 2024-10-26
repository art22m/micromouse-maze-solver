package solver

import "jackson/internal/maze"

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

func getNeighboursWithDirection(x, y int) (res []PositionWithDirection) {
	if checkPositionNotFinish(x-1, y) {
		res = append(res, PositionWithDirection{Position{x - 1, y}, maze.Down})
	}
	if checkPositionNotFinish(x+1, y) {
		res = append(res, PositionWithDirection{Position{x + 1, y}, maze.Up})
	}
	if checkPositionNotFinish(x, y-1) {
		res = append(res, PositionWithDirection{Position{x, y - 1}, maze.Left})
	}
	if checkPositionNotFinish(x, y+1) {
		res = append(res, PositionWithDirection{Position{x, y + 1}, maze.Right})
	}
	return res
}

func getNeighboursNotFinish(x, y int) (res []Position) {
	if checkPositionNotFinish(x-1, y) {
		res = append(res, Position{x - 1, y})
	}
	if checkPositionNotFinish(x+1, y) {
		res = append(res, Position{x + 1, y})
	}
	if checkPositionNotFinish(x, y-1) {
		res = append(res, Position{x, y - 1})
	}
	if checkPositionNotFinish(x, y+1) {
		res = append(res, Position{x, y + 1})
	}
	return res
}

func validPosition(x, y int) bool {
	return 0 <= x && x < Height && 0 <= y && y <= Width
}

func checkPositionNotFinish(x, y int) bool {
	if FinishXFrom <= x && x <= FinishXTo && FinishYFrom <= y && y <= FinishYTo {
		return false
	}
	return validPosition(x, y)
}
