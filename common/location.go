package common

import "math"

// A location holds a Cartesian X,Y coordinate with a
// precomputed cost for traveling to or from the depot.
type Location struct {
	X        float64
	Y        float64
	HomeCost float64
}

// The Location constructor precomputes HomeCost so the Location is ready for use.
func NewLocation(x float64, y float64) *Location {
	l := Location{x, y, 0.0}
	l.HomeCost = l.Distance(HomeLocation)
	return &l
}

// Distance computes the Pythagorean distance between two Locations.
func (l *Location) Distance(other *Location) float64 {
	return math.Sqrt(math.Pow(l.X-other.X, 2) + math.Pow(l.Y-other.Y, 2))
}
