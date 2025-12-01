package aoc

type Point struct {
	X, Y int
}

func (p Point) Add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

func (p Point) Manhattan(q Point) int {
	return Abs(p.X-q.X) + Abs(p.Y-q.Y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
