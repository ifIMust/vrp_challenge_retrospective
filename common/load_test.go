package common

import "testing"

func TestLoadHomeCostPickup(t *testing.T) {
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	load := NewLoad(1, l0x, l0y)

	hcp := load.HomeCostPickup()
	if 0 != hcp {
		t.Fatalf("Pickup Home Cost incorrect")
	}
}

func TestLoadHomeCostDropoff(t *testing.T) {
	l0x := NewLocation(0, 0)
	l0y := NewLocation(3, 4)
	load := NewLoad(1, l0x, l0y)

	hcp := load.HomeCostDropoff()
	if 5 != hcp {
		t.Fatalf("Dropoff Home Cost incorrect")
	}
}
