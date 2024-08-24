package tabu

import "testing"

func TestDeepCopyRoute(t *testing.T) {
	origRoutes := make([][]int, 0)
	origRoutes = append(origRoutes, make([]int, 2))
	origRoutes[0][0] = 5
	origRoutes[0][1] = 6

	dup := deepCopyRoute(origRoutes)
	dup[0][0] = 4

	if origRoutes[0][0] != 5 {
		t.Fatalf("origRoutes modified: now reads: %v", origRoutes[0][0])
	}
}
