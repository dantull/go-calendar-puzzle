package main

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"calendar-puzzle/solver"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

func parsePoint(arg string) (int, int) {
	parts := strings.Split(arg, ",")
	if len(parts) != 2 {
		fmt.Printf("Invalid point format: %s\n", arg)
		return 0, 0
	}
	x := 0
	y := 0
	fmt.Sscanf(parts[0], "%d", &x)
	fmt.Sscanf(parts[1], "%d", &y)
	return x, y
}

func main() {
	fmt.Println("Calendar Solver!")

	var cpuprofile = flag.Bool("c", false, "write cpu profile to file")

	if *cpuprofile {
		f, err := os.Create("profile")
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	width := 6
	height := 7
	ps := append(geom.Grid(width, height), geom.Point{X: 6, Y: 2}, geom.Point{X: 6, Y: 3},
		geom.Point{X: 6, Y: 4}, geom.Point{X: 6, Y: 5}, geom.Point{X: 6, Y: 6},
		geom.Point{X: 6, Y: 7}, geom.Point{X: 5, Y: 7}, geom.Point{X: 4, Y: 7})

	b := board.NewBoard(ps)

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
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}),
		"Z": *geom.NewShape(true, 1, []geom.Point{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}}),
	}

	verbose := flag.Bool("v", false, "verbose output")
	multi := flag.Int("m", 1, "multiple solutions")

	flag.Parse()

	origin := geom.Point{X: 0, Y: 0}

	// try to place each shape in position 0, 0
	for label, shape := range labeledShapes {
		vs := geom.Variants(&shape)
		placed := false

		for _, v := range vs {
			if *verbose {
				lines := geom.Stringify(v, func(p geom.Point) string {
					for _, pv := range v {
						if p == pv {
							return label
						}
					}

					return " "
				})
				for _, s := range lines {
					fmt.Printf("%s\n", s)
				}
				fmt.Println()
			}
			remove := board.FillPoints(b, &v, origin, label)

			if remove != nil {
				placed = true
				(*remove)()
			}
		}
		if !placed {
			fmt.Printf("Failed to place %s shape\n", label)
			return
		}
	}

	fills := []geom.Point{}
	for _, arg := range flag.Args() {
		x, y := parsePoint(arg)
		fills = append(fills, geom.Point{X: x, Y: y})
	}

	board.FillPoints(b, &fills, geom.Point{X: 0, Y: 0}, "*")
	dumpBoard()

	solverStepper := solver.CreateSolver(b, labeledShapes, 4)

	done := false
	tofind := *multi

	for {
		more := solverStepper(func(inspector solver.Inspector, event solver.Event) {
			if event.Kind == "solved" {
				dumpBoard()
				fmt.Println("Solved!")
				tofind--
				done = tofind <= 0
			} else if *verbose {
				dumpBoard()
				fmt.Printf("Event: %s (%s)\n", event.Kind, event.Label)
			}
		})

		if !more || done {
			break
		}
	}
}
