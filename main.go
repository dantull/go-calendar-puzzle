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
	height := 4
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

	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	labeledShapes := make(map[string]geom.Shape, len(labels))
	for _, label := range labels {
		labeledShapes[label] = *geom.NewShape(true, 4, shapePoints)
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
