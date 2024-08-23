package greedy

import (
	"sort"

	"github.com/ifIMust/vrp_challenge/common"
)

// Try to assign all loads using the closest location to the driver's current location
func AssignRoutes(loads []*common.Load) [][]int {
	// assignments is the primary output
	assignments := make([][]int, 0)

	remainingLoads := common.AsMap(loads)

	// Used to check for task completion
	numLoads := len(loads)
	loadsCompleted := 0

	// driver is the driver currently being assigned
	for driver := 0; loadsCompleted < numLoads; driver += 1 {
		assignments = append(assignments, make([]int, 0))
		greedy(remainingLoads, assignments, driver, common.HomeLocation, &loadsCompleted)
	}
	return assignments
}

// Assign the nearest location possible, as many times as possible, to this driver.
func greedy(remainingLoads common.LoadMap,
	assignments [][]int,
	driver int,
	location *common.Location,
	loadsCompleted *int) {

	minutesUsed := 0.0

	for len(remainingLoads) > 0 {
		// Sort remaining locations
		sorter := common.NewLoadSorter(remainingLoads, location)
		sort.Sort(sorter)

		nextLocation := sorter.Pop()
		nextLocationCost := location.Distance(nextLocation.Pickup) + nextLocation.Cost
		nextLocationMinCost := nextLocationCost + nextLocation.HomeCost()

		// Check if this driver's job is done.
		if nextLocationMinCost+minutesUsed > common.MaxMinutesPerDriver {
			return
		}

		// Assign the closest pickup to this driver
		assignments[driver] = append(assignments[driver], nextLocation.Index)
		delete(remainingLoads, nextLocation.Index)
		minutesUsed += nextLocationCost
		location = nextLocation.Dropoff
		*loadsCompleted = *loadsCompleted + 1
	}
}
