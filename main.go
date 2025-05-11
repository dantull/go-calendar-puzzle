package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"fmt"
	"math/rand"
	"strings"
)

func grid(width int, height int) []geom.Point {
	var points []geom.Point
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			points = append(points, geom.Point{X: x, Y: y})
		}
	}
	return points
}

func offsetAll(ps []geom.Point, offset geom.Point) []geom.Point {
	var offsetPoints []geom.Point
	for _, p := range ps {
		offsetPoints = append(offsetPoints, geom.AddPoints(p, offset))
	}
	return offsetPoints
}

func main() {
	width := 12
	height := 8
	ps := grid(width, height)
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

	shape := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
	}

	undos := make([]func(), 0)

	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}

	for _, label := range labels {
		h := rand.Intn(height - 1)
		w := rand.Intn(width - 1)

		r := board.FillPoints(b, offsetAll(shape, geom.Point{X: h, Y: w}), label)
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
