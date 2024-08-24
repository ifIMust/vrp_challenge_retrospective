package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ifIMust/vrp_challenge/input"
	"github.com/ifIMust/vrp_challenge/more_branch"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s input_file\n", os.Args[0])
		os.Exit(1)
	}

	loads := input.ReadFile(os.Args[1])

	assignments := more_branch.AssignRoutes(loads)
	//assignments := assemble_branch.AssignRoutes(loads)

	//assignments = tabu.TabuSearch(assignments, loads)

	// output the results from the result structures
	for _, driver := range assignments {
		fmt.Println(formatSlice(driver))
	}
}

func formatSlice(slice []int) string {
	defaultFormat := fmt.Sprintf("%v", slice)
	return strings.ReplaceAll(defaultFormat, " ", ",")
}
