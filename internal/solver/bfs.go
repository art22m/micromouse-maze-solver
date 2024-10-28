package solver

import (
	"math"
	"slices"

	"jackson/internal/queue"
)

func (f *FloodFill) shortestPath() (path []Position) {
	f.printVisited()

	dist := make([][]int, height)
	for x := 0; x < height; x++ {
		dist[x] = make([]int, width)
		for y := 0; y < width; y++ {
			dist[x][y] = math.MaxInt
		}
	}
	dist[f.pos.x][f.pos.y] = 0

	parent := map[Position]*Position{}
	parent[f.pos] = nil

	q := queue.Queue[Position]{}
	q.Push(f.pos)

	var finish *Position
	for !q.Empty() {
		curr := q.Pop()
		nb := f.getVisitedOpenNeighbours(curr)
		for _, n := range nb {
			if dist[n.x][n.y] != math.MaxInt {
				continue
			}
			dist[n.x][n.y] = dist[curr.x][curr.y] + 1
			parent[n] = &curr
			q.Push(n)
			if finish == nil && f.isFinish(n) {
				finish = &n
			}
		}
	}

	if finish == nil {
		panic("no finish?")
	}

	curr := *finish
	path = append(path, curr)
	for {
		if parent[curr] == nil {
			break
		}
		path = append(path, *parent[curr])
		curr = *parent[curr]
	}
	slices.Reverse(path)
	return path
}

func (f *FloodFill) getVisitedOpenNeighbours(pos Position) (res []Position) {
	for _, n := range f.getNeighbours(pos) {
		if !f.isOpen(pos, n) {
			continue
		}
		if !f.visited[n.x][n.y] {
			continue
		}
		res = append(res, n)
	}
	return res
}
