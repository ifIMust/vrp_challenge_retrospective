package common

// LoadDistance precomputes the distance from load dropoff to pickup
type LoadDistance struct {
	loads     LoadMap
	distances map[int]map[int]float64
}

func NewLoadDistance(loads []*Load) *LoadDistance {
	ld := LoadDistance{}
	ld.loads = AsMap(loads)
	ld.distances = make(map[int]map[int]float64)
	if len(loads) == 0 {
		return &ld
	}

	for _, li := range loads {
		ld.distances[li.Index] = make(map[int]float64)
		for _, lj := range loads {
			if li != lj {
				d := li.Dropoff.Distance(lj.Pickup)
				ld.distances[li.Index][lj.Index] = d
			}
		}
	}
	return &ld
}

// Return precomputed distance from i's Dropoff to j's Pickup
func (ld *LoadDistance) Distance(i, j int) float64 {
	return ld.distances[i][j]
}

// Use precomputed distances to determine the route cost
func (ld *LoadDistance) RouteCost(routes [][]int) float64 {
	drivers := 0
	minutes := 0.0
	for _, driver := range routes {
		drivers += 1
		minutes += ld.MinutesFromRoute(driver, true)
	}
	return float64(drivers)*500.0 + minutes
}

func (ld *LoadDistance) MinutesFromRoute(route []int, includeDepotTime bool) float64 {
	minutes := 0.0
	numLoads := len(route)
	if numLoads > 0 {
		if includeDepotTime {
			minutes += ld.loads[route[0]].HomeCostPickup()
			minutes += ld.loads[route[numLoads-1]].HomeCostDropoff()
		}
		lastLoad := route[0]
		for _, l := range route {
			minutes += ld.loads[l].Cost
			if l != lastLoad {
				minutes += ld.Distance(lastLoad, l)
			}
			lastLoad = l
		}
	}
	return minutes
}
