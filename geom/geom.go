package geom

import (
	"strings"
)

type Point struct {
	X int
	Y int
}

type Shape struct {
	chiral    bool
	rotations int
	points    []Point
}

type Mapper func(Point) Point

func flipPoint(p Point) Point {
	return Point{X: p.Y, Y: -p.X}
}

func identity(p Point) Point {
	return p
}

var rotates = [4]Mapper{
	identity, // 0 degrees rotation
	func(p Point) Point {
		return Point{X: -p.Y, Y: p.X} // 90 degrees rotation
	},
	func(p Point) Point {
		return Point{X: -p.X, Y: -p.Y} // 180 degrees rotation
	},
	func(p Point) Point {
		return Point{X: p.Y, Y: -p.X} // 270 degrees rotation
	},
}

func NewShape(chiral bool, rotations int, points []Point) *Shape {
	return &Shape{
		chiral:    chiral,
		rotations: rotations,
		points:    points,
	}
}

func AddPoints(a, b Point) Point {
	return Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func SubtractPoints(a, b Point) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func Bounds(ps []Point) [2]Point {
	var minX int
	var minY int
	var maxX int
	var maxY int

	if len(ps) > 0 {
		minX = ps[0].X
		minY = ps[0].Y
		maxX = ps[0].X
		maxY = ps[0].Y
	} else {
		minX = 0
		minY = 0
		maxX = 0
		maxY = 0
	}

	for _, p := range ps {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return [2]Point{{X: minX, Y: minY}, {X: maxX, Y: maxY}}
}

func Variants(s *Shape) [][]Point {
	var flips []Mapper
	if s.chiral {
		flips = []Mapper{identity, flipPoint}
	} else {
		flips = []Mapper{identity}
	}
	vs := make([][]Point, 0, s.rotations*len(flips))

	for _, flip := range flips {
		for i := 0; i <= s.rotations; i++ {
			rotate := rotates[i%4]
			v := make([]Point, len(s.points))
			for i, p := range s.points {
				v[i] = rotate(flip(p))
			}
			vs = append(vs, v)
		}
	}

	return vs
}

func Stringify(ps []Point, convert func(Point) string) []string {
	bounds := Bounds(ps)
	width := bounds[1].X - bounds[0].X + 1
	height := bounds[1].Y - bounds[0].Y + 1

	strs := make([]string, height)
	for i := range strs {
		sb := strings.Builder{}
		for j := 0; j < width; j++ {
			sb.WriteString(convert(AddPoints(Point{X: j, Y: i}, bounds[0])))
		}
		strs[i] = sb.String()
	}

	return strs
}

func Grid(width int, height int) []Point {
	var points []Point
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			points = append(points, Point{X: x, Y: y})
		}
	}
	return points
}

func OffsetAll(ps []Point, offset Point) []Point {
	var offsetPoints []Point
	for _, p := range ps {
		offsetPoints = append(offsetPoints, AddPoints(p, offset))
	}
	return offsetPoints
}
