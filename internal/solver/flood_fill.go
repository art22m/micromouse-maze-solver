package solver

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/sirupsen/logrus"

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

	Logger *logrus.Entry
}

type FloodFill struct {
	flood   [][]int
	visited [][]bool
	cells   [][]ma.Wall

	moveForwardOnly bool

	finishFrom Position
	finishTo   Position

	dir ma.Direction
	pos Position

	mo mo.Mover

	iteration int

	logger *logrus.Entry
}

func NewFloodFill(config FloodFillConfig) *FloodFill {
	flood := make([][]int, height)
	cells := make([][]ma.Wall, height)
	visited := make([][]bool, height)
	for i := 0; i < height; i++ {
		flood[i] = make([]int, width)
		cells[i] = make([]ma.Wall, width)
		visited[i] = make([]bool, width)
	}

	return &FloodFill{
		flood:   flood,
		cells:   cells,
		visited: visited,

		mo:              config.Mover,
		pos:             config.StartPosition,
		dir:             config.StartDirection,
		moveForwardOnly: config.MoveForwardOnly,

		finishFrom: Position{finishXFrom, finishYFrom},
		finishTo:   Position{finishXTo, finishYTo},

		logger: config.Logger,
	}
}

func (f *FloodFill) runFastPath(
	visited [][]bool,
	cells [][]ma.Wall,
	pos Position,
	dir ma.Direction,
) {
	f.logger.Println("start fast path")

	f.visited = visited
	f.cells = cells
	f.finishFrom = Position{finishXFrom, finishYFrom}
	f.finishTo = Position{finishXTo, finishYTo}
	f.pos = pos
	f.dir = dir

	path := f.shortestPath()
	fmt.Println("shortest path:")
	for _, p := range path {
		fmt.Print(p.String(), "->")
	}
	fmt.Println("\n--------------")

	if !f.pos.Equal(path[0]) {
		panic("should be equal to current position")
	}

	for i := 1; i < len(path); i++ {
		f.logger.Printf("-------------\nfast path iteration #%d", i)
		if f.isFinish(f.pos) {
			break
		}
		f.move(f.calculateDirection(path[i]))
	}
}

func (f *FloodFill) Solve() {
	f.AskUser()
}

func (f *FloodFill) AskUser() {
	fmt.Println("\n" +
		"(1) \t [Flood Fill] From start, current position won't be changed \n" +
		"(2) \t [Flood Fill] From finish, current position won't be changed \n" +
		"(3) \t [Fast Path] Go from start to finish, current position would be (0,0) and UP direction\n" +
		"(4) \t [Fast Path] Go from start to finish, current position won't be changed\n" +
		"(5) \t Exit",
	)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Println("!!!! scan error", scanner.Err())
	}

	switch scanner.Text() {
	case "1":
		f.startToFinish()
	case "2":
		f.finishToStart()
	case "3":
		f.runFastPath(f.visited, f.cells, Position{x: 0, y: 0}, ma.Up)
	case "4":
		f.runFastPath(f.visited, f.cells, f.pos, f.dir)
	case "5", "exit":
		return
	default:
		fmt.Println("Invalid choice")
	}
	f.AskUser()
}

func (f *FloodFill) startToFinish() {
	f.logger.Info("finding path from start to finish")

	f.flood = make([][]int, height)
	for i := 0; i < height; i++ {
		f.flood[i] = make([]int, width)
	}

	f.dummyFloodFill()
	f.start()
}

func (f *FloodFill) finishToStart() {
	f.logger.Info("finding path from finish to start")

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
		f.logger.Warnf("!!! iteration #%d", f.iteration)

		f.setVisited()
		f.updateWalls()
		if f.isFinish(f.pos) {
			break
		}
		f.smartFloodFill()
		f.move(f.getNextPosition)
	}

	f.logger.Info("finish was reached")
	f.printFlood()
	f.printWalls()
}

func (f *FloodFill) move(getNextPosition func() PositionWithDirection) {
	nextPos := getNextPosition()
	f.logger.Infof("want to go to %v\n", nextPos.String())

	newDir, moveForward := f.rotateIfNeeded(nextPos)
	f.logger.Infof("prev dir=%v, new dir=%v\n", f.dir.String(), newDir.String())

	if moveForward {
		f.mo.Forward(1)
	} else {
		f.mo.Backward(1)
	}

	f.logger.Infof("prev pos=%v, new pos=%v\n", f.pos.String(), nextPos.Position.String())

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
	f.logger.Infof("got state: wall=%v\n", state.Wall.String())
	f.updateWallsIfNeeded(f.pos, state.Wall)
	f.updateNeighboursWallsIfNeeded(f.pos, state.Wall)
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
	st := stack.Stack[Position]{}
	st.Push(f.pos)
	for !st.Empty() {
		topPos := st.Pop()
		minPos := f.getOpenNeighbourWithSmallestFlood(topPos)

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
