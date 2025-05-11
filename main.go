package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"fmt"
	"math/rand"
	"strings"
)

func main() {
	width := 12
	height := 8
	ps := geom.Grid(width, height)
	b := board.NewBoard(ps)

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

	shapePoints := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
	}

	undos := make([]func(), 0)

	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}

	variants := geom.Variants(geom.NewShape(true, 4, shapePoints))

	for {
		if len(undos) == len(labels) {
			break
		}

		h := rand.Intn(height - 1)
		w := rand.Intn(width - 1)
		v := rand.Intn(len(variants))

		r := board.FillPoints(b, geom.OffsetAll(variants[v], geom.Point{X: w, Y: h}), labels[len(undos)])
		if r != nil {
			undos = append(undos, *r)
			dumpBoard()
		}
	}

	for _, u := range undos {
		u()
		dumpBoard()
	}
}
