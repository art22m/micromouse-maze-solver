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
	_, err := m.move("forward", cell*180)
	if err != nil {
		m.logger.Fatal(err)
	}
}

func (m *DummyMover) Backward(cell int) {
	_, err := m.move("backward", cell*180)
	if err != nil {
		m.logger.Fatal(err)
	}
}

func (m *DummyMover) Left() {
	_, err := m.move("left", 90)
	if err != nil {
		m.logger.Fatal(err)
	}
}

func (m *DummyMover) Right() {
	_, err := m.move("right", 90)
	if err != nil {
		m.logger.Fatal(err)
	}
}

func (m *DummyMover) Rotate() {
	_, err := m.move("right", 180)
	if err != nil {
		m.logger.Fatal(err)
	}
}

func (m *DummyMover) CellState(d maze.Direction) Cell {
	resp, err := m.getSensor()
	if err != nil {
		m.logger.Fatal(err)
	}
	cell := resp.ToCell(d)
	m.logger.Infof("wall: %s, dir: %s", cell, d)
	return cell
}
