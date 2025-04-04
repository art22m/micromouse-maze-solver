package mover

import (
	"github.com/sirupsen/logrus"

	"jackson/internal/maze"
)

type DummyMover struct {
	baseMover
}

func NewDummyMover(logger *logrus.Entry, sensorsIP, motorsIP string, id string) *DummyMover {
	return &DummyMover{
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
			logger:    logger,
		},
	}
}

func (m *DummyMover) Forward(cell int) {
	m.move("forward", cell*180)
}

func (m *DummyMover) Backward(cell int) {
	m.move("backward", cell*180)
}

func (m *DummyMover) Left() {
	m.move("left", 90)
}

func (m *DummyMover) Right() {
	m.move("right", 90)
}

func (m *DummyMover) Rotate() {
	m.move("right", 180)
}

func (m *DummyMover) CellState(d maze.Direction) Cell {
	resp := m.getSensor()
	cell := resp.ToCell(d)
	m.logger.Infof("wall: %s, dir: %s", cell, d)
	return cell
}

func (m *DummyMover) Reset() {

}
