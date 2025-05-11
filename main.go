package board

import (
	"fmt"
	"go-calendar-puzzle/board"
	"go-calendar-puzzle/geom"
)

func Main() {
	b := board.NewBoard([]geom.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: 2},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 1, Y: 2},
	})

	c := board.FillPoints(b, []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 0}}, "A")
	d := board.FillPoints(b, []geom.Point{{X: 0, Y: 1}, {X: 1, Y: 1}}, "B")
	e := board.FillPoints(b, []geom.Point{{X: 0, Y: 0}, {X: 0, Y: 1}}, "C")

	fmt.Printf("%+v\n", b)

	fmt.Printf("%+v\n", board.RemainingPoints(b))

	if c != nil {
		(*c)()
	}
	fmt.Printf("%+v\n", b)
	if d != nil {
		(*d)()
	}
	fmt.Printf("%+v\n", b)
	if e == nil {
		fmt.Println("e is nil")
	} else {
		(*e)()
	}
}
