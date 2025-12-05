package day05

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start, end int
}

type Solver struct {
	ranges []Range
	ids    []int
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day05", "input.txt"))
	base := strings.Split(string(input), "\n\n")

	ranges := strings.Split(base[0], "\n")
	parsedRanges := make([]Range, len(ranges))
	for i, line := range ranges {
		startEnd := strings.Split(line, "-")
		start, _ := strconv.Atoi(startEnd[0])
		end, _ := strconv.Atoi(startEnd[1])
		parsedRanges[i] = Range{start: start, end: end}
	}

	ids := strings.Split(base[1], "\n")
	parsedIds := make([]int, len(ids))
	for i, id := range ids {
		parsedIds[i], _ = strconv.Atoi(id)
	}

	return &Solver{
		ranges: parsedRanges,
		ids:    parsedIds,
	}
}

func (s *Solver) Part1() (interface{}, error) {
	numberOfValid := 0
	for _, id := range s.ids {
		for _, r := range s.ranges {
			if id >= r.start && id <= r.end {
				numberOfValid++
				break
			}
		}
	}
	return numberOfValid, nil
}

func (s *Solver) Part2() (interface{}, error) {
	sort.Slice(s.ranges, func(i, j int) bool {
		return s.ranges[i].start < s.ranges[j].start
	})

	numberOfValid := 0
	currentStart := s.ranges[0].start
	currentEnd := s.ranges[0].end

	for _, r := range s.ranges[1:] {
		if r.start <= currentEnd+1 {
			if r.end > currentEnd {
				currentEnd = r.end
			}
		} else {
			numberOfValid += currentEnd - currentStart + 1
			currentStart = r.start
			currentEnd = r.end
		}
	}
	numberOfValid += currentEnd - currentStart + 1

	return numberOfValid, nil
}
