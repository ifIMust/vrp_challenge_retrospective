package conc_branch

import (
	"math"
	"sort"
	"sync"

	"github.com/ifIMust/vrp_challenge/common"
)

// Try to improve on the naive branch approach, by one or more of the following:
// Using goroutines to improve CPU utilization and (hopefully) pruning
// Using a sort or deterministic order for loads
// Memoizing the cost of traveling between 2 loads
// Improving bound function by incorporating minimum cost of remaining loads

type ConcurrentBranchBoundSearcher struct {
	loads []*common.Load

	bestMutex sync.RWMutex
	bestCost  float64
	bestRoute [][]int

	//semaphore chan int
	waitGroup sync.WaitGroup
	done      chan int
}

func NewConcurrentBranchBoundSearcher(loads []*common.Load) *ConcurrentBranchBoundSearcher {
	c := ConcurrentBranchBoundSearcher{}
	c.loads = loads
	c.bestCost = math.Inf(1)
	//c.semaphore = make(chan int, MaxGoroutines)
	c.done = make(chan int)
	return &c
}

func (c *ConcurrentBranchBoundSearcher) lowestCost() float64 {
	c.bestMutex.RLock()
	defer c.bestMutex.RUnlock()
	return c.bestCost
}

func (c *ConcurrentBranchBoundSearcher) postResult(cost float64, route [][]int) {
	c.bestMutex.Lock()
	defer c.bestMutex.Unlock()

	if cost < c.bestCost {
		c.bestCost = cost
		c.bestRoute = route
	}
}

func (c *ConcurrentBranchBoundSearcher) GetRoutes() [][]int {
	// Set up initial state / head of graph for first driver
	assignments := make([][]int, 0)
	assignments = append(assignments, make([]int, 0))
	remainingLoads := common.AsMap(c.loads)
	driver := 0

	c.search(remainingLoads, assignments, driver, common.HomeLocation, 0.0, 0.0)
	c.waitGroup.Wait()
	return c.bestRoute
}

type BranchTask struct {
	Load              *common.Load
	remainingLoads    common.LoadMap
	assignments       [][]int
	driver            int
	location          *common.Location
	driverMinutesUsed float64
	totalMinutesUsed  float64
}

func (c *ConcurrentBranchBoundSearcher) search(
	remainingLoads common.LoadMap,
	assignments [][]int,
	driver int,
	location *common.Location,
	driverMinutesUsed float64,
	totalMinutesUsed float64) {

	// Is all the work assigned for this branch?
	if len(remainingLoads) == 0 {
		// Account for sending the last driver home
		totalMinutesUsed += location.HomeCost

		c.postResult(totalMinutesUsed, assignments)
		return
	}

	sorter := common.NewLoadSorter(remainingLoads, location)
	sort.Sort(sorter)

	// Precompute overall cost remaining, regardless of next branch choice
	var minCost = 0.0
	var lowestHomeCost = math.Inf(1)
	for _, load := range remainingLoads {
		minCost += load.Cost
		lowestHomeCost = min(lowestHomeCost, load.HomeCost())
	}

	maxBranches := min(sorter.Len(), 1)
	for branch := 0; branch < maxBranches; branch += 1 {
		nearbyLoad := sorter.LoadEntries[branch].Load

		c.waitGroup.Add(1)
		//c.semaphore <- 1
		go func() {
			// Check if this branch should be considered
			if bound(nearbyLoad, totalMinutesUsed+minCost, location) < c.lowestCost() {
				// Duplicate these to avoid entanglement with other branches.
				remainingLoadsCopy := remainingLoads.Duplicate()
				assignmentsCopy := deepCopyAssigments(assignments)

				// Check if current driver can handle this Load
				if bound(nearbyLoad, driverMinutesUsed+nearbyLoad.HomeCost()+nearbyLoad.Cost, location) > common.MaxMinutesPerDriver {
					// Current driver can't do this load.
					// This branch will continue with a new driver starting at the depot location.
					c.search(remainingLoadsCopy,
						append(assignmentsCopy, make([]int, 0)),
						driver+1,
						common.HomeLocation,
						0.0,
						totalMinutesUsed+location.HomeCost)
				} else {
					// Assign this work to current driver
					assignmentsCopy[driver] = append(assignmentsCopy[driver], nearbyLoad.Index)
					delete(remainingLoadsCopy, nearbyLoad.Index)
					additionalMinutes := location.Distance(nearbyLoad.Pickup) + nearbyLoad.Cost

					c.search(remainingLoadsCopy,
						assignmentsCopy,
						driver,
						nearbyLoad.Dropoff,
						driverMinutesUsed+additionalMinutes,
						totalMinutesUsed+additionalMinutes)
				}
			}
			c.waitGroup.Done()
		}()
	}
}

func AssignRoutes(loads []*common.Load) [][]int {
	searcher := NewConcurrentBranchBoundSearcher(loads)
	return searcher.GetRoutes()
}

// Bounding function for branch pruning or driver capacity checks
// prevMinutes includes prior travel, plus minimum total travel needed for completion of driver route or total problem
func bound(load *common.Load, prevMinutes float64, location *common.Location) float64 {
	return prevMinutes + location.Distance(load.Pickup)
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
