package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"fmt"
	"strings"
)

func main() {
	ps := []geom.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: 2},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 1, Y: 2},
	}
	b := board.NewBoard(ps)

	c := board.FillPoints(b, []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 0}}, "A")
	d := board.FillPoints(b, []geom.Point{{X: 0, Y: 1}, {X: 1, Y: 1}}, "B")
	e := board.FillPoints(b, []geom.Point{{X: 0, Y: 0}, {X: 0, Y: 1}}, "C")

	var boardAt = func(p geom.Point) string {
		if label := board.LabelAt(b, p); label != nil {
			return *label
		}
		return "-"
	}

	var dumpBoard = func() {
		fmt.Println()
		fmt.Printf("%+v\n", strings.Join(geom.Stringify(ps, boardAt), "\n"))
		fmt.Println()
	}

	dumpBoard()

	fmt.Printf("%+v\n", board.RemainingPoints(b))

	if c != nil {
		(*c)()
	}

	dumpBoard()

	if d != nil {
		(*d)()
	}

	dumpBoard()

	if e == nil {
		fmt.Println("e is nil")
	} else {
		(*e)()
	}
}
