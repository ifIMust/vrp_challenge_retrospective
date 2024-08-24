package greedy

import (
	"sort"

	"github.com/ifIMust/vrp_challenge/common"
)

// Assign all loads by always using the closest location to the driver's current location.
func AssignRoutes(loads []*common.Load) ([][]int, float64) {
	// assignments is the primary output.
	assignments := make([][]int, 0)

	// Each iteration, loads are read from remainingLoads and sorted.
	// Loads are deleted from the remainingLoads when assigned to a route.
	remainingLoads := common.AsMap(loads)

	// Used to check for task completion
	numLoads := len(loads)
	loadsCompleted := 0

	minutesUsed := 0.0

	// driver is the driver currently being assigned
	for driver := 0; loadsCompleted < numLoads; driver += 1 {
		// Create a new empty route for the new driver
		assignments = append(assignments, make([]int, 0))
		// Assign nearby locations until the driver's day is full.
		minutesUsed += greedy(remainingLoads, assignments, driver, common.HomeLocation, &loadsCompleted)
	}
	return assignments, minutesUsed
}

// Assign the nearest location possible, as many times as possible, to this driver.
func greedy(remainingLoads common.LoadMap,
	assignments [][]int,
	driver int,
	location *common.Location,
	loadsCompleted *int) float64 {

	minutesUsed := 0.0

	for len(remainingLoads) > 0 {
		// Sort remaining locations
		sorter := common.NewLoadSorter(remainingLoads, location)
		sort.Sort(sorter)

		nextLoad := sorter.Pop()
		nextLoadCost := location.Distance(nextLoad.Pickup) + nextLoad.Cost
		nextLoadMinCost := nextLoadCost + nextLoad.HomeCostDropoff()

		// Check if this driver's job is done.
		if nextLoadMinCost+minutesUsed > common.MaxMinutesPerDriver {
			return minutesUsed
		}

		// Assign the closest pickup to this driver
		assignments[driver] = append(assignments[driver], nextLoad.Index)
		delete(remainingLoads, nextLoad.Index)
		minutesUsed += nextLoadCost
		location = nextLoad.Dropoff
		*loadsCompleted = *loadsCompleted + 1
	}
	return minutesUsed
}
