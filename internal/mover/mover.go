package mover

type Cell struct {
	Left, Right, Forward, Backward int
}

type Mover interface {
	Forward(int)
	Backward(int)

	Left()
	Right()

	CellState() Cell
}
