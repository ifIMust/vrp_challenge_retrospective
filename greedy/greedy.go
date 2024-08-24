package greedy

import (
	"sort"

	"github.com/ifIMust/vrp_challenge/common"
)

// Assign all loads by always using the closest location to the driver's current location.
func AssignRoutes(loads []*common.Load) [][]int {
	// assignments is the primary output.
	assignments := make([][]int, 0)

	// Each iteration, loads are read from remainingLoads and sorted.
	// Loads are deleted from the remainingLoads when assigned to a route.
	remainingLoads := common.AsMap(loads)

	// Used to check for task completion
	numLoads := len(loads)
	loadsCompleted := 0

	// driver is the driver currently being assigned
	for driver := 0; loadsCompleted < numLoads; driver += 1 {
		// Create a new empty route for the new driver
		assignments = append(assignments, make([]int, 0))
		// Assign nearby locations until the driver's day is full.
		loadsCompleted += greedy(remainingLoads, assignments, driver)
	}
	return assignments
}

// Assign the nearest location possible, as many times as possible, to this driver.
// Return the number of loads completed by this driver.
// remainingLoads and assignments are modified by this function.
func greedy(remainingLoads common.LoadMap,
	assignments [][]int,
	driver int) int {
	// Initial state for a new driver
	loadsCompleted := 0
	minutesUsed := 0.0
	location := common.HomeLocation

	for len(remainingLoads) > 0 {
		// Sort locations by pickup proximity to driver location
		sorter := common.NewLoadSorter(remainingLoads, location)
		sort.Sort(sorter)

		nextLoad := sorter.Pop()
		nextLoadCost := location.Distance(nextLoad.Pickup) + nextLoad.Cost
		// nextLoadMinCost includes returning to the depot.
		nextLoadMinCost := nextLoadCost + nextLoad.HomeCostDropoff()

		// Check if this driver has time to take the load.
		if nextLoadMinCost+minutesUsed > common.MaxMinutesPerDriver {
			// Not enough time.
			return loadsCompleted
		}

		// Assign the closest pickup to the driver.
		assignments[driver] = append(assignments[driver], nextLoad.Index)
		delete(remainingLoads, nextLoad.Index)
		minutesUsed += nextLoadCost
		location = nextLoad.Dropoff
		loadsCompleted += 1
	}
	return loadsCompleted
}
