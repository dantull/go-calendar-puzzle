package board

import "fmt"

type PointId int

var directions = [4]Point{
	{0, 1},  // up
	{1, 0},  // right
	{0, -1}, // down
	{-1, 0}, // left
}

type Point struct {
	X int
	Y int
}

type Board struct {
	unfilled map[PointId]bool
	filled   map[PointId]string
	all      []Point
	adder    func(Point, Point) Point
	encoder  func(Point) PointId
}

func AddPoints(a, b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

const X_LIMIT = 16

func encodePoint(p Point) PointId {
	return PointId(p.X*X_LIMIT + p.Y)
}

func makeSet(ps []Point, enc func(Point) PointId) map[PointId]bool {
	set := make(map[PointId]bool, len(ps))
	for _, p := range ps {
		set[enc(p)] = true
	}

	return set
}

func LabelAt(b *Board, p Point) *string {
	id := b.encoder(p)
	if _, ok := b.unfilled[id]; ok {
		return nil
	}
	if label, ok := b.filled[id]; ok {
		return &label
	}
	return nil
}

func FillPoints(b *Board, ps []Point, label string) *func() {
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

func RemainingPoints(b *Board) []Point {
	remaining := make([]Point, 0, len(b.unfilled))
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

func spreadPoints(b *Board, p Point, limit int, accum map[PointId]bool) {
	ep := b.encoder(p)
	if len(accum) < limit {
		_, ok := b.unfilled[ep]
		if ok {
			_, ok := accum[ep]
			if !ok {
				accum[ep] = true
				for _, d := range directions {
					fmt.Println("test dir: ", d)
					spreadPoints(b, b.adder(p, d), limit, accum)
					if len(accum) >= limit {
						break
					}
				}
			}
		}
	}
}

func CountFill(b *Board, p Point, limit int) int {
	reached := make(map[PointId]bool)
	spreadPoints(b, p, limit, reached)
	return len(reached)
}

func NewBoard(points []Point) *Board {
	return &Board{
		unfilled: makeSet(points, encodePoint),
		filled:   make(map[PointId]string),
		all:      points,
		adder:    AddPoints,
		encoder:  encodePoint,
	}
}
