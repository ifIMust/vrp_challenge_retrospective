package common

import "math"

type Location struct {
	X        float64
	Y        float64
	HomeCost float64
}

func NewLocation(x float64, y float64) *Location {
	l := Location{x, y, 0.0}
	l.HomeCost = l.Distance(HomeLocation)
	return &l
}

func (l *Location) Distance(other *Location) float64 {
	return math.Sqrt(math.Pow(l.X-other.X, 2) + math.Pow(l.Y-other.Y, 2))
}
