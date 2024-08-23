package common

func QuickCost(numDrivers int, minutesDriven float64) float64 {
	return float64(numDrivers)*500.0 + minutesDriven
}
