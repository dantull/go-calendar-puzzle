package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"fmt"
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
	width := 6
	height := 4
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

	labels := []string{"A", "B", "C"}

	for i, label := range labels {
		r := board.FillPoints(b, offsetAll(shape, geom.Point{X: i, Y: i}), label)
		if r != nil {
			undos = append(undos, *r)
		}
		dumpBoard()
	}

	for _, u := range undos {
		u()
		dumpBoard()
	}
}
