package day11

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
)

type Solver struct {
	edges map[string][]string
}

func New() *Solver {
	input, _ := os.ReadFile(
		filepath.Join("internal", "solutions", "day11", "input.txt"),
	)

	edges := parseEdges(bytes.NewReader(input))

	return &Solver{edges: edges}
}

func parseEdges(r *bytes.Reader) map[string][]string {
	edges := make(map[string][]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) < 5 {
			continue
		}

		from := string(line[:3])
		rawTo := bytes.Split(line[5:], []byte{' '})

		to := make([]string, len(rawTo))
		for i, v := range rawTo {
			to[i] = string(v)
		}

		edges[from] = to
	}

	return edges
}

func dfsPart1(memo map[string]int, edges map[string][]string, cur string) int {
	if cur == "out" {
		return 1
	}

	if v, ok := memo[cur]; ok {
		return v
	}

	var res int
	for _, next := range edges[cur] {
		res += dfsPart1(memo, edges, next)
	}

	memo[cur] = res
	return res
}

func (s *Solver) Part1() (interface{}, error) {
	memo := make(map[string]int)
	return dfsPart1(memo, s.edges, "you"), nil
}

type serverState struct {
	name    string
	seenFFT bool
	seenDAC bool
}

func dfsPart2(
	memo map[serverState]int,
	edges map[string][]string,
	cur string,
	seenFFT, seenDAC bool,
) int {
	if cur == "out" {
		if seenFFT && seenDAC {
			return 1
		}
		return 0
	}

	state := serverState{cur, seenFFT, seenDAC}
	if v, ok := memo[state]; ok {
		return v
	}

	var res int
	for _, next := range edges[cur] {
		res += dfsPart2(
			memo,
			edges,
			next,
			seenFFT || next == "fft",
			seenDAC || next == "dac",
		)
	}

	memo[state] = res
	return res
}

func (s *Solver) Part2() (interface{}, error) {
	memo := make(map[serverState]int)
	return dfsPart2(memo, s.edges, "svr", false, false), nil
}
