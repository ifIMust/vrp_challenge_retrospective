package tabu

import (
	"math"
	"reflect"
	"slices"
	"time"

	"github.com/ifIMust/vrp_challenge/common"
)

type Route [][]int

// Driver routes this size or smaller will be selected for moving loads to other routes
const maxSourceRouteSize = 5

// Total Tabu search loops
const iterations = 90

// Size of Tabu list
const tabuSize = 20

const timeLimitSeconds = 29 * time.Second

// CandidateResult is the output of a concurrently evaluated candidate solution.
// If 'good' is true, the main function should check the score
// against the best candidate score.
type CandidateResult struct {
	score     float64
	candidate Route
	good      bool
}

// Try to improve a solution by exploring similar solutions.
func TabuSearch(route Route, loads []*common.Load) Route {
	// Goroutines are used to process the candidate queue. When complete,
	// they send the result to resultChan.
	resultChan := make(chan CandidateResult)

	// The LoadDistance object is used to compute route costs.
	ld := common.NewLoadDistance(loads)

	bestScore := ld.RouteCost(route)
	bestSolution := deepCopyRoute(route)
	bestCandidate := bestSolution

	// The initial best candidate is the route given as input.
	var candidates []Route = make([]Route, 1)
	candidates[0] = bestCandidate

	// The Tabu list tracks recent best candidates.
	// Avoiding candidates on this list helps to explore worse solutions
	// on the way to better ones.
	var tabu []Route = make([]Route, 0)

	startTime := time.Now()
	for i := 0; i < iterations; i += 1 {
		// Stop if we are out of time.
		if time.Now().Sub(startTime) >= timeLimitSeconds {
			break
		}

		bestCandidateScore := math.Inf(1)

		// Generate the "neighboring space" of the previous best candidate.
		candidates = getNeighbors(bestCandidate)

		iterationComplete := make(chan struct{})

		// Read results from the result channel, and determine which,
		// if any, candidate was the best from the batch.
		go func() {
			for ci := 0; ci < len(candidates); ci += 1 {
				result := <-resultChan
				if result.good && result.score < bestCandidateScore {
					bestCandidateScore = result.score
					bestCandidate = result.candidate
				}
			}
			iterationComplete <- struct{}{}
		}()

		// Find valid candidates that are not on the Tabu list.
		// This creates a large, but finite, amount of goroutines to handle
		// the expensive tasks of computing the route costs, and performing a deep
		// comparison against the Tabu list.
		for _, candidate := range candidates {
			go func() {
				var result CandidateResult
				result.candidate = candidate

				if isValid(candidate, ld) {
					result.score = ld.RouteCost(candidate)
					if !isTabu(candidate, tabu) {
						result.good = true
					}
				}
				resultChan <- result
			}()
		}

		<-iterationComplete

		// If there were no valid, non-Tabu candidates, the search is complete.
		if bestCandidateScore == math.Inf(1) {
			break
		}

		// If we found a new best solution, save it.
		if bestCandidateScore < bestScore {
			bestScore = bestCandidateScore
			bestSolution = bestCandidate
		}

		// The best candidate is added to the Tabu list, so we will not consider it again soon.
		tabu = append(tabu, bestCandidate)

		// Remove old Tabu entries, to enforce a deliberately short memory.
		if len(tabu) > tabuSize {
			tabu = tabu[1:]
		}
	}
	return bestSolution
}

// A route that exactly matches any route on this list
// has already been explored recently.
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
		if ld.MinutesFromRoute(driverRoute) > common.MaxMinutesPerDriver {
			return false
		}
	}
	return true
}

// Focus on trying to insert single loads from smallest routes into other routes.
func getNeighbors(route Route) []Route {
	neighbors := make([]Route, 0)
	// Look through all driver routes for ones the right size
	for i, sourceRoute := range route {
		sourceRouteSz := len(sourceRoute)
		if sourceRouteSz <= maxSourceRouteSize {
			// This route is small enough to tamper with.
			// Try moving loads from the route to all positions in all other routes.
			for n, targetRoute := range route {
				// Iterate the target route to insert a load at every position
				for targetRouteIndex := 0; targetRouteIndex < len(targetRoute)+1; targetRouteIndex += 1 {
					// Iterate the source route to move a load from every position
					for sourceRouteIdx := 0; sourceRouteIdx < sourceRouteSz; sourceRouteIdx += 1 {
						if !(i == n && sourceRouteIdx == targetRouteIndex) {
							neighbor := deepCopyRoute(route)
							// insert load at new position
							neighbor[n] = slices.Insert(neighbor[n], targetRouteIndex, neighbor[i][sourceRouteIdx])

							// remove the load we just copied
							if i == n && sourceRouteIdx > targetRouteIndex {
								// When moving a load within the same route, the source index may change.
								neighbor[i] = slices.Delete(neighbor[i], sourceRouteIdx+1, sourceRouteIdx+2)
							} else {
								neighbor[i] = slices.Delete(neighbor[i], sourceRouteIdx, sourceRouteIdx+1)
							}
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
