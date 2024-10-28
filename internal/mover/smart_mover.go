package mover

import (
	"log"
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
	lastMovementTime time.Time

	baseMover
}

func NewSmartMover(logger *logrus.Entry, sensorsIP, motorsIP string, id string) *SmartMover {
	log.SetPrefix("smart-mover: ")
	sm := &SmartMover{
		zeroAngle:        0,
		state:            stateDefault,
		lastMovementTime: time.Now(),
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
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

}

func (sm *SmartMover) Backward(cells int) {

}

func (sm *SmartMover) CellState(d maze.Direction) Cell {
	resp, err := sm.getSensor()
	if err != nil {
		sm.logger.Fatal(err)
	}
	cell := resp.ToCell(d)
	sm.logger.Infof("wall: %s, dir: %s", cell, d)
	return cell
}
