package board

import (
	"calendar-puzzle/geom"
)

type PointId int

var directions = [4]geom.Point{
	{X: 0, Y: 1},  // up
	{X: 1, Y: 0},  // right
	{X: 0, Y: -1}, // down
	{X: -1, Y: 0}, // left
}

type Board struct {
	unfilled map[geom.Point]bool
	filled   map[geom.Point]string
	all      []geom.Point
}

func makeSet(ps []geom.Point) map[geom.Point]bool {
	set := make(map[geom.Point]bool, len(ps))
	for _, p := range ps {
		set[p] = true
	}

	return set
}

func LabelAt(b *Board, p geom.Point) *string {
	if _, ok := b.unfilled[p]; ok {
		return nil
	}
	if label, ok := b.filled[p]; ok {
		return &label
	}
	return nil
}

func FillPoints(b *Board, ps *[]geom.Point, offset geom.Point, label string) *func() {
	eps := make([]geom.Point, len(*ps))

	for i, p := range *ps {
		p = geom.AddPoints(p, offset)
		if _, ok := b.unfilled[p]; ok {
			eps[i] = p
		} else {
			return nil
		}
	}

	for _, p := range eps {
		b.filled[p] = label
		delete(b.unfilled, p)
	}

	undo := func() {
		for _, ep := range eps {
			delete(b.filled, ep)
			b.unfilled[ep] = true
		}
	}

	return &undo
}

func RemainingPoints(b *Board) []geom.Point {
	remaining := make([]geom.Point, 0, len(b.unfilled))
	for _, p := range b.all {
		if _, ok := b.unfilled[p]; ok {
			remaining = append(remaining, p)
		}
	}
	return remaining
}

func spreadPoints(b *Board, p geom.Point, limit int, accum map[geom.Point]bool) {
	if len(accum) < limit {
		_, ok := accum[p]
		if !ok {
			_, ok := b.unfilled[p]
			if ok {
				accum[p] = true
				for _, d := range directions {
					spreadPoints(b, geom.AddPoints(p, d), limit, accum)
					if len(accum) >= limit {
						break
					}
				}
			}
		}
	}
}

func CountFill(b *Board, p geom.Point, limit int) int {
	reached := make(map[geom.Point]bool)
	spreadPoints(b, p, limit, reached)
	return len(reached)
}

func NewBoard(points []geom.Point) *Board {
	return &Board{
		unfilled: makeSet(points),
		filled:   make(map[geom.Point]string),
		all:      points,
	}
}
