package common

import "math"

// LoadDistance precomputes the distance from load dropoff to pickup
type LoadDistance struct {
	loads           LoadMap
	distances       map[int]map[int]float64
	averageDistance float64
}

func NewLoadDistance(loads []*Load) *LoadDistance {
	ld := LoadDistance{}
	ld.loads = AsMap(loads)
	ld.distances = make(map[int]map[int]float64)
	if len(loads) == 0 {
		return &ld
	}

	var distanceSum float64
	var numDistances int
	for _, li := range loads {
		ld.distances[li.Index] = make(map[int]float64)
		for _, lj := range loads {
			if li != lj {
				d := li.Dropoff.Distance(lj.Pickup)
				ld.distances[li.Index][lj.Index] = d
				distanceSum += d
				numDistances += 1
			}
		}
	}
	ld.averageDistance = distanceSum / float64(numDistances)
	return &ld
}

// Return precomputed distance from i's Dropoff to j's Pickup
func (ld *LoadDistance) Distance(i, j int) float64 {
	return ld.distances[i][j]
}

func (ld *LoadDistance) AverageDistance() float64 {
	return ld.averageDistance
}

// Average distance between loads not excluded.
func (ld *LoadDistance) AverageDistanceRemaining(exclude map[int]struct{}) float64 {
	numDistances := 0
	distanceSum := 0.0

	for i, nestedMap := range ld.distances {
		_, iExcluded := exclude[i]
		if !iExcluded {
			for j, distance := range nestedMap {
				_, jExcluded := exclude[j]
				if !jExcluded {
					distanceSum += distance
					numDistances += 1
				}
			}
		}
	}

	if distanceSum == 0 {
		return 0.0
	}
	return distanceSum / float64(numDistances)
}

// Find the Load nearest to each the given index (Dropoff nearest to Pickup).
// Exclude all indices found in exclude.
// Return the index and the distance.
// Prefers valid input. Returns 0 index if input is invalid
func (ld *LoadDistance) Nearest(index int, exclude map[int]struct{}) (int, float64) {
	distance := math.Inf(1)
	nearest := 0
	m, ok := ld.distances[index]
	if !ok {
		return nearest, distance
	}

	for i, d := range m {
		_, excluded := exclude[i]
		if !excluded && d < distance {
			distance = d
			nearest = i
		}
	}

	return nearest, distance
}

// Like Nearest, but look for Load with closest Dropoff to the indexed Load's Pickup
func (ld *LoadDistance) NearestBefore(index int, exclude map[int]struct{}) (int, float64) {
	distance := math.Inf(1)
	nearest := 0
	for i, nestedMap := range ld.distances {
		_, excluded := exclude[i]
		if !excluded {
			d := nestedMap[index]
			if d < distance {
				distance = d
				nearest = i
			}
		}
	}

	return nearest, distance
}

// Find the pair of Loads nearest each other, excluding indices found in exclude.
// Return the indices (in order) and their distance.
// May return 0 indices if no pair is found (due to exclusion)
func (ld *LoadDistance) NearestPair(exclude map[int]struct{}) (int, int, float64) {
	minDistance := math.Inf(1)
	nearestDrop := 0
	nearestPick := 0

	for i, nestedMap := range ld.distances {
		_, iExcluded := exclude[i]
		if !iExcluded {
			for j, distance := range nestedMap {
				_, jExcluded := exclude[j]
				if !jExcluded && distance < minDistance {
					minDistance = distance
					nearestDrop = i
					nearestPick = j
				}
			}
		}
	}
	return nearestDrop, nearestPick, minDistance
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

func (ld *LoadDistance) NearestPickupFromHome(exclude map[int]struct{}) int {
	index := 0
	closest := math.Inf(1)
	for i, l := range ld.loads {
		_, excluded := exclude[i]
		if !excluded {
			if l.HomeCostPickup() < closest {
				closest = l.HomeCostPickup()
				index = i
			}
		}
	}
	return index
}
