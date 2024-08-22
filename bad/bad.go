package bad

import "github.com/ifIMust/vrp_challenge/input"

// This package provides an algorithm to solve the VRP challenge very badly.
// It's goal is to be even easier to write than brute-force, and establish the interface.

func AssignRoutes(loads []*input.Load) [][]int {
	assignments := make([][]int, 0, 1)

	driver0 := make([]int, 1)
	for _, load := range loads {
		driver0 = append(driver0, load.Index)
	}
	assignments = append(assignments, driver0)
	return assignments
}
