package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ifIMust/vrp_challenge/conc_branch"
	"github.com/ifIMust/vrp_challenge/input"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s input_file\n", os.Args[0])
		os.Exit(1)
	}

	loads := input.ReadFile(os.Args[1])

	//assignments := greedy.AssignRoutes(loads)
	//assignments := naive_branch.AssignRoutes(loads)
	assignments := conc_branch.AssignRoutes(loads)

	// output the results from the result structures
	for _, driver := range assignments {
		fmt.Println(formatSlice(driver))
	}
}

func formatSlice(slice []int) string {
	defaultFormat := fmt.Sprintf("%v", slice)
	return strings.ReplaceAll(defaultFormat, " ", ",")
}
