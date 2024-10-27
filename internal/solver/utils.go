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

func getNeighboursWithDirection(pos Position) (res []PositionWithDirection) {
	if validPosition(pos.Shift(maze.Down)) {
		res = append(res, PositionWithDirection{pos.Shift(maze.Down), maze.Down})
	}
	if validPosition(pos.Shift(maze.Up)) {
		res = append(res, PositionWithDirection{pos.Shift(maze.Up), maze.Up})
	}
	if validPosition(pos.Shift(maze.Left)) {
		res = append(res, PositionWithDirection{pos.Shift(maze.Left), maze.Left})
	}
	if validPosition(pos.Shift(maze.Right)) {
		res = append(res, PositionWithDirection{pos.Shift(maze.Right), maze.Right})
	}
	return res
}

func getNeighboursNotFinish(pos Position) (res []Position) {
	if checkPositionNotFinish(pos.Shift(maze.Down)) {
		res = append(res, pos.Shift(maze.Down))
	}
	if checkPositionNotFinish(pos.Shift(maze.Up)) {
		res = append(res, pos.Shift(maze.Up))
	}
	if checkPositionNotFinish(pos.Shift(maze.Left)) {
		res = append(res, pos.Shift(maze.Left))
	}
	if checkPositionNotFinish(pos.Shift(maze.Right)) {
		res = append(res, pos.Shift(maze.Right))
	}
	return res
}

func validPosition(pos Position) bool {
	return 0 <= pos.x && pos.x < height && 0 <= pos.y && pos.y < width
}

func checkPositionNotFinish(pos Position) bool {
	if finishXFrom <= pos.x && pos.x <= finishXTo && finishYFrom <= pos.y && pos.y <= finishYTo {
		return false
	}
	return validPosition(pos)
}
