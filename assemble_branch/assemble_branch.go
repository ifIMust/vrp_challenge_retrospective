package assemble_branch

import (
	"container/heap"
	"math"

	"github.com/ifIMust/vrp_challenge/common"
	"github.com/ifIMust/vrp_challenge/greedy"
)

// Try a branch and bound approach using precomputed nearest loads.
// The solutions space should expand in polynomial size rather than factorial
type BranchBoundSearcher struct {
	loads         []*common.Load
	loadMap       common.LoadMap
	bestCost      float64
	bestRoute     [][]int
	loadDistances *common.LoadDistance
}

func NewBranchBoundSearcher(loads []*common.Load) *BranchBoundSearcher {
	c := BranchBoundSearcher{}
	c.loads = loads
	c.loadMap = common.AsMap(loads)
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
	lowerBound       float64
	visited          map[int]struct{}
	assignments      [][]int
	driver           int
	totalMinutesUsed float64
	addFront         bool
	addBack          bool
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
		lowerBound:       c.lowestCost(),
		visited:          make(map[int]struct{}),
		assignments:      assignments,
		driver:           0,
		totalMinutesUsed: 0.0,
		addFront:         true,
		addBack:          true,
	}
	heap.Push(&queue, startNode)

	for queue.Len() > 0 {
		node := heap.Pop(&queue).(*SearchItem)

		if len(node.visited) == len(c.loads) {
			// candidate complete solution
			lastStopIndex := len(node.assignments[node.driver]) - 1
			node.totalMinutesUsed += c.loadMap[0].HomeCostPickup()
			node.totalMinutesUsed += c.loadMap[lastStopIndex].HomeCostDropoff()

			c.postResult(node.totalMinutesUsed, node.assignments)
		} else {
			// incomplete solution

			// If the current driver has no loads, choose the two loads with the least deadhead time
			// If those two loads can't be done in one day, a fallback plan is needed.
			if len(node.assignments[node.driver]) == 0 {
				drop, pick, _ := c.loadDistances.NearestPair(node.visited)
				if drop == 0 || pick == 0 {
					// Failed to assign a pair... only one left
					c.startRouteNearHome(node)
					continue
				}
				route := []int{drop, pick}
				if c.loadDistances.MinutesFromRoute(route, true) > common.MaxMinutesPerDriver {
					// The two nearest loads are sadly too much work for one driver.
					c.startRouteNearHome(node)
					continue
				}
				newNode := &SearchItem{}
				newNode.assignments = deepCopyAssigments(node.assignments)
				newNode.assignments[node.driver] = route
				newNode.visited = duplicateVisited(node.visited)
				newNode.visited[drop] = struct{}{}
				newNode.visited[pick] = struct{}{}
				newNode.driver = node.driver
				minutes := c.loadDistances.MinutesFromRoute(route, false)
				newNode.totalMinutesUsed = node.totalMinutesUsed + minutes
				newNode.addFront = true
				newNode.addBack = true

				// Estimate if this branch is worth considering
				newNode.lowerBound = c.bound(newNode)
				if newNode.lowerBound < c.lowestCost() {
					heap.Push(&queue, newNode)
				}
			} else {
				addedBack := false
				backLoad := 0
				addedFront := false
				frontLoad := 0
				// The driver already has a route started. Try to expand upon it.
				if node.addBack {
					driverRoute := node.assignments[node.driver]
					driverNumLoads := len(driverRoute)
					lastLoad := node.assignments[node.driver][driverNumLoads-1]
					nearestLoad, distance := c.loadDistances.Nearest(lastLoad, node.visited)
					driverRoute = append(driverRoute, nearestLoad)

					if c.loadDistances.MinutesFromRoute(driverRoute, true) <= common.MaxMinutesPerDriver {
						// The driver can handle the extra load
						newNode := &SearchItem{}
						newNode.assignments = deepCopyAssigments(node.assignments)
						newNode.assignments[node.driver] = driverRoute
						newNode.visited = duplicateVisited(node.visited)
						newNode.visited[nearestLoad] = struct{}{}
						newNode.driver = node.driver
						addedMinutes := distance + c.loadMap[nearestLoad].Cost
						newNode.totalMinutesUsed = node.totalMinutesUsed + addedMinutes
						newNode.addFront = false
						newNode.addBack = true

						newNode.lowerBound = c.bound(newNode)
						if newNode.lowerBound < c.lowestCost() {
							heap.Push(&queue, newNode)
							addedBack = true
							backLoad = nearestLoad
						}
					}
				}
				if node.addFront {
					driverRoute := node.assignments[node.driver]
					lastLoad := node.assignments[node.driver][0]
					nearestLoad, distance := c.loadDistances.NearestBefore(lastLoad, node.visited)
					driverRoute = append([]int{nearestLoad}, driverRoute...)

					if c.loadDistances.MinutesFromRoute(driverRoute, true) <= common.MaxMinutesPerDriver {
						// The driver can handle the extra load
						newNode := &SearchItem{}
						newNode.assignments = deepCopyAssigments(node.assignments)
						newNode.assignments[node.driver] = driverRoute
						newNode.visited = duplicateVisited(node.visited)
						newNode.visited[nearestLoad] = struct{}{}
						newNode.driver = node.driver
						addedMinutes := distance + c.loadMap[nearestLoad].Cost
						newNode.totalMinutesUsed = node.totalMinutesUsed + addedMinutes
						newNode.addFront = true
						newNode.addBack = false

						newNode.lowerBound = c.bound(newNode)
						if newNode.lowerBound < c.lowestCost() {
							heap.Push(&queue, newNode)
							addedFront = true
							frontLoad = nearestLoad
						}
					}
				}

				if addedBack && addedFront {
					// the third branch: continue with both
					driverRoute := node.assignments[node.driver]
					prevFront := driverRoute[0]
					prevBack := driverRoute[len(driverRoute)-1]
					driverRoute = append([]int{frontLoad}, driverRoute...)
					driverRoute = append(driverRoute, backLoad)

					if c.loadDistances.MinutesFromRoute(driverRoute, true) <= common.MaxMinutesPerDriver {
						// The driver can handle the extra loads
						newNode := &SearchItem{}
						newNode.assignments = deepCopyAssigments(node.assignments)
						newNode.assignments[node.driver] = driverRoute
						newNode.visited = duplicateVisited(node.visited)
						newNode.visited[frontLoad] = struct{}{}
						newNode.visited[backLoad] = struct{}{}
						newNode.driver = node.driver

						addedMinutes := c.loadMap[frontLoad].Cost + c.loadMap[backLoad].Cost
						addedMinutes += c.loadDistances.Distance(frontLoad, prevFront)
						addedMinutes += c.loadDistances.Distance(prevBack, backLoad)
						newNode.totalMinutesUsed = node.totalMinutesUsed + addedMinutes
						newNode.addFront = true
						newNode.addBack = true

						newNode.lowerBound = c.bound(newNode)
						if newNode.lowerBound < c.lowestCost() {
							heap.Push(&queue, newNode)
						}
					}
				} else if !addedBack && !addedFront {
					// This driver is done for the day. Add a node to start new driver's route
					newNode := &SearchItem{}
					newNode.assignments = deepCopyAssigments(node.assignments)
					newNode.assignments = append(newNode.assignments, make([]int, 0))
					newNode.driver = node.driver + 1
					newNode.visited = duplicateVisited(node.visited)
					newNode.totalMinutesUsed = node.totalMinutesUsed
					newNode.addFront = true
					newNode.addBack = true
					heap.Push(&queue, newNode)
				}
			}
		}
	}
	return c.bestRoute
}

func (c *BranchBoundSearcher) startRouteNearHome(node *SearchItem) {
	closeLoad := c.loadDistances.NearestPickupFromHome(node.visited)
	driverRoute := append(node.assignments[node.driver], closeLoad)
	newNode := &SearchItem{}
	newNode.assignments = deepCopyAssigments(node.assignments)
	newNode.assignments[node.driver] = driverRoute
	newNode.visited = duplicateVisited(node.visited)
	newNode.visited[closeLoad] = struct{}{}
	newNode.driver = node.driver
	minutes := c.loadDistances.MinutesFromRoute(driverRoute, false)
	newNode.totalMinutesUsed = node.totalMinutesUsed + minutes
	newNode.addFront = false
	newNode.addBack = true
}

// Bounding function for branch pruning
func (c *BranchBoundSearcher) bound(node *SearchItem) float64 {
	remainingLoadNum := len(c.loads) - len(node.visited)

	// Time cost remaining includes precomputed pickup->delivery cost for all loads,
	// travel time to next node, and lowest possible return to depot time.
	// Also includes time from depot to start node for thecurrent driver,
	// as it is not included until the end.
	currentDriverRouteStartTime := c.loadMap[node.assignments[node.driver][0]].HomeCostPickup()

	totalMinutesMinimum := node.totalMinutesUsed + currentDriverRouteStartTime + c.minRemainingMinutes(node.visited)

	avgDistancePerLoad := c.loadDistances.AverageDistance()

	// Add approximate travel time between all remaining loads
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
	var lowestHomeCost = math.Inf(1)
	for _, load := range c.loads {
		_, exclude := visited[load.Index]
		if !exclude {
			minCost += load.Cost
			lowestHomeCost = min(lowestHomeCost, load.HomeCostDropoff())
		}
	}
	minCost += lowestHomeCost
	return minCost
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
