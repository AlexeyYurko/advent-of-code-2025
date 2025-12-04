package day04

import (
	"os"
	"path/filepath"
	"strings"
)

type Solver struct {
	grid [][]byte
	rows int
	cols int
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day04", "input.txt"))
	lines := strings.Split(string(input), "\n")

	rows := len(lines)
	cols := len(lines[0])
	grid := make([][]byte, rows)
	for i, line := range lines {
		grid[i] = []byte(line)
	}

	return &Solver{
		grid: grid,
		rows: rows,
		cols: cols,
	}
}

func countNeighbors(grid [][]byte, i, j, rows, cols int) int {
	count := 0
	if i > 0 && j > 0 && grid[i-1][j-1] == '@' {
		count++
	}
	if i > 0 && grid[i-1][j] == '@' {
		count++
	}
	if i > 0 && j+1 < cols && grid[i-1][j+1] == '@' {
		count++
	}
	if j > 0 && grid[i][j-1] == '@' {
		count++
	}
	if j+1 < cols && grid[i][j+1] == '@' {
		count++
	}
	if i+1 < rows && j > 0 && grid[i+1][j-1] == '@' {
		count++
	}
	if i+1 < rows && grid[i+1][j] == '@' {
		count++
	}
	if i+1 < rows && j+1 < cols && grid[i+1][j+1] == '@' {
		count++
	}
	return count
}

func (s *Solver) Part1() (interface{}, error) {
	if s.rows == 0 {
		return 0, nil
	}

	count := 0
	rows, cols := s.rows, s.cols
	grid := s.grid

	for i := 0; i < rows; i++ {
		row := grid[i]
		for j := 0; j < cols; j++ {
			if row[j] != '@' {
				continue
			}
			if countNeighbors(grid, i, j, rows, cols) < 4 {
				count++
			}
		}
	}
	return count, nil
}

func (s *Solver) Part2() (interface{}, error) {
	if s.rows == 0 {
		return 0, nil
	}

	totalAts := 0
	for _, row := range s.grid {
		totalAts += strings.Count(string(row), "@")
	}

	grid := s.grid
	rows, cols := s.rows, s.cols

	toRemove := make([][2]int, 0, totalAts)

	totalRemoved := 0

	for {
		toRemove = toRemove[:0]

		for i := 0; i < rows; i++ {
			row := grid[i]
			for j := 0; j < cols; j++ {
				if row[j] != '@' {
					continue
				}
				if countNeighbors(grid, i, j, rows, cols) < 4 {
					toRemove = append(toRemove, [2]int{i, j})
				}
			}
		}

		if len(toRemove) == 0 {
			break
		}

		for _, p := range toRemove {
			grid[p[0]][p[1]] = '.'
		}
		totalRemoved += len(toRemove)
	}

	return totalRemoved, nil
}
