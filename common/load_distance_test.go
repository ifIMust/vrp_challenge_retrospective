package common

import "testing"

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
