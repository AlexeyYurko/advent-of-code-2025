package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/AlexeyYurko/advent-of-code-2025/internal/runner"
)

func main() {
	day := flag.Int("day", 0, "day to run (1-25)")
	flag.Parse()

	if *day < 1 || *day > 25 {
		log.Fatal("day must be between 1 and 25")
	}

	result, err := runner.Run(*day)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Day %d:\n", *day)
	fmt.Printf("Part 1: %v\n", result.Part1)
	fmt.Printf("Part 2: %v\n", result.Part2)
}
