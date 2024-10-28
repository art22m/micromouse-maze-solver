package maze

import (
	"fmt"
	"strings"
)

type Direction int

const (
	Left  Direction = 1 << 0
	Up              = 1 << 1
	Right           = 1 << 2
	Down            = 1 << 3
)

func (d Direction) TurnsCount(dc Direction) int {
	if d == dc {
		return 0
	}

	switch d {
	case Left, Right:
		switch dc {
		case Up, Down:
			return 1
		default:
			return 2
		}
	case Up, Down:
		switch dc {
		case Left, Right:
			return 1
		default:
			return 2
		}
	}

	panic("invalid directions")
}

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	}
	return "unknown"
}

type Wall int

const (
	L Wall = 1 << 0
	U      = 1 << 1
	R      = 1 << 2
	D      = 1 << 3

	//LD = L | D
	//RD = R | D
	//UR = U | R
	//UL = U | L
	//LR = L | R
	//UD = U | D
	//
	//LRD = L | R | D
	//URD = U | R | D
	//LUR = L | U | R
	//LUD = L | U | D

	Empty = 0
	Full  = L | U | R | D

	Unknown = 42
)

func (w Wall) Contains(wc Wall) bool {
	return w&wc == wc
}

// Add adds wall
// NOTE: wc should be only L R U D
func (w *Wall) Add(wc Wall) {
	*w |= wc
}

func (w Wall) String() string {
	var sb strings.Builder
	if w.Contains(U) {
		sb.WriteByte('U')
	}
	if w.Contains(R) {
		sb.WriteByte('R')
	}
	if w.Contains(D) {
		sb.WriteByte('D')
	}
	if w.Contains(L) {
		sb.WriteByte('L')
	}
	if sb.Len() == 0 {
		return "x"
	}
	return sb.String()
}

func abs(v int) int {
	if v > 0 {
		return v
	}
	return -v
}

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}

	panic(fmt.Errorf("Direction.Opposite: invalid d (%d)", d))
}

func (d Direction) LocalTo(orientation Direction) Direction {
	if orientation == Up {
		return d
	}

	if d == orientation {
		return Up
	}

	if d.Opposite() == orientation {
		return Down
	}

	if orientation == Right {
		switch d {
		case Up:
			return Left
		case Down:
			return Right
		}
	}

	if orientation == Left {
		switch d {
		case Up:
			return Right
		case Down:
			return Left
		}
	}

	if orientation == Down {
		switch d {
		case Left:
			return Right
		case Right:
			return Left
		}
	}

	panic(fmt.Errorf("LocalTo: invalid combination of d (%d) and orientation (%d)", d, orientation))
}

func (d Direction) GlobalFrom(orientation Direction) Direction {
	if d == orientation || d.Opposite() == orientation {
		return d.LocalTo(orientation).Opposite()
	}

	return d.LocalTo(orientation)
}
