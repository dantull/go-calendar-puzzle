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

type SolverCase struct {
	name      string
	chiral    bool
	rotations int
	events    []string
}

func TestSolver(t *testing.T) {
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
