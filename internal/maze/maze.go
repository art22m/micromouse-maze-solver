package maze

type Direction int

const (
	Up Direction = iota + 1
	Right
	Down
	Left
)

func (d Direction) TurnsCount(dc Direction) int {
	if d == dc {
		return 0
	}
	if abs(int(d-dc)) == 1 || abs(int(d-dc)) == 3 {
		return 1
	}
	return 2
}

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

// Add adds wall
// NOTE: wc should be only L R U D
func (w Wall) Add(wc Wall) {
	w &= wc
}

func abs(v int) int {
	if v > 0 {
		return v
	}
	return -v
}
