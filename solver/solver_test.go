package solver_test

import (
	"testing"

	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"calendar-puzzle/solver"
)

func TestTrivialSolve(t *testing.T) {
	ps := geom.Grid(2, 2)
	b := board.NewBoard(ps)
	shape := geom.NewShape(false, 0, ps)

	stepper := solver.CreateSolver(b, map[string]geom.Shape{"*": *shape}, 4)

	calls := 0
	cb := func(inspector solver.Inspector, event solver.Event) {
		calls++

		if event.Kind != "solved" {
			t.Errorf("Expected event kind 'solved', got '%s'", event.Kind)
		}

		if calls > 1 {
			t.Errorf("Expected only one call to the callback, got %d", calls)
		}
	}

	more := true
	for {
		if !more {
			break
		}

		more = stepper(cb)
	}

	if calls != 1 {
		t.Errorf("Expected one call to the callback, got %d", calls)
	}
}

func TestFlipSolve(t *testing.T) {
	shape_ps := []geom.Point{
		{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}}

	board_ps := []geom.Point{
		{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 1}}

	b := board.NewBoard(board_ps)
	shape := geom.NewShape(true, 3, shape_ps)

	stepper := solver.CreateSolver(b, map[string]geom.Shape{"*": *shape}, 3)

	solved := false
	for {
		more := stepper(func(inspector solver.Inspector, event solver.Event) {
			if event.Kind == "solved" {
				solved = true
			}
		})

		if !more {
			break
		}
	}

	if !solved {
		t.Errorf("Expected to solve the board, but it was not solved")
	}
}

type SolverCase struct {
	name      string
	chiral    bool
	rotations int
	events    []string
}

func TestBasicSolverCases(t *testing.T) {
	cases := []SolverCase{
		{
			name:      "impossible",
			rotations: 0,
			chiral:    false,
			events:    []string{"failed"},
		},
		{
			name:      "rotation",
			rotations: 1,
			chiral:    false,
			events:    []string{"solved"},
		},
		{
			name:      "extra rotations",
			rotations: 3,
			chiral:    false,
			events:    []string{"solved", "solved"},
		},
		{
			name:      "flip",
			rotations: 1,
			chiral:    true,
			events:    []string{"solved", "solved"},
		},
		{
			name:      "all rotations and flip",
			rotations: 3,
			chiral:    true,
			events:    []string{"solved", "solved", "solved", "solved"},
		},
	}

	p := geom.Grid(1, 2)
	b := board.NewBoard(geom.Grid(2, 1))

	for _, c := range cases {
		shape := geom.NewShape(c.chiral, c.rotations, p)
		stepper := solver.CreateSolver(b, map[string]geom.Shape{"*": *shape}, 2)
		events := make([]string, 0)

		for {
			more := stepper(func(inspector solver.Inspector, event solver.Event) {
				events = append(events, event.Kind)
			})

			if !more {
				break
			}
		}

		if len(events) != len(c.events) {
			t.Errorf("Expected %d events, got %d", len(c.events), len(events))
		}

		for i, event := range events {
			if event != c.events[i] {
				t.Errorf("Expected event %d to be '%s', got '%s'", i, c.events[i], event)
			}
		}
	}
}

func TestSolverCase(t *testing.T) {
	ps := geom.Grid(2, 3)
	b := board.NewBoard(ps)
	shape := geom.NewShape(false, 3, []geom.Point{
		{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}})

	stepper := solver.CreateSolver(b, map[string]geom.Shape{"x": *shape, "o": *shape}, 2)

	counts := map[string]int{
		"failed": 4,
		"placed": 8,
		"solved": 4,
	}

	for {
		more := stepper(func(inspector solver.Inspector, event solver.Event) {
			counts[event.Kind]--
		})

		if !more {
			break
		}
	}

	for k, v := range counts {
		if v != 0 {
			t.Errorf("Expected all counts to be 0, got %d for %s", v, k)
		}
	}
}
