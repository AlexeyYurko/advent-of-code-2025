package day02

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Solver struct {
	input    string
	idRanges []string
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day02", "input.txt"))
	s := &Solver{input: string(input)}
	s.idRanges = strings.Split(s.input, ",")
	return s
}

func (s *Solver) Part1() (interface{}, error) {
	result := 0
	for _, idRange := range s.idRanges {
		var start, end int
		fmt.Sscanf(idRange, "%d-%d", &start, &end)

		for i := start; i <= end; i++ {
			str := strconv.Itoa(i)
			n := len(str)
			if n%2 != 0 {
				continue
			}
			if str[:n/2] == str[n/2:] {
				result += i
			}
		}
	}
	return result, nil
}

func isRepetition(s string) bool {
	n := len(s)
	if n < 2 {
		return false
	}
	for pLen := 1; pLen <= n/2; pLen++ {
		if n%pLen != 0 {
			continue
		}
		for i := pLen; i < n; i++ {
			if s[i] != s[i%pLen] {
				goto nextPeriod
			}
		}
		return true
	nextPeriod:
	}
	return false
}

func (s *Solver) Part2() (interface{}, error) {
	result := 0
	for _, idRange := range s.idRanges {
		var start, end int
		fmt.Sscanf(idRange, "%d-%d", &start, &end)

		for i := start; i <= end; i++ {
			str := strconv.Itoa(i)
			if isRepetition(str) {
				result += i
			}
		}
	}
	return result, nil
}
