package main

import (
	"fmt"
	"os"

	"github.com/ifIMust/vrp_challenge/input"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s input_file\n", os.Args[0])
		os.Exit(1)
	}

	// read the input data
	loads := input.ReadFile(os.Args[1])
	fmt.Println("Found loads: ", len(loads))

	// call one of the algorithms
	// output the results from the result structures

}
