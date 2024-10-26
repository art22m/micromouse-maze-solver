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
	L Wall = 1 << 0
	U      = 1 << 1
	R      = 1 << 2
	D      = 1 << 3

	LD = L | D
	RD = R | D
	UR = U | R
	UL = U | L
	LR = L | R
	UD = U | D

	LRD = L | R | D
	URD = U | R | D
	LUR = L | U | R
	LUD = L | U | D

	Empty = 0
	Full  = L | U | R | D

	Unknown = 42
)

func (w Wall) Contains(wc Wall) bool {
	return w&wc == wc
}
