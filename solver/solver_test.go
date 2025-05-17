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
