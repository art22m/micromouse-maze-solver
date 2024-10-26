package mover

import (
	"fmt"
	"math"
)

type SmartMover struct {
	angle int // from 0 to 360 degrees
	front int
	back  int
	left  int
	right int

	baseMover
}

func NewSmartMover(sensorsIP, motorsIP string, id string) *SmartMover {
	return &SmartMover{
		angle: 0,
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
		},
	}
}

const (
	Front     = 0
	Right     = 90
	Down      = 180
	Left      = 270
	Tolerance = 5 // допустимая погрешность
)

func (m *SmartMover) closestDirectionAndAngle() (string, int) {
	directions := map[string]int{
		"Front": Front,
		"Right": Right,
		"Down":  Down,
		"Left":  Left,
	}

	minDiff := 360.0
	closest := ""
	var angleDiff int

	for dir, angle := range directions {
		diff := angle - m.angle
		if diff > 180 {
			diff -= 360
		} else if diff < -180 {
			diff += 360
		}

		if math.Abs(float64(diff)) < minDiff {
			minDiff = math.Abs(float64(diff))
			closest = dir
			angleDiff = diff
		}
	}

	return closest, angleDiff
}

func (m *SmartMover) isNotAimedAtCenter() bool {
	_, angleDiff := m.closestDirectionAndAngle()
	return math.Abs(float64(angleDiff)) > Tolerance
}

// Метод для центрирования робота к ближайшей оси
func (m *SmartMover) centering() {
	_, angleDiff := m.closestDirectionAndAngle()

	if angleDiff > 0 {
		m.RotateRight(int(angleDiff))
	} else if angleDiff < 0 {
		m.RotateRight(int(-angleDiff))
	} else {
		fmt.Println("Робот уже отцентрован.")
	}
}

func (m *SmartMover) Forward(cell int) {
	if m.isNotAimedAtCenter() {
		m.centering()
	}
	m.move("forward", cell*180)
	// transform cell parameter to mm
	// send command to mouse
	// check position and angle
	// save angle
}

func (m *SmartMover) Backward(cell int) {
	// same as forward
	m.move("backward", cell*180)
}

func (m *SmartMover) RotateLeft(degrees int) {
	_, err := m.move("left", degrees)
	if err != nil {
		return
	}
}

func (m *SmartMover) Left() {
	// get current angle from memory
	m.RotateLeft(90)
}

func (m *SmartMover) RotateRight(degrees int) {
	_, err := m.move("right", degrees)
	if err != nil {
		return
	}
}
func (m *SmartMover) Right() {
	m.RotateRight(90)
}

func (m *SmartMover) CellState() Cell {
	//resp, err := m.getSensor()
	return Cell{}
}

func (m *SmartMover) Rotate() {
	m.RotateRight(180)
}
