package solver

import (
	"fmt"
	"log"
	"sort"
	"time"

	ma "jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/queue"
	"jackson/internal/stack"
)

type FloodFillConfig struct {
	StartDirection ma.Direction
	StartPosition  Position

	MoveForwardOnly bool

	Mover mo.Mover
}

type FloodFill struct {
	flood       [][]int
	cachedFlood [][]int
	visited     [][]bool
	cells       [][]ma.Wall

	moveForwardOnly bool

	finishFrom Position
	finishTo   Position

	dir ma.Direction
	pos Position

	mo mo.Mover

	iteration int
}

func NewFloodFill(config FloodFillConfig) *FloodFill {
	log.SetPrefix("flood-fill: ")

	flood := make([][]int, height)
	cells := make([][]ma.Wall, height)
	visited := make([][]bool, height)
	for i := 0; i < height; i++ {
		flood[i] = make([]int, width)
		cells[i] = make([]ma.Wall, width)
		visited[i] = make([]bool, width)
	}

	ff := &FloodFill{
		flood:       flood,
		cachedFlood: nil,
		cells:       cells,
		visited:     visited,

		mo:              config.Mover,
		pos:             config.StartPosition,
		dir:             config.StartDirection,
		moveForwardOnly: config.MoveForwardOnly,

		finishFrom: Position{finishXFrom, finishYFrom},
		finishTo:   Position{finishXTo, finishYTo},
	}
	ff.dummyFloodFill()
	return ff
}

func (f *FloodFill) Solve() ([][]bool, [][]ma.Wall) {
	f.printFlood()
	f.printWalls()

	f.startToFinish()
	time.Sleep(100 * time.Millisecond)
	f.finishToStart()

	return f.visited, f.cells
}

func (f *FloodFill) startToFinish() {
	log.Println("finding path from start to finish")
	f.start()
}

func (f *FloodFill) finishToStart() {
	log.Println("finding path from finish to start")

	f.cachedFlood = f.flood
	f.flood = make([][]int, height)
	for i := 0; i < height; i++ {
		f.flood[i] = make([]int, width)
	}

	f.finishFrom = Position{0, 0}
	f.finishTo = Position{0, 0}

	f.dummyFloodFill()
	f.start()
}

func (f *FloodFill) start() {
	for {
		f.iteration++
		log.Println("-------------")
		log.Printf("iteration #%d", f.iteration)

		f.visited[f.pos.x][f.pos.y] = true
		if f.getFlood(f.pos) == 0 {
			f.updateWalls()
			break
		}

		f.updateWalls()
		f.smartFloodFill()
		f.move()
	}
	log.Println("finish was reached")
	f.printFlood()
	f.printWalls()
}

func (f *FloodFill) move() {
	nextPos := f.getNextPosition()
	log.Printf("want to go to %v\n", nextPos.String())

	newDir, moveForward := f.rotateIfNeeded(nextPos)
	log.Printf("prev dir=%v, new dir=%v\n", f.dir.String(), newDir.String())

	if moveForward {
		f.mo.Forward(1)
	} else {
		f.mo.Backward(1)
	}

	log.Printf("prev pos=%v, new pos=%v\n", f.pos.String(), nextPos.Position.String())

	f.dir = newDir
	f.pos = nextPos.Position
}

func (f *FloodFill) rotateIfNeeded(nextPos PositionWithDirection) (ma.Direction, bool) {
	switch {
	case f.dir.TurnsCount(nextPos.Direction) == 0:
		return f.dir, true
	case f.dir.TurnsCount(nextPos.Direction) == 2:
		if f.moveForwardOnly {
			f.mo.Rotate()
			return nextPos.Direction, true
		}
		return f.dir, false

	default:
		switch f.dir {
		case ma.Left:
			if nextPos.Direction == ma.Up {
				f.mo.Right()
			} else {
				f.mo.Left()
			}
			return nextPos.Direction, true
		case ma.Right:
			if nextPos.Direction == ma.Up {
				f.mo.Left()
			} else {
				f.mo.Right()
			}
			return nextPos.Direction, true
		case ma.Down:
			if nextPos.Direction == ma.Left {
				f.mo.Right()
			} else {
				f.mo.Left()
			}
			return nextPos.Direction, true
		case ma.Up:
			if nextPos.Direction == ma.Left {
				f.mo.Left()
			} else {
				f.mo.Right()
			}
			return nextPos.Direction, true
		}
		panic("invalid diff turn")
	}
}

func (f *FloodFill) getNextPosition() PositionWithDirection {
	res := make([]PositionWithDirection, 0, 4)
	for _, n := range f.getNeighboursWithDirection(f.pos) {
		if !f.isOpen(f.pos, n.Position) {
			continue
		}
		res = append(res, n)
	}

	if len(res) == 0 {
		panic("no next position")
	}

	sort.Slice(res, func(i, j int) bool {
		return f.getFlood(res[i].Position) < f.getFlood(res[j].Position) ||
			(f.getFlood(res[i].Position) == f.getFlood(res[j].Position) &&
				f.dir.TurnsCount(res[i].Direction) < f.dir.TurnsCount(res[j].Direction))
	})

	return res[0]
}

func (f *FloodFill) updateWalls() {
	state := f.mo.CellState(f.dir)
	log.Printf("got state: wall=%v\n", state.Wall.String())

	f.updateWallsIfNeeded(f.pos, state.Wall)
	f.updateNeighboursWallsIfNeeded(f.pos, state.Wall)
}

func (f *FloodFill) updateWallsIfNeeded(pos Position, wall ma.Wall) {
	if !f.validPosition(pos) {
		return
	}
	f.cells[pos.x][pos.y] |= wall
}

func (f *FloodFill) updateNeighboursWallsIfNeeded(pos Position, wall ma.Wall) {
	if wall.Contains(ma.L) {
		f.updateWallsIfNeeded(pos.Shift(ma.Left), ma.R)
	}
	if wall.Contains(ma.U) {
		f.updateWallsIfNeeded(pos.Shift(ma.Up), ma.D)
	}
	if wall.Contains(ma.R) {
		f.updateWallsIfNeeded(pos.Shift(ma.Right), ma.L)
	}
	if wall.Contains(ma.D) {
		f.updateWallsIfNeeded(pos.Shift(ma.Down), ma.U)
	}
}

func (f *FloodFill) dummyFloodFill() {
	visited := make(map[Position]struct{}, height*width)
	q := queue.Queue[Position]{}

	for x := f.finishFrom.x; x <= f.finishTo.x; x++ {
		for y := f.finishFrom.y; y <= f.finishTo.y; y++ {
			pos := Position{x: x, y: y}
			f.setFlood(pos, 0)
			visited[pos] = struct{}{}
			q.Push(pos)
		}
	}

	for !q.Empty() {
		frontPos := q.Pop()
		nb := f.getOpenNeighbours(frontPos)
		for _, n := range nb {
			if _, ok := visited[n]; ok {
				continue
			}
			visited[n] = struct{}{}
			f.setFlood(n, f.getFlood(frontPos)+1)
			q.Push(n)
		}
	}
}

func (f *FloodFill) smartFloodFill() {
	f.printFlood()
	st := stack.Stack[Position]{}
	st.Push(f.pos)
	for !st.Empty() {
		topPos := st.Pop()
		minPos := f.getOpenNeighbourWithSmallestFlood(topPos)

		fmt.Println(topPos.String())
		fmt.Println(minPos.String())

		if f.getFlood(topPos)-1 == f.getFlood(minPos) {
			continue
		}

		f.setFlood(topPos, f.getFlood(minPos)+1)
		for _, n := range f.getOpenNeighbours(topPos) {
			if f.isFinish(n) {
				continue
			}
			st.Push(n)
		}
	}
}

func (f *FloodFill) setFlood(pos Position, val int) {
	f.flood[pos.x][pos.y] = val
}

func (f *FloodFill) getFlood(pos Position) int {
	return f.flood[pos.x][pos.y]
}

func (f *FloodFill) getCell(pos Position) ma.Wall {
	return f.cells[pos.x][pos.y]
}

func (f *FloodFill) printFlood() {
	log.Println("----flood-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-4v", f.flood[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}

func (f *FloodFill) printWalls() {
	log.Println("----walls-----")
	for i := height - 1; i >= 0; i-- {
		for j := 0; j < width; j++ {
			fmt.Printf("%-4v", f.cells[i][j])
		}
		fmt.Println()
	}
	log.Println("-------------")
}
