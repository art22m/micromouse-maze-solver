package mover

import (
	"log"
	"math"
	"time"

	"jackson/internal/maze"
)

type SmartMover struct {
	angle int // from 0 to 360 degrees
	front int
	back  int
	left  int
	right int

	baseMover
}

const (
	angleUpdateTime = 1
	frontUpdateTime = 1
	backUpdateTime  = 1
)

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
		m.RotateLeft(int(angleDiff))
	} else if angleDiff < 0 {
		m.RotateRight(int(-angleDiff))
	} else {
		log.Println("Робот уже отцентрован.")
	}

	time.Sleep(angleUpdateTime)
	state, _ := m.getSensor()
	m.angle = state.Imu.Yaw
}

func (m *SmartMover) Forward(cell int) {
	if m.isNotAimedAtCenter() {
		m.centering()
	}

	for i := 0; i < cell; i++ {
		time.Sleep(frontUpdateTime)
		dist := m.calcFrontDistance()
		m.move("forward", dist)
		time.Sleep(angleUpdateTime)
		state, _ := m.getSensor()
		m.angle = state.Imu.Yaw
		_, angle := m.closestDirectionAndAngle()
		if angle >= 2 {
			m.RotateRight(angle * 2)
		} else if angle <= -2 {
			m.RotateLeft(int(math.Abs(float64(angle)) * 2))
		}
	}
}

func (m *SmartMover) calcFrontDistance() int {
	frontDiff := 49.0
	state, _ := m.getSensor()
	if state.Laser.Front > 270 {
		return 180
	}
	_, angle := m.closestDirectionAndAngle()
	return int(math.Round(float64(state.Laser.Front) - frontDiff/math.Cos(math.Abs(float64(angle))*(math.Pi/180.0))))
}

func (m *SmartMover) Backward(cell int) {
	if m.isNotAimedAtCenter() {
		m.centering()
	}

	for i := 0; i < cell; i++ {
		time.Sleep(backUpdateTime)
		dist := m.calcBackDistance()
		m.move("backward", dist)
		time.Sleep(angleUpdateTime)
		state, _ := m.getSensor()
		m.angle = state.Imu.Yaw
		_, angle := m.closestDirectionAndAngle()
		if angle >= 2 {
			m.RotateRight(angle * 2)
		} else if angle <= -2 {
			m.RotateLeft(int(math.Abs(float64(angle)) * 2))
		}
	}
}

func (m *SmartMover) calcBackDistance() int {
	backDiff := 49.0
	state, _ := m.getSensor()
	if state.Laser.Back > 270 {
		return 180
	}
	_, angle := m.closestDirectionAndAngle()
	return int(math.Round(float64(state.Laser.Back) - backDiff/math.Cos(math.Abs(float64(angle))*(math.Pi/180.0))))
}

func (m *SmartMover) RotateLeft(degrees int) {
	_, err := m.move("left", degrees)
	if err != nil {
		return
	}
}

func (m *SmartMover) Left() {
	time.Sleep(angleUpdateTime)
	state, _ := m.getSensor()
	m.angle = state.Imu.Yaw

	direction, angleDiff = m.getSensor()
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

func (m *SmartMover) CellState(d maze.Direction) Cell {
	//resp, err := m.getSensor()
	return Cell{}
}

func (m *SmartMover) Rotate() {
	m.RotateRight(180)
}
