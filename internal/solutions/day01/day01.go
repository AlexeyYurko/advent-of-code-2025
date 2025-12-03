package day01

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Solver struct {
	input string
	steps []Steps
}

type Steps struct {
	direction byte
	turns     int
}

func New() *Solver {
	input, _ := os.ReadFile(filepath.Join("internal", "solutions", "day01", "input.txt"))
	s := &Solver{input: string(input)}
	s.steps, _ = prepareData(s.input)
	return s
}

func prepareData(s string) ([]Steps, error) {
	var openingSteps []Steps

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		direction := line[0]
		turns, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("error converting line: %s", line)
		}
		openingSteps = append(openingSteps, Steps{direction: direction, turns: turns})
	}
	return openingSteps, nil
}

func (s *Solver) Part1() (interface{}, error) {
	position := 50
	password := 0
	for _, step := range s.steps {
		if step.direction == 'L' {
			position -= step.turns
		} else {
			position += step.turns
		}
		position %= 100
		if position == 0 {
			password++
		}
	}
	return password, nil
}

func countHits(start int, dir byte, steps int) int {
	if steps <= 0 {
		return 0
	}
	if dir == 'R' {
		return (start + steps) / 100
	}
	if start == 0 {
		return steps / 100
	}
	if steps < start {
		return 0
	}
	return 1 + (steps-start)/100
}

func (s *Solver) Part2() (interface{}, error) {
	position := 50
	password := 0

	for _, step := range s.steps {
		password += countHits(position, step.direction, step.turns)

		if step.direction == 'R' {
			position = (position + step.turns) % 100
		} else {
			position = (position - step.turns%100 + 100) % 100
		}
	}
	return password, nil
}
