package day08

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AlexeyYurko/advent-of-code-2025/internal/aoc"
)

type Solver struct {
	points    []aoc.Point3D
	edges     []Edge
	numPoints int
	dsu       *DSU
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day08", "input.txt"))
	var points []aoc.Point3D
	for _, line := range strings.Split(string(input), "\n") {
		pointsData := strings.Split(line, ",")
		x, _ := strconv.Atoi(pointsData[0])
		y, _ := strconv.Atoi(pointsData[1])
		z, _ := strconv.Atoi(pointsData[2])
		points = append(points, aoc.Point3D{X: x, Y: y, Z: z})
	}

	n := len(points)

	edges := make([]Edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := distSq(points[i], points[j])
			edges = append(edges, Edge{dist: d, i: i, j: j})
		}
	}

	sort.Slice(edges, func(a, b int) bool {
		return edges[a].dist < edges[b].dist
	})

	return &Solver{
		points:    points,
		edges:     edges,
		numPoints: n,
		dsu:       NewDSU(n),
	}
}

type DSU struct {
	parent []int
	size   []int
}

type Edge struct {
	dist int64
	i, j int
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) bool {
	xr, yr := d.Find(x), d.Find(y)
	if xr == yr {
		return false
	}
	if d.size[xr] < d.size[yr] {
		xr, yr = yr, xr
	}
	d.parent[yr] = xr
	d.size[xr] += d.size[yr]
	return true
}

func (d *DSU) ComponentSizes() []int {
	sizes := make([]int, 0)
	for i := range d.parent {
		if d.parent[i] == i {
			sizes = append(sizes, d.size[i])
		}
	}
	return sizes
}

func distSq(p, q aoc.Point3D) int64 {
	dx := int64(p.X - q.X)
	dy := int64(p.Y - q.Y)
	dz := int64(p.Z - q.Z)
	return dx*dx + dy*dy + dz*dz
}

func (s *Solver) Part1() (interface{}, error) {
	pairsToConnect := 1000

	for i := 0; i < len(s.edges) && i < pairsToConnect; i++ {
		e := s.edges[i]
		s.dsu.Union(e.i, e.j)
	}

	sizes := s.dsu.ComponentSizes()
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	for len(sizes) < 3 {
		sizes = append(sizes, 1)
	}

	product := sizes[0] * sizes[1] * sizes[2]
	return product, nil
}

func (s *Solver) Part2() (interface{}, error) {
	unions := 0
	var lastEdge Edge

	for _, e := range s.edges {
		if s.dsu.Union(e.i, e.j) {
			lastEdge = e
			unions++
			if unions == s.numPoints-1 {
				break
			}
		}
	}

	x1 := s.points[lastEdge.i].X
	x2 := s.points[lastEdge.j].X
	return x1 * x2, nil
}
