package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ifIMust/vrp_challenge/greedy"
	"github.com/ifIMust/vrp_challenge/input"
	"github.com/ifIMust/vrp_challenge/tabu"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s input_file\n", os.Args[0])
		os.Exit(1)
	}

	loads, err := input.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Failed to read input file:", err.Error())
		os.Exit(2)
	}

	// Use Greedy algorithm to find a valid solution quickly
	assignments := greedy.AssignRoutes(loads)

	// Perform Tabu Search on the solution, to find an improved solution
	assignments = tabu.TabuSearch(assignments, loads)

	for _, driver := range assignments {
		fmt.Println(formatSlice(driver))
	}
}

// Modify the standard formatting for a slice, replacing commas with spaces
func formatSlice(slice []int) string {
	defaultFormat := fmt.Sprintf("%v", slice)
	return strings.ReplaceAll(defaultFormat, " ", ",")
}
