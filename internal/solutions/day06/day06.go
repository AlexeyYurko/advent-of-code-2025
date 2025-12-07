package day06

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Range struct {
	start, end int
}

type Solver struct {
	n1parsed      []string
	n2parsed      []string
	n3parsed      []string
	n4parsed      []string
	opParsed      []string
	lines         []string
	numberOfChars int
}

func parseData(input string) (parsedString []string) {
	for _, data := range strings.Fields(input) {
		if data == "" {
			continue
		}
		parsedString = append(parsedString, data)
	}
	return parsedString
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day06", "input.txt"))
	base := strings.Split(string(input), "\n")

	n1 := base[0]
	n2 := base[1]
	n3 := base[2]
	n4 := base[3]
	op := base[4]
	return &Solver{
		n1parsed:      parseData(n1),
		n2parsed:      parseData(n2),
		n3parsed:      parseData(n3),
		n4parsed:      parseData(n4),
		opParsed:      parseData(op),
		lines:         base,
		numberOfChars: len(base[0]),
	}
}

func (s *Solver) Part1() (interface{}, error) {
	result := 0
	for i := 0; i < len(s.n1parsed); i++ {
		n1, _ := strconv.Atoi(s.n1parsed[i])
		n2, _ := strconv.Atoi(s.n2parsed[i])
		n3, _ := strconv.Atoi(s.n3parsed[i])
		n4, _ := strconv.Atoi(s.n4parsed[i])
		op := s.opParsed[i]
		if op == `+` {
			result += n1 + n2 + n3 + n4
		} else if op == `*` {
			result += n1 * n2 * n3 * n4
		}
	}
	return result, nil
}

func (s *Solver) Part2() (interface{}, error) {
	result := 0
	var numbers []int
	for i := s.numberOfChars - 1; i >= 0; i-- {
		number := fmt.Sprintf("%s%s%s%s", string(s.lines[0][i]), string(s.lines[1][i]), string(s.lines[2][i]), string(s.lines[3][i]))
		convertedNumber, _ := strconv.Atoi(strings.TrimSpace(number))
		if convertedNumber == 0 {
			continue
		}
		numbers = append(numbers, convertedNumber)

		opByte := s.lines[4][i]
		if opByte == '+' || opByte == '*' {
			var blockResult int
			if opByte == '+' {
				blockResult = 0
				for _, n := range numbers {
					blockResult += n
				}
			} else {
				blockResult = 1
				for _, n := range numbers {
					blockResult *= n
				}
			}
			result += blockResult

			numbers = numbers[:0]
		}
	}

	return result, nil
}
