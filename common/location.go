package common

import "math"

type Location struct {
	X float64
	Y float64
}

func (l *Location) Distance(other *Location) float64 {
	return math.Sqrt(math.Pow(l.X-other.X, 2) + math.Pow(l.Y-other.Y, 2))
}

func (l *Location) HomeCost() float64 {
	return l.Distance(HomeLocation)
}
