package mover

import (
	"log"
	"math"
	"time"

	"jackson/internal/maze"

	"github.com/sirupsen/logrus"
)

type smartMoverState int

const (
	stateDefault smartMoverState = iota + 1
	stateRotatingRight
	stateRotatingLeft
	stateRotating180
)

type SmartMover struct {
	config           SmartMoverConfig
	zeroAngle        int
	state            smartMoverState
	cellState        *CellResp
	cellStateValid   bool
	lastMovementTime time.Time
	facingAngle      int

	baseMover
}

func NewSmartMover(logger *logrus.Entry, sensorsIP, motorsIP string, id string) *SmartMover {
	log.SetPrefix("smart-mover: ")
	sm := &SmartMover{
		zeroAngle:        0,
		state:            stateDefault,
		cellState:        nil,
		cellStateValid:   false,
		lastMovementTime: time.Now(),
		facingAngle:      0,
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
			logger:    logger,
		},
	}

	return sm
}

func (sm *SmartMover) Left() {
	if sm.state == stateDefault {
		sm.state = stateRotatingLeft
		return
	}

	sm.logger.Errorf("Left is called when current state is %d", sm.state)

	switch sm.state {
	case stateRotatingLeft:
		sm.state = stateRotating180
	case stateRotating180:
		sm.state = stateRotatingRight
	case stateRotatingRight:
		sm.state = stateDefault
	}
}

func (sm *SmartMover) Right() {
	if sm.state == stateDefault {
		sm.state = stateRotatingRight
		return
	}

	sm.logger.Errorf("Right is called when current state is %d", sm.state)

	switch sm.state {
	case stateRotatingRight:
		sm.state = stateRotating180
	case stateRotating180:
		sm.state = stateRotatingLeft
	case stateRotatingLeft:
		sm.state = stateDefault
	}
}

// 180 turn
func (sm *SmartMover) Rotate() {
	if sm.state == stateDefault {
		sm.state = stateRotating180
		return
	}

	sm.logger.Errorf("Rotate is called when current state is %d", sm.state)

	switch sm.state {
	case stateDefault:
		sm.state = stateRotating180
	case stateRotatingRight:
		sm.state = stateRotatingLeft
	case stateRotating180:
		sm.state = stateDefault
	case stateRotatingLeft:
		sm.state = stateRotatingRight
	}
}

func (sm *SmartMover) Forward(cells int) {
	for _ = range cells {
		sm.forward(sm.config.robot.CellSize)
	}
}

func (sm *SmartMover) Backward(cells int) {
	for _ = range cells {
		sm.backward(sm.config.robot.CellSize)
	}
}

func (sm *SmartMover) forward(distance int) {
	if sm.state == stateDefault {
		sm.move("forward", sm.getCalibrated(distance, sm.config.calibration.ForwardRatio))
		sm.onMovement()
		sm.fixFacingDirection()
		return
	}

	targetX := 90
	targetY := -90
	if sm.state == stateRotating180 {
		sm.facingAngle = (sm.facingAngle + 180) % 360
	}
	if sm.state == stateRotatingRight {
		targetX = 270
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 90) % 360
	}
	if sm.state == stateRotatingLeft {
		targetX = -270
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 270) % 360
	}

	angle, dist := sm.getFacingCorrectionAndDistance(targetX, targetY)
	if angle > 0 {
		sm.move("right", sm.getCalibrated(angle, sm.config.calibration.TurnRightRatio))
	} else {
		sm.move("left", sm.getCalibrated(angle, sm.config.calibration.TurnLeftRatio))
	}

	sm.fixFacingDirection()

	sm.state = stateDefault
	sm.forward(dist)
}

func (sm *SmartMover) backward(distance int) {
	if sm.state == stateDefault {
		sm.move("backward", sm.getCalibrated(distance, sm.config.calibration.BackwardRatio))
		sm.onMovement()
		sm.fixFacingDirection()
		return
	}

	targetX := 90
	targetY := -90

	if sm.state == stateRotating180 {
		sm.facingAngle = (sm.facingAngle + 180) % 360
	}
	if sm.state == stateRotatingLeft {
		targetX = 270
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 270) % 360
	}
	if sm.state == stateRotatingRight {
		targetX = -270
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 90) % 360
	}

	sm.getFacingCorrectionAndDistance(targetX, targetY)

	angle, dist := sm.getFacingCorrectionAndDistance(targetX, targetY)
	angle += 180
	if angle > 180 {
		angle = -360 + angle
	}
	if angle > 0 {
		sm.move("right", sm.getCalibrated(angle, sm.config.calibration.TurnRightRatio))
	} else {
		sm.move("left", sm.getCalibrated(angle, sm.config.calibration.TurnLeftRatio))
	}

	sm.fixFacingDirection()

	sm.state = stateDefault
	sm.backward(dist)
}

func (sm *SmartMover) CellState(d maze.Direction) Cell {
	if sm.state != stateDefault {
		sm.logger.Fatalf("attempt to get CellState while sm.state = %d", sm.state)
	}

	sm.refreshCellState()

	cell := sm.cellState.ToCell(d)
	sm.logger.Infof("wall: %s, dir: %s", cell, d)
	return cell
}

func (sm *SmartMover) getFacingCorrectionAndDistance(targetX, targetY int) (int, int) {
	return calcVector(sm.cellState.ToCell(maze.Up).Wall, &CellResp{
		Laser: sm.cellState.Laser,
		Imu: struct {
			Roll  float64 "json:\"roll\""
			Pitch float64 "json:\"pitch\""
			Yaw   float64 "json:\"yaw\""
		}{0, 0, float64(int(sm.cellState.Imu.Yaw-float64(sm.zeroAngle)+360) % 360)},
	}, targetX, targetY)
}

func (sm *SmartMover) onMovement() {
	sm.lastMovementTime = time.Now()
	sm.cellStateValid = false
}

func (sm *SmartMover) fixFacingDirection() {
	sm.refreshCellState()

	currentRotation := (int(sm.cellState.Imu.Yaw) + 360 - sm.zeroAngle) % 360
	diff := (sm.facingAngle - currentRotation + 360) % 360
	absDiff := int(math.Abs(float64(diff)))
	if absDiff > sm.config.robot.MaxDerivationFromAxis {
		if diff >= sm.config.calibration.MinTurn {
			if diff > 0 {
				sm.move("right", diff)
			} else {
				sm.move("left", -diff)
			}
		} else {
			if diff > 0 {
				sm.move("left", sm.config.calibration.MinTurn)
				sm.move("right", sm.config.calibration.MinTurn+diff)
			} else {
				sm.move("right", sm.config.calibration.MinTurn)
				sm.move("left", sm.config.calibration.MinTurn-diff)
			}
		}

		sm.fixFacingDirection()
	}
}

func (sm *SmartMover) refreshCellState() {
	if sm.cellStateValid {
		return
	}

	if time.Since(sm.lastMovementTime) < sm.config.robot.SensorsDelay {
		time.Sleep(sm.config.robot.SensorsDelay - time.Since(sm.lastMovementTime))
	}

	sm.cellState = sm.getSensor()
}

func (sm *SmartMover) getCalibrated(value int, calibration float32) int {
	return int(math.Round(float64(value) * float64(calibration)))
}
