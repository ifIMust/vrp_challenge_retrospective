package tabu

import (
	"math"
	"reflect"
	"slices"

	"github.com/ifIMust/vrp_challenge/common"
)

type Route [][]int

const iterations = 64
const tabuSize = iterations / 2

// Try to improve a solution by exploring similar solutions.
func TabuSearch(route Route, loads []*common.Load) Route {
	ld := common.NewLoadDistance(loads)
	bestScore := ld.RouteCost(route)
	bestSolution := deepCopyRoute(route)
	bestCandidate := bestSolution
	var candidates []Route = make([]Route, 1)
	candidates[0] = bestCandidate
	var tabu []Route = make([]Route, 0)

	for i := 0; i < iterations; i += 1 {
		candidates = getNeighbors(bestCandidate)
		bestCandidateScore := math.Inf(1)

		for _, c := range candidates {
			if !isTabu(c, tabu) && isValid(c, ld) {
				score := ld.RouteCost(c)

				if score < bestCandidateScore {
					bestCandidateScore = score
					bestCandidate = c
				}
			}
		}

		if bestCandidateScore == math.Inf(1) {
			break
		}

		if bestCandidateScore < bestScore {
			bestScore = bestCandidateScore
			bestSolution = bestCandidate
		}

		tabu = append(tabu, bestCandidate)
		if len(tabu) > tabuSize {
			tabu = tabu[1:]
		}

	}
	return bestSolution
}

func isTabu(route Route, tabu []Route) bool {
	for _, t := range tabu {
		if reflect.DeepEqual(route, t) {
			return true
		}
	}
	return false
}

func isValid(route Route, ld *common.LoadDistance) bool {
	for _, driverRoute := range route {
		if ld.MinutesFromRoute(driverRoute, true) > common.MaxMinutesPerDriver {
			return false
		}
	}
	return true
}

// Focus on trying to insert single load routes into other routes.
func getNeighbors(route Route) []Route {
	neighbors := make([]Route, 0)
	for i, driverRoute := range route {
		driverRouteSz := len(driverRoute)
		if driverRouteSz == 1 {
			// Try moving this single load everywhere else
			for n, modifiedDriverRoute := range route {
				if i != n {
					for o := 0; o < len(modifiedDriverRoute)+1; o += 1 {
						//for o, _ := range modifiedDriverRoute {
						neighbor := deepCopyRoute(route)
						// insert load at new position
						neighbor[n] = slices.Insert(neighbor[n], o, route[i][0])

						// remove entire previous driver slot
						neighbor = slices.Delete(neighbor, i, i+1)
						neighbors = append(neighbors, neighbor)
					}
				}
			}
		}
	}
	// fmt.Println("neighbors:")
	// for i, n := range neighbors {
	// 	fmt.Println(i, n)
	// }

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
