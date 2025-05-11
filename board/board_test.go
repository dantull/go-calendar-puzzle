package board_test

import (
	"testing"

	"go-calendar-puzzle/board"
	"go-calendar-puzzle/geom"
)

func TestNewBoard(t *testing.T) {
	points := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
	b := board.NewBoard(points)

	if len(board.RemainingPoints(b)) != len(points) {
		t.Errorf("Expected %d unfilled points", len(points))
	}
	for _, p := range points {
		if board.LabelAt(b, p) != nil {
			t.Errorf("Expected no label at point %v", p)
		}
	}
}

func TestFillPoints(t *testing.T) {
	points := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
	b := board.NewBoard(points)

	fillA := []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 0}}
	fillB := []geom.Point{{X: 1, Y: 1}}

	cleanupA := board.FillPoints(b, fillA, "A")
	if cleanupA == nil {
		t.Fatal("Expected non-nil cleanup function")
	}

	if len(board.RemainingPoints(b)) != 2 {
		t.Errorf("Expected 2 unfilled points, got %d", len(board.RemainingPoints(b)))
	}

	cleanupB := board.FillPoints(b, fillB, "B")
	if cleanupB == nil {
		t.Fatal("Expected non-nil cleanup function")
	}

	if len(board.RemainingPoints(b)) != 1 {
		t.Errorf("Expected 1 unfilled point, got %d", len(board.RemainingPoints(b)))
	}

	if board.FillPoints(b, fillA, "X") != nil {
		t.Error("Expected nil cleanup function for overlapping fill")
	}

	for _, p := range fillA {
		if label := board.LabelAt(b, p); label == nil || *label != "A" {
			t.Errorf("Expected label 'A' at point %v, got %v", p, label)
		}
	}

	for _, p := range fillB {
		if label := board.LabelAt(b, p); label == nil || *label != "B" {
			t.Errorf("Expected label 'B' at point %v, got %v", p, label)
		}
	}

	(*cleanupA)()
	if len(board.RemainingPoints(b)) != 3 {
		t.Errorf("Expected 2 unfilled points after cleanupA, got %d", len(board.RemainingPoints(b)))
	}

	(*cleanupB)()
	if len(board.RemainingPoints(b)) != 4 {
		t.Errorf("Expected 4 unfilled points after cleanupB, got %d", len(board.RemainingPoints(b)))
	}

	for _, p := range fillA {
		if label := board.LabelAt(b, p); label != nil {
			t.Errorf("Expected no label at point %v after cleanupA, got %v", p, label)
		}
	}

	for _, p := range fillB {
		if label := board.LabelAt(b, p); label != nil {
			t.Errorf("Expected no label at point %v after cleanupB, got %v", p, label)
		}
	}
}

func TestReachable(t *testing.T) {
	points := []geom.Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
	b := board.NewBoard(points)

	fillA := []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 0}}
	fillB := []geom.Point{{X: 1, Y: 1}}

	cleanupA := board.FillPoints(b, fillA, "A")
	if cleanupA == nil {
		t.Fatal("Expected non-nil cleanup function")
	}

	if board.CountFill(b, fillA[0], 1) != 0 {
		t.Error("Expected 0 reachable points from fillA[0]")
	}

	if board.CountFill(b, fillB[0], 3) != 2 {
		t.Errorf("Expected 2 reachable points from fillB[0]")
	}

	cleanupB := board.FillPoints(b, fillB, "B")
	if cleanupB == nil {
		t.Fatal("Expected non-nil cleanup function")
	}

	(*cleanupA)()

	if board.CountFill(b, fillA[1], 4) != 3 {
		t.Errorf("Expected 2 reachable points from fillA[1]")
	}

	if board.CountFill(b, fillA[1], 2) != 2 {
		t.Errorf("Expected 2 reachable points from fillA[1]")
	}
}
