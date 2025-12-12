package day09

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AlexeyYurko/advent-of-code-2025/internal/aoc"
)

type Solver struct {
	points []aoc.Point
}

func New() *Solver {
	input, err := os.ReadFile(filepath.Join("internal", "solutions", "day09", "input.txt"))
	if err != nil {
		panic(err)
	}

	var points []aoc.Point
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		xs := strings.TrimSpace(parts[0])
		ys := strings.TrimSpace(parts[1])
		x, errX := strconv.Atoi(xs)
		y, errY := strconv.Atoi(ys)
		if errX != nil || errY != nil {
			continue
		}
		points = append(points, aoc.Point{X: x, Y: y})
	}

	return &Solver{points: points}
}

type interval struct {
	L, R int
}

func mergeIntervals(iv []interval) []interval {
	if len(iv) == 0 {
		return iv
	}
	sort.Slice(iv, func(i, j int) bool {
		if iv[i].L == iv[j].L {
			return iv[i].R < iv[j].R
		}
		return iv[i].L < iv[j].L
	})
	out := make([]interval, 0, len(iv))
	cur := iv[0]
	for i := 1; i < len(iv); i++ {
		if iv[i].L <= cur.R+1 {
			if iv[i].R > cur.R {
				cur.R = iv[i].R
			}
		} else {
			out = append(out, cur)
			cur = iv[i]
		}
	}
	out = append(out, cur)
	return out
}

func containsInterval(iv []interval, L, R int) bool {
	i := sort.Search(len(iv), func(i int) bool { return iv[i].L > L })
	if i == 0 {
		return false
	}
	i--
	return iv[i].L <= L && iv[i].R >= R
}

func (s *Solver) Part1() (interface{}, error) {
	maxArea := 0
	n := len(s.points)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := s.points[i]
			p2 := s.points[j]

			dx := aoc.Abs(p2.X - p1.X)
			dy := aoc.Abs(p2.Y - p1.Y)

			if dx > 0 && dy > 0 {
				area := (dx + 1) * (dy + 1)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea, nil
}

func (s *Solver) Part2() (interface{}, error) {
	n := len(s.points)
	if n < 2 {
		return 0, nil
	}

	type edge struct{ x1, y1, x2, y2 int }
	edges := make([]edge, 0, n)
	for i := 0; i < n; i++ {
		a := s.points[i]
		b := s.points[(i+1)%n]
		edges = append(edges, edge{a.X, a.Y, b.X, b.Y})
	}

	minY, maxY := s.points[0].Y, s.points[0].Y
	minX, maxX := s.points[0].X, s.points[0].X
	for _, p := range s.points {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
	}

	rowIntervals := make(map[int][]interval)

	for _, e := range edges {
		if e.y1 == e.y2 {
			y := e.y1
			L := e.x1
			R := e.x2
			if L > R {
				L, R = R, L
			}
			rowIntervals[y] = append(rowIntervals[y], interval{L, R})
		}
	}

	for y := minY; y <= maxY; y++ {
		xs := make([]int, 0, 16)
		for _, e := range edges {
			if e.x1 == e.x2 {
				x := e.x1
				y1, y2 := e.y1, e.y2
				if y1 > y2 {
					y1, y2 = y2, y1
				}
				if y >= y1 && y < y2 {
					xs = append(xs, x)
				}
			}
		}
		if len(xs) == 0 {
			continue
		}
		sort.Ints(xs)
		for i := 0; i+1 < len(xs); i += 2 {
			L := xs[i]
			R := xs[i+1]
			if R-1 >= L {
				rowIntervals[y] = append(rowIntervals[y], interval{L, R - 1})
			}
		}
	}

	for y, iv := range rowIntervals {
		rowIntervals[y] = mergeIntervals(iv)
	}

	allowedRange := func(y, L, R int) bool {
		iv, ok := rowIntervals[y]
		if !ok || len(iv) == 0 {
			return false
		}
		return containsInterval(iv, L, R)
	}

	maxArea := 0
	for i := 0; i < n; i++ {
		p1 := s.points[i]
		for j := i + 1; j < n; j++ {
			p2 := s.points[j]
			if p1.X == p2.X || p1.Y == p2.Y {
				continue
			}
			L := min(p1.X, p2.X)
			R := max(p1.X, p2.X)
			T := min(p1.Y, p2.Y)
			B := max(p1.Y, p2.Y)

			ok := true
			for y := T; y <= B; y++ {
				if !allowedRange(y, L, R) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}

			area := (aoc.Abs(p1.X-p2.X) + 1) * (aoc.Abs(p1.Y-p2.Y) + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea, nil
}
