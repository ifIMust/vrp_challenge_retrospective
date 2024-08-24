package tabu

import (
	"fmt"

	"github.com/ifIMust/vrp_challenge/common"
)

type Route [][]int

const iterations = 1

// Try to improve a solution by exploring similar solutions.
func TabuSearch(route Route, loads []*common.Load) Route {
	//ld := common.NewLoadDistance(loads)
	//bestScore := ld.RouteCost(route)
	bestSolution := deepCopyRoute(route)
	bestCandidate := bestSolution
	var candidates []Route = make([]Route, 1)
	candidates[0] = bestCandidate
	//var tabu []Route = make([]Route, 0)

	for i := 0; i < iterations; i += 1 {
		candidates = getNeighbors(bestCandidate)
	}
	return bestSolution
}

// Focus on trying to insert single load routes into other routes.
func getNeighbors(route Route) []Route {
	neighbors := make([]Route, 0)
	for i, driverRoute := range route {
		if len(driverRoute) == 1 {
			// Try moving this single load everywhere else
			for n, modifiedDriverRoute := range route {
				if i != n {
					for o, _ := range modifiedDriverRoute {
						_ = o
						neighbor := deepCopyRoute(route)
						// insert load at new position
						neighbor[n] = append(append(neighbor[n][0:o], route[i][0]), neighbor[n][o:(len(neighbor[n])-1)]...)
						// remove entire previous driver slot
						neighbor = append(neighbor[0:i], neighbor[i:len(neighbor)-1]...)
						neighbors = append(neighbors, neighbor)
					}
				}
			}
		}
	}
	fmt.Println("neighbors:")
	for i, n := range neighbors {
		fmt.Println(i, n)
	}

	return neighbors
}

// To avoid the nested slices from being entangled when branching, manually copy them
func deepCopyRoute(a Route) Route {
	result := make([][]int, 0, len(a))
	for _, v := range a {
		nested := make([]int, len(v))
		copy(nested, v)
		result = append(result, nested)
	}
	return result
}

// func deepCopyRoute(a Route) Route {
// 	result := make([][]int, 0, len(a))
// 	for _, v := range a {
// 		nested := make([]int, len(v))
// 		copy(nested, v)
// 		result = append(result, v)
// 	}
// 	return result
// }
