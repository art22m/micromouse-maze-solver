package maze

type Orientation int

const (
	Up Orientation = iota + 1
	Right
	Down
	Left
)

type Wall int

const (
	L Wall = iota + 1
	U
	R
	D

	LD
	RD
	UR
	LU

	LR
	UD
	LRD
	URD

	LUR
	LUD
	Empty
	Full
)
