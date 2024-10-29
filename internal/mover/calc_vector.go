package mover

import (
	"math"

	"jackson/internal/maze"
)

func fromDegToRad(angle float64) float64 {
	return (angle * math.Pi) / 180.0
}

func calcVector(walls maze.Wall, state *CellResp, targetX int, targetY int) (rotate int, forward int) {
	// cell size
	cellSize := 168.0

	// sensor shift inside mouse
	fromCenterToFrontSensor, fromCenterToBackSensor := 20.0, 20.0

	// side sensor shift from front
	fromFrontSideSensorsShift := 20.0

	// mouse size
	mouseLen := 70.0
	mouseWidth := 60.0

	// transform angle to relative base
	angle := float64(int(state.Imu.Yaw) % 90)

	/*
		                           | /
		if angle righter then axis |/  then sideFromAxis = 1

								 \ |
		if angle lefter then axis \| then sideFromAxis = -1
	*/
	sideFromAxis := 1.0
	if angle > 45 {
		angle = 90 - angle
		sideFromAxis = -1.0
	}
	angle = fromDegToRad(angle)

	diagonalX, diagonalY := -1.0, -1.0

	// xDiagonal calculation
	if walls.Contains(maze.L) {
		diagonalX =
			(state.Laser.Left * math.Cos(angle)) -
				sideFromAxis*((mouseLen/2)-fromFrontSideSensorsShift)*math.Sin(angle) +
				(mouseWidth/2)/math.Cos(angle)
	} else if walls.Contains(maze.R) {
		diagonalX = cellSize -
			((state.Laser.Right * math.Cos(angle)) +
				sideFromAxis*((mouseLen/2)-fromFrontSideSensorsShift)*math.Sin(angle) +
				(mouseWidth/2)/math.Cos(angle))
	}

	// yDiagonal calculation
	if walls.Contains(maze.D) {
		diagonalY = (state.Laser.Back + mouseLen/2 - fromCenterToBackSensor) * math.Cos(angle)
	} else if walls.Contains(maze.U) {
		diagonalY = cellSize - (state.Laser.Back+mouseLen/2-fromCenterToFrontSensor)*math.Cos(angle)
	}

	// calc distance as sqrt( (x_d - x_t)^2 + (y_d - y_t)^2 )
	distance := math.Sqrt((diagonalX-float64(targetX))*(diagonalX-float64(targetX)) + (diagonalY-float64(targetY))*(diagonalY-float64(targetY)))

	// angle
	theta := math.Atan2(float64(targetY)-diagonalY, float64(targetX)-diagonalX)
	phi := math.Pi/2.0 - theta - sideFromAxis*angle
	return int(phi * (180.0 / math.Pi)), int(distance)
}
