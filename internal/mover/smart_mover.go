package mover

import (
	"fmt"
	"log"
	"math"
	"time"

	"jackson/internal/maze"
)

type SmartMover struct {
	angle      int // from 0 to 360 degrees
	startAngle int
	front      int
	back       int
	left       int
	right      int

	targetAimAngle int
	state          *CellResp

	rightPlus int
	leftPlus  int

	baseMover
}

const (
	angleUpdateTime = 500 * time.Millisecond
	frontUpdateTime = 500 * time.Millisecond
	backUpdateTime  = 500 * time.Millisecond
	allUpdateTime   = 500 * time.Millisecond
)

func NewSmartMover(sensorsIP, motorsIP string, id string) *SmartMover {
	log.SetPrefix("smart-mover: ")
	sm := &SmartMover{
		angle: 0,
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
		},
	}
	sm.Calibrate()
	time.Sleep(allUpdateTime)
	return sm
}

func (m *SmartMover) Calibrate() {
	m.state, _ = m.getSensor()
	m.startAngle = int(m.state.Imu.Yaw)
	m.angle = m.startAngle

	time.Sleep(allUpdateTime)
	m.RotateRight(120)
	time.Sleep(allUpdateTime)
	m.state, _ = m.getSensor()
	m.angle = int(m.state.Imu.Yaw)
	m.rightPlus = (m.angle-m.startAngle+360)%360 - 120
	log.Println("!!! RIGHT ERROR", m.rightPlus, "DEGREES")

	prevAngle := m.angle
	time.Sleep(allUpdateTime)
	m.RotateLeft(120)
	time.Sleep(allUpdateTime)
	m.state, _ = m.getSensor()
	m.angle = int(m.state.Imu.Yaw)
	m.rightPlus = (prevAngle-m.angle+360)%360 - 120
	log.Println("!!! LEFT ERROR", m.rightPlus, "DEGREES")
}

const (
	Front     = 0
	Right     = 90
	Down      = 180
	Left      = 270
	Tolerance = 5 // допустимая погрешность
)

// вращает робота в сторону от стены если он слишком близко
func (m *SmartMover) fixCenter() int {
	if m.state.Laser.Left > 100 || m.state.Laser.Right > 100 {
		// справа или слева нет стены
		return 0
	}

	if m.state.Laser.Left < 20 {
		return -5
	} else if m.state.Laser.Right < 20 {
		return 5
	}

	return 0
}

func (m *SmartMover) fixAngle() int {
	fmt.Println("!!!", m.angle, m.startAngle, m.targetAimAngle)
	diff := (m.angle+360-m.startAngle)%360 - m.targetAimAngle
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}

	if math.Abs(float64(diff)) < 10 {
		return diff + m.fixCenter()
	}

	return diff
}

func (m *SmartMover) isNotAimedAtCenter() bool {
	angleDiff := m.fixAngle()
	return math.Abs(float64(angleDiff)) > Tolerance
}

// Метод для центрирования робота к нужной оси
func (m *SmartMover) centering() {
	angleDiff := m.fixAngle()
	if angleDiff > 0 {
		m.RotateLeft(int(angleDiff))
	} else if angleDiff < 0 {
		m.RotateRight(int(-angleDiff))
	} else {
		log.Println("Робот уже отцентрован.")
	}

	m.updateAngle()
}

func (m *SmartMover) updateAngle() {
	time.Sleep(angleUpdateTime)
	m.state, _ = m.getSensor()
	m.angle = int(m.state.Imu.Yaw)
}

func (m *SmartMover) Forward(cell int) {
	for i := 0; i < cell; i++ {
		m.updateAngle()
		for m.isNotAimedAtCenter() {
			m.centering()
		}

		m.move("forward", 170)

		dist := m.calcFrontDistance()
		fmt.Println("!!! dist:", dist)
		if dist != -1 {
			if dist < -5 {
				m.move("backward", -dist)
			} else if dist > 5 {
				m.move("forward", dist)
			}
		}
		//m.updateAngle()
		//_, angle := m.fixAngle()
		//if angle >= 2 {
		//	m.RotateRight(angle * 2)
		//} else if angle <= -2 {
		//	m.RotateLeft(int(math.Abs(float64(angle)) * 2))
		//}
	}
}

func (m *SmartMover) calcFrontDistance() int {
	frontDiff := 49.0
	m.state, _ = m.getSensor()
	if m.state.Laser.Front > 100 {
		return -1
	}
	angle := m.fixAngle()
	return int(math.Round(float64(m.state.Laser.Front) - frontDiff/math.Cos(math.Abs(float64(angle))*(math.Pi/180.0))))
}

func (m *SmartMover) Backward(cell int) {
	for m.isNotAimedAtCenter() {
		m.centering()
	}

	for i := 0; i < cell; i++ {
		time.Sleep(backUpdateTime)
		dist := m.calcBackDistance()
		m.move("backward", dist)
		time.Sleep(angleUpdateTime)
		m.state, _ = m.getSensor()
		m.angle = int(m.state.Imu.Yaw)
		angle := m.fixAngle()
		if angle >= 2 {
			m.RotateLeft(angle)
		} else if angle <= -2 {
			m.RotateRight(int(math.Abs(float64(angle))))
		}
	}
}

func (m *SmartMover) calcBackDistance() int {
	backDiff := 49.0
	m.state, _ = m.getSensor()
	if m.state.Laser.Back > 270 {
		return 180
	}
	angle := m.fixAngle()
	return int(math.Round(float64(m.state.Laser.Back) - backDiff/math.Cos(math.Abs(float64(angle))*(math.Pi/180.0))))
}

func (m *SmartMover) RotateLeft(degrees int) {
	degrees += m.leftPlus

	if degrees < 0 {
		degrees -= m.leftPlus
	}

	if degrees < 40 {
		_, err := m.move("right", 120)
		if err != nil {
			return
		}

		_, err = m.move("left", 120+degrees)
		if err != nil {
			return
		}
		return
	}

	_, err := m.move("left", degrees)
	if err != nil {
		return
	}
}

func (m *SmartMover) Left() {
	m.updateAngle()

	angleDiff := m.fixAngle()
	m.RotateLeft(90 + angleDiff)
	m.targetAimAngle = (m.targetAimAngle + 270) % 360
}

func (m *SmartMover) RotateRight(degrees int) {
	degrees += m.rightPlus

	if degrees < 0 {
		degrees -= m.rightPlus
	}

	if degrees < 40 {
		_, err := m.move("left", 120)
		if err != nil {
			return
		}

		_, err = m.move("right", 120+degrees)
		if err != nil {
			return
		}
		return
	}

	_, err := m.move("right", degrees)
	if err != nil {
		return
	}
}
func (m *SmartMover) Right() {
	m.updateAngle()

	angleDiff := m.fixAngle()
	m.RotateRight(90 - angleDiff)
	m.targetAimAngle = (m.targetAimAngle + 90) % 360
}

func (m *SmartMover) Rotate() {
	m.updateAngle()

	angleDiff := m.fixAngle()
	m.RotateRight(180 - angleDiff)
	m.targetAimAngle = (m.targetAimAngle + 180) % 360
}

func (m *SmartMover) CellState(d maze.Direction) Cell {
	time.Sleep(allUpdateTime)
	var err error

	m.centering()

	time.Sleep(allUpdateTime)

	m.state, err = m.getSensor()
	if err != nil {
		log.Fatal(err)
	}
	return m.state.ToCell(d)
}
