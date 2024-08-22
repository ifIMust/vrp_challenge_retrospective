package common

import (
	"math"
	"testing"
)

func TestDistancePerfect(t *testing.T) {
	l0 := Location{0, 3}
	l1 := Location{4, 0}
	d := l0.Distance(&l1)
	if d != 5.0 {
		t.Fatalf("Distance: %f", d)
	}
}

func TestDistanceZero(t *testing.T) {
	l0 := Location{3, 3}
	l1 := Location{3, 3}
	d := l0.Distance(&l1)
	if d != 0.0 {
		t.Fatalf("Distance: %f", d)
	}
}

func TestDistanceOne(t *testing.T) {
	l0 := Location{3, 3}
	l1 := Location{4, 4}
	d := l0.Distance(&l1)
	if d != math.Sqrt(2) {
		t.Fatalf("Distance: %f", d)
	}
}
