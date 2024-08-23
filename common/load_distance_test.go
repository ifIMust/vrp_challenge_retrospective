package common

import (
	"math"
	"testing"
)

const floatCompareThreshold = 1e-7

func almostEqual(x float64, y float64) bool {
	return math.Abs(x-y) <= floatCompareThreshold
}

func TestLoadDistanceNew(t *testing.T) {
	loads := make([]*Load, 2)
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	loads[0] = NewLoad(1, l0x, l0y)

	l1x := NewLocation(0, 0)
	l1y := NewLocation(3, 4)
	loads[1] = NewLoad(2, l1x, l1y)

	ld := NewLoadDistance(loads)
	d0 := ld.Distance(1, 2)
	if d0 != 5.0 {
		t.Fatalf("wrong distance stored")
	}

	if d0 != ld.Distance(2, 1) {
		t.Fatalf("equivalent loads were not the same distance apart")
	}
}

func TestAverageDistance(t *testing.T) {
	loads := make([]*Load, 2)
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	loads[0] = NewLoad(1, l0x, l0y)

	l1x := NewLocation(0, 0)
	l1y := NewLocation(3, 4)
	loads[1] = NewLoad(2, l1x, l1y)

	ld := NewLoadDistance(loads)
	avg := ld.AverageDistance()

	if 5 != avg {
		t.Fatalf("average unexpected: %v", avg)
	}
}

func TestNearest(t *testing.T) {
	loads := make([]*Load, 3)
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	loads[0] = NewLoad(1, l0x, l0y)

	l1x := NewLocation(3, 5)
	l1y := NewLocation(6, 9)
	loads[1] = NewLoad(2, l1x, l1y)

	l2x := NewLocation(5, 5)
	l2y := NewLocation(1, 2)
	loads[2] = NewLoad(3, l2x, l2y)

	exclude := make(map[int]struct{})
	ld := NewLoadDistance(loads)
	nearestIndex, distance := ld.Nearest(1, exclude)
	if 2 != nearestIndex {
		t.Fatalf("nearest index lookup incorrect")
	}
	if distance != 1 {
		t.Fatalf("nearest distance incorrect")
	}

	nearestIndex, distance = ld.Nearest(2, exclude)
	if 3 != nearestIndex {
		t.Fatalf("nearest index lookup incorrect")
	}
	if !almostEqual(distance, math.Sqrt(17)) {
		t.Fatalf("nearest distance incorrect: %v", distance)
	}

	// Check exclusion
	exclude[3] = struct{}{}
	nearestIndex, distance = ld.Nearest(2, exclude)
	if 1 != nearestIndex {
		t.Fatalf("nearest index lookup incorrect")
	}
}

func TestRouteCost(t *testing.T) {
	loads := make([]*Load, 2)
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	loads[0] = NewLoad(1, l0x, l0y)

	l1x := NewLocation(3, 5)
	l1y := NewLocation(6, 9)
	loads[1] = NewLoad(2, l1x, l1y)

	ld := NewLoadDistance(loads)

	route := make([][]int, 1)
	route[0] = make([]int, 2)
	route[0][0] = 1
	route[0][1] = 2

	cost := ld.RouteCost(route)
	expectedCost := 500 + 5 + 1 + 5 + math.Sqrt(36+81)
	if !almostEqual(expectedCost, cost) {
		t.Fatalf("cost calculation failed. expected %v, got %v", expectedCost, cost)
	}
}
