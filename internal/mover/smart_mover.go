package mover

import (
	"fmt"
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
	sm := &SmartMover{
		config:           LoadSmartMoverConfig(),
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

	sm.refreshCellState()
	sm.zeroAngle = int(sm.cellState.Imu.Yaw)

	return sm
}

func (sm *SmartMover) Reset() {
	sm.zeroAngle = 0
	sm.state = stateDefault
	sm.cellState = nil
	sm.cellStateValid = false
	sm.lastMovementTime = time.Now()
	sm.facingAngle = 0

	sm.refreshCellState()
	sm.zeroAngle = int(sm.cellState.Imu.Yaw)
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
		sm.fixFacingDirection(0)
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
		targetX = -90
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 270) % 360
	}

	angleBefore := int(sm.cellState.Imu.Yaw+360-float64(sm.zeroAngle)) % 360
	angle, dist := sm.getFacingCorrectionAndDistance(targetX, targetY)
	if angle > 0 {
		sm.move("right", sm.getCalibrated(angle, sm.config.calibration.TurnRightRatio))
	} else {
		sm.move("left", sm.getCalibrated(-angle, sm.config.calibration.TurnLeftRatio))
	}
	sm.onMovement()

	angleDiff := (int(angleBefore) + angle + 360) % 90
	if angleDiff > 45 {
		angleDiff -= 90
	}
	sm.fixFacingDirection(angleDiff)

	sm.state = stateDefault
	sm.forward(int(float64(dist) * (float64(distance) / 180.0)))
}

func (sm *SmartMover) backward(distance int) {
	if sm.state == stateDefault {
		sm.move("backward", sm.getCalibrated(distance, sm.config.calibration.BackwardRatio))
		sm.onMovement()
		sm.fixFacingDirection(0)
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
		targetX = -90
		targetY = 90
		sm.facingAngle = (sm.facingAngle + 90) % 360
	}

	angle, dist := sm.getFacingCorrectionAndDistance(targetX, targetY)
	angle += 180
	if angle > 180 {
		angle = -360 + angle
	}
	if angle > 0 {
		sm.move("right", sm.getCalibrated(angle, sm.config.calibration.TurnRightRatio))
	} else {
		sm.move("left", sm.getCalibrated(-angle, sm.config.calibration.TurnLeftRatio))
	}
	sm.onMovement()

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
	sm.refreshCellState()

	return calcVector(sm.cellState.ToCell(maze.Up).Wall, &CellResp{
		Laser: sm.cellState.Laser,
		Imu: struct {
			Roll  float64 "json:\"roll\""
			Pitch float64 "json:\"pitch\""
			Yaw   float64 "json:\"yaw\""
		}{0, 0, float64((int(sm.cellState.Imu.Yaw) + 360 - sm.zeroAngle) % 360)},
	}, targetX, targetY)
}

func (sm *SmartMover) onMovement() {
	sm.lastMovementTime = time.Now()
	sm.cellStateValid = false
}

func (sm *SmartMover) fixFacingDirection(angleDiff int) {
	sm.refreshCellState()

	currentRotation := (int(sm.cellState.Imu.Yaw) + 360 - sm.zeroAngle) % 360
	requiredAngle := (sm.facingAngle + angleDiff + 360) % 360
	diff := (requiredAngle - currentRotation)
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}

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

		sm.onMovement()
		sm.fixFacingDirection(angleDiff)
	}
}

func (sm *SmartMover) refreshCellState() {
	if sm.cellStateValid {
		return
	}

	if time.Since(sm.lastMovementTime) < time.Duration(sm.config.robot.SensorsDelayMs)*time.Millisecond {
		time.Sleep(time.Duration(sm.config.robot.SensorsDelayMs)*time.Millisecond - time.Since(sm.lastMovementTime))
	}

	sm.cellState = sm.getSensor()
	sm.cellStateValid = true
}

func (sm *SmartMover) getCalibrated(value int, calibration float32) int {
	return int(math.Round(float64(value) * float64(calibration)))
}

func (sm *SmartMover) Calibrate() {
	config := SmartMoverCalibrationConfig{}

	sm.config.calibration = SmartMoverCalibrationConfig{}
	sm.config.calibration.MinTurn = 90

	sm.refreshCellState()
	yawBefore := sm.cellState.Imu.Yaw
	sm.move("right", 90)
	sm.onMovement()
	sm.refreshCellState()
	config.TurnRightRatio = float32(int(sm.cellState.Imu.Yaw-yawBefore+360)%360) / 90.0
	fmt.Println("!!! TRR", config.TurnRightRatio)
	sm.config.calibration.TurnRightRatio = config.TurnRightRatio

	sm.refreshCellState()
	yawBefore = sm.cellState.Imu.Yaw
	sm.move("left", 90)
	sm.onMovement()
	sm.refreshCellState()
	config.TurnLeftRatio = float32(int(yawBefore-sm.cellState.Imu.Yaw+360)%360) / 90.0
	fmt.Println("!!! TLL", config.TurnLeftRatio)
	sm.config.calibration.TurnLeftRatio = config.TurnLeftRatio

	sm.refreshCellState()
	yawBefore = sm.cellState.Imu.Yaw
	sm.move("right", 180)
	sm.onMovement()
	sm.refreshCellState()
	config.Turn180Ratio = float32(int(sm.cellState.Imu.Yaw-yawBefore+360)%360) / 180.0

	config.MinTurn = 60
	// for i := 10; i < 100; i++ {
	// 	sm.refreshCellState()
	// 	yawBefore = sm.cellState.Imu.Yaw
	// 	sm.move("right", i)
	// 	sm.onMovement()
	// 	sm.refreshCellState()
	// 	if int(sm.cellState.Imu.Yaw-yawBefore+360)%360 > i*2/3 {
	// 		config.MinTurn = i + 10
	// 		fmt.Println("!!! MT", config.MinTurn)
	// 		break
	// 	}
	// }

	sm.config.calibration.ForwardRatio = 1
	sm.config.calibration.BackwardRatio = 1

	sm.fixFacingDirection(0)
	sm.Right()
	sm.forward(1)
	sm.refreshCellState()
	backBefore := sm.cellState.Laser.Back
	fmt.Println(backBefore)
	sm.forward(120)
	sm.refreshCellState()
	backAfter := sm.cellState.Laser.Back
	fmt.Println(backAfter)
	fmt.Println("!!!", backAfter-backBefore)
	config.ForwardRatio = (float32(backAfter) - float32(backBefore)) / 120.0
	fmt.Println("!!! FWR", config.ForwardRatio)

	sm.fixFacingDirection(0)
	sm.refreshCellState()
	backBefore = sm.cellState.Laser.Back
	sm.backward(100)
	sm.refreshCellState()
	backAfter = sm.cellState.Laser.Back
	config.BackwardRatio = (float32(backBefore) - float32(backAfter)) / 120.0

	saveCalibrationConfig(config)
}
