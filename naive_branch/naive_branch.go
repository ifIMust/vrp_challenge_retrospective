package naive_branch

import (
	"math"
	"sort"

	"github.com/ifIMust/vrp_challenge/common"
)

// Try to improve on the greedy approach, by branching at location selection.
// Branch on whether each driver takes the next nearest load from the collection, or the second next.
// Prune branches that are worse than the best seen solution.
// Single threaded, with hopes of an improved concurrent version to follow.
func AssignRoutes(loads []*common.Load) [][]int {
	var bestRoute [][]int
	lowestCost := math.Inf(1)

	assignments := make([][]int, 0)
	assignments = append(assignments, make([]int, 0))
	remainingLoads := common.AsMap(loads)
	driver := 0

	search(remainingLoads, assignments, driver, common.HomeLocation, 0.0, 0.0, &bestRoute, &lowestCost)
	return bestRoute
}

func bound(load *common.Load, prevMinutes float64, location *common.Location) float64 {
	return prevMinutes + location.Distance(load.Pickup) + load.Cost() + load.HomeCost()
}

// To avoid the nested slices from being entangled when branching, manually copy them
func deepCopyAssigments(a [][]int) [][]int {
	result := make([][]int, 0, len(a))
	for _, v := range a {
		nested := make([]int, len(v))
		copy(nested, v)
		result = append(result, v)
	}
	return result

}

func search(
	remainingLoads common.LoadMap,
	assignments [][]int,
	driver int,
	location *common.Location,
	driverMinutesUsed float64,
	totalMinutesUsed float64,
	bestRoute *[][]int,
	lowestCost *float64) {

	// Is all the work assigned for this branch?
	if len(remainingLoads) == 0 {
		// Update best route if this is the best
		if totalMinutesUsed < *lowestCost {
			*lowestCost = totalMinutesUsed + location.HomeCost()
			*bestRoute = assignments

		}
		return
	}

	// Sort loads by those nearest to currrent location
	sorter := common.NewLoadSorter(remainingLoads, location)
	sort.Sort(sorter)

	//	for branch := 0; branch < sorter.Len(); branch += 1 {
	for branch := 0; branch < 1; branch += 1 { // TODO expand loop condition
		nearbyLoad := sorter.LoadEntries[branch].Load

		// Check if this branch should be considered
		if bound(nearbyLoad, totalMinutesUsed, location) < *lowestCost {
			// Duplicate these to avoid entanglement with other branches.
			remainingLoadsCopy := remainingLoads.Duplicate()
			assignmentsCopy := deepCopyAssigments(assignments)

			// Check if current driver can handle this Load
			if bound(nearbyLoad, driverMinutesUsed, location) > common.MaxMinutesPerDriver {
				// Current driver can't do this load.
				// This branch will continue with a new driver starting at the depot location.
				search(remainingLoadsCopy,
					append(assignmentsCopy, make([]int, 0)),
					driver+1,
					common.HomeLocation,
					0.0,
					totalMinutesUsed+location.HomeCost(),
					bestRoute,
					lowestCost)
			} else {
				// Assign this work to current driver
				assignmentsCopy[driver] = append(assignmentsCopy[driver], nearbyLoad.Index)
				delete(remainingLoadsCopy, nearbyLoad.Index)
				additionalMinutes := location.Distance(nearbyLoad.Pickup) + nearbyLoad.Cost()
				driverMinutesUsed += additionalMinutes
				totalMinutesUsed += additionalMinutes
				location = nearbyLoad.Dropoff

				search(remainingLoadsCopy,
					assignmentsCopy,
					driver,
					location,
					driverMinutesUsed,
					totalMinutesUsed,
					bestRoute,
					lowestCost)
			}
		}
	}
}
