package day12

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type shape struct {
	cells [][2]int
}

type region struct {
	w, h int
	qs   []int
}

type Solver struct {
	shapes  []shape
	regions []region
}

func New() *Solver {
	input, _ := os.ReadFile(
		filepath.Join("internal", "solutions", "day12", "input.txt"),
	)

	shapes, regions := parseInput(string(input))

	return &Solver{shapes: shapes, regions: regions}
}

func parseInput(raw string) ([]shape, []region) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")

	var shapes []shape
	var regions []region

	i := 0
	for i < len(lines) {
		line := lines[i]
		line = strings.TrimSpace(line)
		if line == "" {
			i++
			continue
		}
		if len(line) >= 2 && strings.HasSuffix(line, ":") && !strings.Contains(line, "x") {
			grid := lines[i+1 : i+4]
			var cells [][2]int
			for r, row := range grid {
				for c, ch := range row {
					if ch == '#' {
						cells = append(cells, [2]int{r, c})
					}
				}
			}
			shapes = append(shapes, shape{cells: cells})
			i += 4
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			dim := line[:idx]
			parts := strings.Split(dim, "x")
			if len(parts) == 2 {
				w, _ := strconv.Atoi(parts[0])
				h, _ := strconv.Atoi(parts[1])
				rest := strings.TrimSpace(line[idx+1:])
				nums := strings.Fields(rest)
				qs := make([]int, len(nums))
				for j, n := range nums {
					qs[j], _ = strconv.Atoi(n)
				}
				regions = append(regions, region{w: w, h: h, qs: qs})
			}
		}
		i++
	}

	return shapes, regions
}

type placement struct {
	dr, dc  int
	variant [][2]int
}

func normalize(cells [][2]int) [][2]int {
	minR, minC := cells[0][0], cells[0][1]
	for _, c := range cells {
		if c[0] < minR {
			minR = c[0]
		}
		if c[1] < minC {
			minC = c[1]
		}
	}
	out := make([][2]int, len(cells))
	for k, c := range cells {
		out[k] = [2]int{c[0] - minR, c[1] - minC}
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a][0] != out[b][0] {
			return out[a][0] < out[b][0]
		}
		return out[a][1] < out[b][1]
	})
	return out
}

func keyOf(cells [][2]int) string {
	var b strings.Builder
	for _, c := range cells {
		b.WriteByte(byte('0' + c[0]))
		b.WriteByte(byte('0' + c[1]))
	}
	return b.String()
}

func orientations(cells [][2]int) [][][2]int {
	seen := map[string]bool{}
	var result [][][2]int
	current := make([][2]int, len(cells))
	copy(current, cells)
	for r := 0; r < 4; r++ {
		for f := 0; f < 2; f++ {
			n := normalize(current)
			k := keyOf(n)
			if !seen[k] {
				seen[k] = true
				result = append(result, n)
			}
			current = flip(current)
		}
		current = rotate(current)
	}
	return result
}

func rotate(cells [][2]int) [][2]int {
	out := make([][2]int, len(cells))
	for i, c := range cells {
		out[i] = [2]int{c[1], -c[0]}
	}
	return out
}

func flip(cells [][2]int) [][2]int {
	out := make([][2]int, len(cells))
	for i, c := range cells {
		out[i] = [2]int{c[0], -c[1]}
	}
	return out
}

func placementsFor(sh shape, w, h int) []placement {
	variants := orientations(sh.cells)
	var res []placement
	for _, v := range variants {
		maxR, maxC := 0, 0
		for _, c := range v {
			if c[0] > maxR {
				maxR = c[0]
			}
			if c[1] > maxC {
				maxC = c[1]
			}
		}
		for dr := 0; dr <= h-1-maxR; dr++ {
			for dc := 0; dc <= w-1-maxC; dc++ {
				res = append(res, placement{dr: dr, dc: dc, variant: v})
			}
		}
	}
	return res
}

func canPlace(grid []byte, w int, dr, dc int, v [][2]int) bool {
	for _, c := range v {
		if grid[(dr+c[0])*w+dc+c[1]] == 1 {
			return false
		}
	}
	return true
}

func place(grid []byte, w int, dr, dc int, v [][2]int, val byte) {
	for _, c := range v {
		grid[(dr+c[0])*w+dc+c[1]] = val
	}
}

const slackThreshold = 345

func canFit(w, h int, qs []int, shapes []shape) bool {
	cellCounts := make([]int, len(shapes))
	totalCells := 0
	for i := range shapes {
		cellCounts[i] = len(shapes[i].cells)
	}
	for i := range qs {
		totalCells += qs[i] * cellCounts[i]
	}
	area := w * h
	if totalCells > area {
		return false
	}
	if area-totalCells > slackThreshold {
		return true
	}

	type shapePlacements struct {
		idx  int
		list []placement
	}
	var sps []shapePlacements
	for idx := range shapes {
		if qs[idx] > 0 {
			sps = append(sps, shapePlacements{idx: idx, list: placementsFor(shapes[idx], w, h)})
		}
	}
	sort.Slice(sps, func(a, b int) bool {
		return len(sps[a].list) < len(sps[b].list)
	})

	var slots []int
	for _, sp := range sps {
		for k := 0; k < qs[sp.idx]; k++ {
			slots = append(slots, sp.idx)
		}
	}

	listByIdx := map[int][]placement{}
	for _, sp := range sps {
		listByIdx[sp.idx] = sp.list
	}

	grid := make([]byte, w*h)

	var backtrack func(si, start int) int
	backtrack = func(si, start int) int {
		if si == len(slots) {
			return 1
		}
		idx := slots[si]
		lst := listByIdx[idx]
		sameNext := si+1 < len(slots) && slots[si+1] == idx
		for k := start; k < len(lst); k++ {
			p := lst[k]
			if canPlace(grid, w, p.dr, p.dc, p.variant) {
				place(grid, w, p.dr, p.dc, p.variant, 1)
				nextStart := 0
				if sameNext {
					nextStart = k + 1
				}
				r := backtrack(si+1, nextStart)
				place(grid, w, p.dr, p.dc, p.variant, 0)
				if r == 1 {
					return 1
				}
			}
		}
		return 0
	}

	return backtrack(0, 0) == 1
}

func (s *Solver) Part1() (interface{}, error) {
	count := 0
	for _, r := range s.regions {
		if canFit(r.w, r.h, r.qs, s.shapes) {
			count++
		}
	}
	return count, nil
}

func (s *Solver) Part2() (interface{}, error) {
	return nil, nil
}
