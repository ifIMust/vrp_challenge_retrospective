package more_branch

import (
	"container/heap"
	"math"

	"github.com/ifIMust/vrp_challenge/common"
	"github.com/ifIMust/vrp_challenge/greedy"
)

// Try to improve on the naive branch approach, by one or more of the following:
// Using a min-heap and always processing the node with the lowest lower bound.
// Improving bound function by using the actual cost function, driver estimates,
// estimates travel between loads and minimum cost of remaining loads.
// Memoizing the cost of traveling between 2 loads (not utilized).
type BranchBoundSearcher struct {
	loads         []*common.Load
	bestCost      float64
	bestRoute     [][]int
	loadDistances *common.LoadDistance
}

func NewBranchBoundSearcher(loads []*common.Load) *BranchBoundSearcher {
	c := BranchBoundSearcher{}
	c.loads = loads
	c.loadDistances = common.NewLoadDistance(loads)

	var minutes float64

	// Quickly solve with "greedy" to serve as an upper bound on performance
	c.bestRoute, minutes = greedy.AssignRoutes(loads)
	c.bestCost = common.QuickCost(len(c.bestRoute), minutes)
	return &c
}

func (c *BranchBoundSearcher) lowestCost() float64 {
	return c.bestCost
}

func (c *BranchBoundSearcher) postResult(minutesDriven float64, route [][]int) {
	cost := common.QuickCost(len(route), minutesDriven)
	if cost < c.bestCost {
		c.bestCost = cost
		c.bestRoute = route
	}
}

// SearchItem is a partial or complete solution state used for the
// branch and bound strategy. It is used to sort solutions by lowest cost estimate
// in the priority queue.
type SearchItem struct {
	lowerBound        float64
	visited           map[int]struct{}
	assignments       [][]int
	driver            int
	location          *common.Location
	driverMinutesUsed float64
	totalMinutesUsed  float64
}

// LoadPriorityQueue implements the container.heap interface
type LoadPriorityQueue []*SearchItem

func NewLoadPriorityQueue() LoadPriorityQueue {
	return make([]*SearchItem, 0)
}

func (pq LoadPriorityQueue) Len() int {
	return len(pq)
}

func (pq LoadPriorityQueue) Less(i int, j int) bool {
	return pq[i].lowerBound < pq[j].lowerBound
}

func (pq LoadPriorityQueue) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *LoadPriorityQueue) Push(x any) {
	item := x.(*SearchItem)
	*pq = append(*pq, item)
}

func (pq *LoadPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // from docs: don't stop the GC from reclaiming the item eventually
	*pq = old[0 : n-1]
	return item
}

// AssignRoutes is the entry point to this algorithm
func AssignRoutes(loads []*common.Load) [][]int {
	searcher := NewBranchBoundSearcher(loads)
	return searcher.GetRoutes()
}

func (c *BranchBoundSearcher) GetRoutes() [][]int {
	// Set up initial state / head of graph for first driver
	assignments := make([][]int, 0)
	assignments = append(assignments, make([]int, 0))

	queue := NewLoadPriorityQueue()

	startNode := &SearchItem{
		lowerBound:        c.lowestCost(),
		visited:           make(map[int]struct{}),
		assignments:       assignments,
		driver:            0,
		location:          common.HomeLocation,
		driverMinutesUsed: 0.0,
		totalMinutesUsed:  0.0,
	}
	heap.Push(&queue, startNode)

	for queue.Len() > 0 {
		node := heap.Pop(&queue).(*SearchItem)

		if len(c.loads) == len(node.visited) {
			// candidate complete solution
			node.totalMinutesUsed += node.location.HomeCost
			c.postResult(node.totalMinutesUsed, node.assignments)
		} else {
			// partial solution
			// Search all possible branches from this point
			for _, load := range c.loads {
				_, visited := node.visited[load.Index]
				if visited {
					continue
				}

				// Set up a node representing the branch where this Load is handled next
				newNode := &SearchItem{}
				newNode.assignments = deepCopyAssigments(node.assignments)
				newNode.location = load.Dropoff
				newNode.visited = duplicateVisited(node.visited)
				newNode.visited[load.Index] = struct{}{}

				if driverTotalMinutesWithLoad(load, node.driverMinutesUsed, node.location) > common.MaxMinutesPerDriver {
					// New driver
					newNode.driver = node.driver + 1
					newNode.assignments = append(newNode.assignments, make([]int, 1))
					newNode.assignments[newNode.driver][0] = load.Index
					loadTime := load.HomeCostPickup() + load.Cost
					newNode.driverMinutesUsed = loadTime

					// Includes time for sending the last driver back to depot:
					newNode.totalMinutesUsed = node.totalMinutesUsed + node.location.HomeCost + loadTime
				} else {
					// Same driver
					newNode.driver = node.driver
					newNode.assignments[newNode.driver] = append(newNode.assignments[newNode.driver], load.Index)
					loadTime := node.location.Distance(load.Pickup) + load.Cost
					newNode.driverMinutesUsed = node.driverMinutesUsed + loadTime
					newNode.totalMinutesUsed = node.totalMinutesUsed + loadTime
				}

				// Estimate if this branch might be better than the best solution
				// minCost includes this load, so subtract its cost for more accurate bounding
				newNode.lowerBound = c.bound(newNode)
				if newNode.lowerBound < c.lowestCost() {
					heap.Push(&queue, newNode)
				}
			}
		}
	}
	return c.bestRoute
}

// Bounding function for branch pruning
func (c *BranchBoundSearcher) bound(node *SearchItem) float64 {
	minCostOfRemainingLoads := c.minRemainingMinutes(node.visited)
	remainingLoadNum := len(c.loads) - len(node.visited)

	// Time cost remaining includes precomputed pickup->delivery cost for all loads,
	// travel time to next node, and lowest possible return to depot time.
	totalMinutesMinimum := node.totalMinutesUsed + minCostOfRemainingLoads

	avgDistancePerLoad := c.loadDistances.AverageDistance()

	// Add approximate travel time for all remaining loads
	approxMinutes := totalMinutesMinimum + avgDistancePerLoad*float64(remainingLoadNum-1)

	// Estimate drivers needed for all remaining stops based on a heuristic minimum
	const maxLoadsPerDriver = 6.0
	const goodAvgLoadsPerDriver = 3.0

	var extraDriversNeeded int
	if remainingLoadNum > maxLoadsPerDriver {
		// many loads remain, use an average for the cost
		extraDriversNeeded = int(math.Floor(float64(remainingLoadNum) / goodAvgLoadsPerDriver))
	} else {
		// Near the end; don't bound too aggressively
		extraDriversNeeded = 0
	}

	totalDrivers := len(node.assignments) + extraDriversNeeded
	return common.QuickCost(totalDrivers, approxMinutes)
}

func (c *BranchBoundSearcher) minRemainingMinutes(visited map[int]struct{}) float64 {
	var minCost = 0.0
	for _, load := range c.loads {
		_, exclude := visited[load.Index]
		if !exclude {
			minCost += load.Cost
		}
	}
	return minCost
}

// Calculate the impact of adding a Load for this driver, to compare against the daily maximum
// Includes driving to the pickup, dropoff, and depot.
func driverTotalMinutesWithLoad(load *common.Load, prevMinutes float64, location *common.Location) float64 {
	return prevMinutes + location.Distance(load.Pickup) + load.Cost + load.HomeCostDropoff()
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

func duplicateVisited(source map[int]struct{}) map[int]struct{} {
	result := make(map[int]struct{})
	for i, j := range source {
		result[i] = j
	}
	return result
}
