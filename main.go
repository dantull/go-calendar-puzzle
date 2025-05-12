package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"calendar-puzzle/solver"
	"fmt"
	"strings"
)

func main() {
	width := 6
	height := 7
	ps := append(geom.Grid(width, height), geom.Point{X: 6, Y: 2}, geom.Point{X: 6, Y: 3},
		geom.Point{X: 6, Y: 4}, geom.Point{X: 6, Y: 5}, geom.Point{X: 6, Y: 6},
		geom.Point{X: 6, Y: 7}, geom.Point{X: 5, Y: 7}, geom.Point{X: 4, Y: 7})

	b := board.NewBoard(ps)
	board.FillPoints(b, []geom.Point{{X: 4, Y: 0}, {X: 3, Y: 3}, {X: 3, Y: 6}}, "X")

	var boardAt = func(p geom.Point) string {
		if label := board.LabelAt(b, p); label != nil {
			return *label
		}
		for _, pe := range ps {
			if pe == p {
				return "-"
			}
		}
		return " "
	}

	var dumpBoard = func() {
		fmt.Println()
		fmt.Printf("%+v\n", strings.Join(geom.Stringify(ps, boardAt), "\n"))
		fmt.Println()
	}

	dumpBoard()

	labeledShapes := map[string]geom.Shape{
		"I": *geom.NewShape(false, 1, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}}),
		"L": *geom.NewShape(true, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}}),
		"S": *geom.NewShape(true, 1, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}}),
		"J": *geom.NewShape(true, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 0, Y: 1}}),
		"N": *geom.NewShape(true, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 3, Y: 1}}),
		"U": *geom.NewShape(false, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 2, Y: 1}}),
		"T": *geom.NewShape(false, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}}),
		"P": *geom.NewShape(true, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}),
		"V": *geom.NewShape(false, 3, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 2}}),
		"Z": *geom.NewShape(true, 1, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}}),
	}

	solverStepper := solver.CreateSolver(b, labeledShapes, 3)

	done := false

	for {
		more := solverStepper(func(inspector solver.Inspector, event solver.Event) {

			if event.Kind == "solved" {
				dumpBoard()
				fmt.Println("Solved!")
				done = true
			}
		})

		if !more || done {
			break
		}
	}
}
