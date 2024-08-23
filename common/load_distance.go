package common

// LoadDistance precomputes the distance from load dropoff to pickup
type LoadDistance struct {
	distances       map[int]map[int]float64
	averageDistance float64
}

func NewLoadDistance(loads []*Load) *LoadDistance {
	ld := LoadDistance{}
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

func (ld *LoadDistance) Distance(i, j int) float64 {
	return ld.distances[i][j]
}

func (ld *LoadDistance) AverageDistance() float64 {
	return ld.averageDistance
}
