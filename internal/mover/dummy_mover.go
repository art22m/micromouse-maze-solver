package mover

import (
	"log"

	"jackson/internal/maze"
)

type DummyMover struct {
	baseMover
}

func NewDummyMover(sensorsIP, motorsIP string, id string) *DummyMover {
	log.SetPrefix("dummy-mover: ")
	return &DummyMover{
		baseMover: baseMover{
			motorsIP:  motorsIP,
			sensorsIP: sensorsIP,
			id:        id,
		},
	}
}

func (m *DummyMover) Forward(cell int) {
	_, err := m.move("forward", cell*180)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DummyMover) Backward(cell int) {
	_, err := m.move("backward", cell*180)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DummyMover) Left() {
	_, err := m.move("left", 90)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DummyMover) Right() {
	_, err := m.move("right", 90)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DummyMover) Rotate() {
	_, err := m.move("right", 180)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DummyMover) CellState(d maze.Direction) Cell {
	resp, err := m.getSensor()
	if err != nil {
		log.Fatal(err)
	}
	return resp.ToCell(d)
}
