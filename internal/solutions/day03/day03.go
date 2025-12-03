package day03

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Solver struct {
	input string
	banks []string
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day03", "input.txt"))
	s := &Solver{input: string(input)}
	s.banks = strings.Split(s.input, "\n")
	return s
}

func (s *Solver) Part1() (interface{}, error) {
	result := 0
	for _, bank := range s.banks {
		n := len(bank)
		maxOne := -1
		pos := -1
		maxSecond := -1

		for i, b := range bank {
			d := int(b - '0')

			if i < n-1 && d > maxOne {
				maxOne = d
				pos = i
				maxSecond = -1
			}

			if pos != -1 && i > pos {
				if d > maxSecond {
					maxSecond = d
				}
			}
		}

		result += maxOne*10 + maxSecond
	}
	return result, nil
}

func (s *Solver) Part2() (interface{}, error) {
	var total int64
	for _, bank := range s.banks {
		n := len(bank)
		toRemove := n - 12

		stack := make([]byte, 0, 12+toRemove)

		for i := 0; i < n; i++ {
			d := bank[i]
			for toRemove > 0 && len(stack) > 0 && d > stack[len(stack)-1] {
				stack = stack[:len(stack)-1]
				toRemove--
			}
			stack = append(stack, d)
		}

		if toRemove > 0 {
			stack = stack[:len(stack)-toRemove]
		}

		resultStr := string(stack[:12])

		val, _ := strconv.ParseInt(resultStr, 10, 64)
		total += val
	}
	return total, nil
}
