package day07

import (
	"os"
	"path/filepath"
	"strings"
)

type Solver struct {
	grid   []string
	startR int
	startC int
	height int
	width  int
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day07", "input.txt"))
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")

	var startR, startC int
	for r, line := range lines {
		if c := strings.IndexByte(line, 'S'); c != -1 {
			startR, startC = r, c
			break
		}
	}

	return &Solver{
		grid:   lines,
		startR: startR,
		startC: startC,
		height: len(lines),
		width:  len(lines[0]),
	}
}

func (s *Solver) Part1() (interface{}, error) {
	type pos struct{ r, c int }
	queue := []pos{{s.startR, s.startC}}
	visited := make([][]bool, s.height)
	for i := range visited {
		visited[i] = make([]bool, s.width)
	}
	visited[s.startR][s.startC] = true

	splits := 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		r, c := curr.r, curr.c

		nextR := r + 1
		if nextR >= s.height {
			continue
		}

		cell := s.grid[nextR][c]
		switch cell {
		case '.':
			if !visited[nextR][c] {
				visited[nextR][c] = true
				queue = append(queue, pos{nextR, c})
			}
		case '^':
			splits++
			if c-1 >= 0 && !visited[nextR][c-1] {
				visited[nextR][c-1] = true
				queue = append(queue, pos{nextR, c - 1})
			}
			if c+1 < s.width && !visited[nextR][c+1] {
				visited[nextR][c+1] = true
				queue = append(queue, pos{nextR, c + 1})
			}
		default:
			if !visited[nextR][c] {
				visited[nextR][c] = true
				queue = append(queue, pos{nextR, c})
			}
		}
	}

	return splits, nil
}

func (s *Solver) Part2() (interface{}, error) {
	dp := make([]int64, s.width)
	dp[s.startC] = 1

	var total int64 = 0

	for r := s.startR; r < s.height; r++ {
		nextR := r + 1

		nextDP := make([]int64, s.width)

		for c := 0; c < s.width; c++ {
			count := dp[c]
			if count == 0 {
				continue
			}

			if nextR >= s.height {
				total += count
			} else {
				cell := s.grid[nextR][c]
				switch cell {
				case '.':
					nextDP[c] += count
				case '^':
					if c-1 >= 0 {
						nextDP[c-1] += count
					}
					if c+1 < s.width {
						nextDP[c+1] += count
					}
				default:
					nextDP[c] += count
				}
			}
		}

		if nextR < s.height {
			dp = nextDP
		}
	}

	return total, nil
}
