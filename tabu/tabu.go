package tabu

import (
	"math"
	"reflect"
	"slices"

	"github.com/ifIMust/vrp_challenge/common"
)

type Route [][]int

// Driver routes this size or smaller will be selected for moving loads to other routes
const maxSourceRouteSize = 5

// Total Tabu search loops
const iterations = 76

// Size of Tabu list
const tabuSize = 20

// CandidateResult is the output of a concurrently evaluated candidate solution.
// If 'good' is true, the main function should check the score
// against the best candidate score.
type CandidateResult struct {
	score     float64
	candidate Route
	good      bool
}

func handleCandidate(candidate Route, bestCandidateScore float64, ld *common.LoadDistance, tabu []Route, resultChan chan CandidateResult) {
	var result CandidateResult
	result.candidate = candidate

	if isValid(candidate, ld) {
		result.score = ld.RouteCost(candidate)
		if result.score < bestCandidateScore && !isTabu(candidate, tabu) {
			result.good = true
		}
	}
	resultChan <- result
}

// Try to improve a solution by exploring similar solutions.
func TabuSearch(route Route, loads []*common.Load) Route {
	resultChan := make(chan CandidateResult)
	ld := common.NewLoadDistance(loads)

	bestScore := ld.RouteCost(route)
	bestSolution := deepCopyRoute(route)
	bestCandidate := bestSolution
	var candidates []Route = make([]Route, 1)
	candidates[0] = bestCandidate

	// The Tabu list tracks recent best candidates.
	// Avoiding candidates on this list helps to explore worse solutions
	// on the way to better ones.
	var tabu []Route = make([]Route, 0)

	for i := 0; i < iterations; i += 1 {
		candidates = getNeighbors(bestCandidate)
		bestCandidateScore := math.Inf(1)

		for _, c := range candidates {
			go handleCandidate(c, bestCandidateScore, ld, tabu, resultChan)
		}
		for ci := 0; ci < len(candidates); ci += 1 {
			result := <-resultChan
			if result.good && result.score < bestCandidateScore {
				bestCandidateScore = result.score
				bestCandidate = result.candidate
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

// isValid determines whether a solution violates the driver maximum distance constraint.
func isValid(route Route, ld *common.LoadDistance) bool {
	for _, driverRoute := range route {
		if ld.MinutesFromRoute(driverRoute, true) > common.MaxMinutesPerDriver {
			return false
		}
	}
	return true
}

// Focus on trying to insert single loads from smallest routes into other routes.
func getNeighbors(route Route) []Route {
	neighbors := make([]Route, 0)
	// Look through all driver routes for ones the right size
	for i, driverRoute := range route {
		driverRouteSz := len(driverRoute)
		if driverRouteSz <= maxSourceRouteSize {
			// This route is small enough to tamper with.
			// Try moving loads from the route to all positions in all other routes.
			for n, modifiedDriverRoute := range route {
				if i != n { // Don't move the load into the same route.
					for o := 0; o < len(modifiedDriverRoute)+1; o += 1 {
						for sourceRouteIdx := 0; sourceRouteIdx < driverRouteSz; sourceRouteIdx += 1 {
							neighbor := deepCopyRoute(route)
							// insert load at new position
							neighbor[n] = slices.Insert(neighbor[n], o, neighbor[i][sourceRouteIdx])

							// remove the load we just copied
							neighbor[i] = slices.Delete(neighbor[i], sourceRouteIdx, sourceRouteIdx+1)

							// remove entire previous driver route, if we took the last element
							if len(neighbor[i]) == 0 {
								neighbor = slices.Delete(neighbor, i, i+1)
							}
							neighbors = append(neighbors, neighbor)
						}
					}
				}
			}
		}
	}
	return neighbors
}

// Copy slice contents, to avoid having slices refer to same underlying array.
func deepCopyRoute(a Route) Route {
	result := make([][]int, 0, len(a))
	for _, v := range a {
		nested := make([]int, len(v))
		copy(nested, v)
		result = append(result, nested)
	}
	return result
}
