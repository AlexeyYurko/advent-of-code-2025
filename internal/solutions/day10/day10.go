package day10

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Machine struct {
	Target  []int
	Buttons [][]int
	Joltage []int
}

type Solver struct {
	machines []Machine
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day10", "input.txt"))

	machines := parseMachines(string(input))

	return &Solver{machines: machines}
}

func parseMachines(raw string) []Machine {
	lines := strings.Split(strings.TrimSpace(raw), "\n")

	re := regexp.MustCompile(`\[([.#]*)]|\(([^)]*)\)|\{([^}]*)}`)

	var machines []Machine

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m := Machine{}
		matches := re.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			switch {
			case match[1] != "": // [target]
				m.Target = parseTarget(match[1])
			case match[2] != "": // (button indices)
				indices := parseIntList(match[2])
				m.Buttons = append(m.Buttons, indices)
			case match[3] != "": // {joltage}
				joltage := parseIntList(match[3])
				m.Joltage = joltage
			}
		}

		machines = append(machines, m)
	}

	return machines
}

func parseTarget(s string) []int {
	out := make([]int, 0, len(s))
	for _, r := range s {
		switch r {
		case '.':
			out = append(out, 0)
		case '#':
			out = append(out, 1)
		}
	}
	return out
}

func parseIntList(s string) []int {
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, _ := strconv.Atoi(part)
		out = append(out, n)
	}
	return out
}

func (s *Solver) Part1() (interface{}, error) {
	const inf = 1 << 30
	totalPresses := 0

	for _, m := range s.machines {
		n := len(m.Target)
		mCount := len(m.Buttons)
		minPress := inf

		for mask := 0; mask < (1 << mCount); mask++ {
			state := make([]int, n)
			presses := 0

			for i := 0; i < mCount; i++ {
				if mask&(1<<i) != 0 {
					presses++
					for _, idx := range m.Buttons[i] {
						if idx >= 0 && idx < n {
							state[idx] ^= 1
						}
					}
				}
			}

			if presses < minPress && slices.Equal(state, m.Target) {
				minPress = presses
			}
		}

		totalPresses += minPress
	}

	return totalPresses, nil
}

func key(target []int) string {
	if len(target) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(target) * 3)
	for i, v := range target {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

func minPresses(buttons [][]int, target []int, memo map[string]int) int {
	if slices.Equal(target, make([]int, len(target))) {
		return 0
	}

	k := key(target)
	if val, ok := memo[k]; ok {
		return val
	}

	n := len(buttons)
	m := len(target)
	best := -1

	for mask := 0; mask < (1 << n); mask++ {
		remainder := make([]int, m)
		copy(remainder, target)

		cost := 0
		possible := true

		for b := 0; b < n; b++ {
			if mask&(1<<b) != 0 {
				cost++
				for i := 0; i < m; i++ {
					remainder[i] -= buttons[b][i]
				}
			}
		}

		for i := 0; i < m; i++ {
			if remainder[i] < 0 || remainder[i]%2 != 0 {
				possible = false
				break
			}
		}

		if !possible {
			continue
		}

		next := make([]int, m)
		for i := 0; i < m; i++ {
			next[i] = remainder[i] / 2
		}

		sub := minPresses(buttons, next, memo)
		if sub != -1 {
			total := cost + 2*sub
			if best == -1 || total < best {
				best = total
			}
		}
	}

	memo[k] = best
	return best
}

func buildButtonMatrix(m Machine) [][]int {
	k := len(m.Joltage)
	matrix := make([][]int, 0, len(m.Buttons))

	for _, btn := range m.Buttons {
		row := make([]int, k)
		for _, idx := range btn {
			if idx >= 0 && idx < k {
				row[idx] = 1
			}
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func (s *Solver) Part2() (interface{}, error) {
	total := 0

	for _, m := range s.machines {
		buttonMatrix := buildButtonMatrix(m)
		memo := make(map[string]int)

		presses := minPresses(buttonMatrix, m.Joltage, memo)
		total += presses
	}

	return total, nil
}
