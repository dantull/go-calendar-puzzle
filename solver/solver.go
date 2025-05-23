package solver

import (
	"calendar-puzzle/board"
	"calendar-puzzle/geom"
	"sort"
)

type ShapeState struct {
	pointIndex   int
	variantIndex int
	remove       *func()
	places       int
	shape        *geom.Shape
	label        string
	baseVariants [][]geom.Point
	points       []geom.Point
}

func newState(shape *geom.Shape, label string, variants [][]geom.Point, points []geom.Point) ShapeState {

	return ShapeState{
		pointIndex:   0,
		variantIndex: 0,
		remove:       nil,
		places:       0,
		shape:        shape,
		label:        label,
		baseVariants: variants,
		points:       points,
	}
}

func stepState(s *ShapeState, b *board.Board, minSize int) bool {
	if s.remove != nil {
		(*s.remove)()
		s.remove = nil
	}

	if s.variantIndex < len(s.baseVariants) {
		s.remove = board.FillPoints(b, &s.baseVariants[s.variantIndex], s.points[s.pointIndex], s.label)

		if s.remove != nil {
			s.places++
		}

		s.variantIndex++
	} else if s.variantIndex == len(s.baseVariants) {
		p := s.points[s.pointIndex]
		if board.CountFill(b, p, minSize) < minSize {
			s.pointIndex = len(s.points)
		} else {
			s.pointIndex++
		}
		s.variantIndex = 0
	}

	return s.pointIndex < len(s.points)
}

func isPlaced(s *ShapeState) bool {
	return s.remove != nil
}

func neverPlaced(s *ShapeState) bool {
	return s.places == 0
}

type Event struct {
	Kind  string
	Label string
	shape *geom.Shape
}

type Inspector func(p geom.Point) string
type SolverCallback func(Inspector, Event)
type SolverStepper func(SolverCallback) bool

func CreateSolver(b *board.Board, labeledShapes map[string]geom.Shape, minSize int) SolverStepper {
	stack := make([]*ShapeState, 0)

	keys := make([]string, 0, len(labeledShapes))
	for k := range labeledShapes {
		keys = append(keys, k)
	}

	allVariants := make(map[string][][]geom.Point, len(labeledShapes))

	for l, s := range labeledShapes {
		vs := geom.Variants(&s)
		allVariants[l] = vs
	}

	sort.Slice(keys, func(i, j int) bool {
		il := len(allVariants[keys[i]][0])
		jl := len(allVariants[keys[j]][0])

		if il == jl {
			return keys[i] > keys[j]
		} else {
			return il > jl
		}
	})

	var nextShape = func() ShapeState {
		label := keys[len(stack)]
		s := labeledShapes[label]

		return newState(&s, label, allVariants[label], board.RemainingPoints(b))
	}

	if len(labeledShapes) > 0 {
		ns := nextShape()
		stack = append(stack, &ns)
	}

	return func(callback SolverCallback) bool {
		if len(stack) == 0 {
			return false
		}

		shapeState := stack[len(stack)-1]

		more := stepState(shapeState, b, minSize)

		if !more {
			if neverPlaced(shapeState) {
				callback(nil, Event{Kind: "failed", Label: shapeState.label, shape: shapeState.shape})
			}

			stack = stack[:len(stack)-1]
			return len(stack) > 0
		} else {
			if isPlaced(shapeState) {
				solved := len(stack) == len(labeledShapes)

				event := "placed"
				if solved {
					event = "solved"
				}

				callback(nil, Event{Kind: event, Label: shapeState.label, shape: shapeState.shape})

				if !solved {
					ns := nextShape()
					stack = append(stack, &ns)
				}
			}

			return true
		}
	}
}
