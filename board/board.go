package board

import (
	"go-calendar-puzzle/geom"
)

type PointId int

var directions = [4]geom.Point{
	{X: 0, Y: 1},  // up
	{X: 1, Y: 0},  // right
	{X: 0, Y: -1}, // down
	{X: -1, Y: 0}, // left
}

type Board struct {
	unfilled map[PointId]bool
	filled   map[PointId]string
	all      []geom.Point
	adder    func(geom.Point, geom.Point) geom.Point
	encoder  func(geom.Point) PointId
}

func AddPoints(a, b geom.Point) geom.Point {
	return geom.Point{X: a.X + b.X, Y: a.Y + b.Y}
}

const X_LIMIT = 16

func encodePoint(p geom.Point) PointId {
	return PointId(p.X*X_LIMIT + p.Y)
}

func makeSet(ps []geom.Point, enc func(geom.Point) PointId) map[PointId]bool {
	set := make(map[PointId]bool, len(ps))
	for _, p := range ps {
		set[enc(p)] = true
	}

	return set
}

func LabelAt(b *Board, p geom.Point) *string {
	id := b.encoder(p)
	if _, ok := b.unfilled[id]; ok {
		return nil
	}
	if label, ok := b.filled[id]; ok {
		return &label
	}
	return nil
}

func FillPoints(b *Board, ps []geom.Point, label string) *func() {
	eps := make([]PointId, len(ps))

	for i, p := range ps {
		id := b.encoder(p)
		if _, ok := b.unfilled[id]; ok {
			eps[i] = id
		} else {
			return nil
		}
	}

	for _, id := range eps {
		b.filled[id] = label
		delete(b.unfilled, id)
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
	for id := range b.unfilled {
		for _, p := range b.all {
			if b.encoder(p) == id {
				remaining = append(remaining, p)
				break
			}
		}
	}
	return remaining
}

func spreadPoints(b *Board, p geom.Point, limit int, accum map[PointId]bool) {
	ep := b.encoder(p)
	if len(accum) < limit {
		_, ok := b.unfilled[ep]
		if ok {
			_, ok := accum[ep]
			if !ok {
				accum[ep] = true
				for _, d := range directions {
					spreadPoints(b, b.adder(p, d), limit, accum)
					if len(accum) >= limit {
						break
					}
				}
			}
		}
	}
}

func CountFill(b *Board, p geom.Point, limit int) int {
	reached := make(map[PointId]bool)
	spreadPoints(b, p, limit, reached)
	return len(reached)
}

func NewBoard(points []geom.Point) *Board {
	return &Board{
		unfilled: makeSet(points, encodePoint),
		filled:   make(map[PointId]string),
		all:      points,
		adder:    AddPoints,
		encoder:  encodePoint,
	}
}
