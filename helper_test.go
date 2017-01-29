package transformation

import (
	"math"
	"testing"
)

func checkDistance(t *testing.T, name string, expected, actual float64) {
	if !equalDistance(actual, expected) {
		t.Errorf("%s: expected %f, actual %f", name, expected, actual)
	}
}

func checkAngle(t *testing.T, name string, expected, actual float64) {
	if !equalAngle(actual, expected) {
		t.Errorf("%s: expected %f, actual %f", name, expected, actual)
	}
}

func checkRegion(t *testing.T, name string, expected, actual geoidRegion) {
	if expected != actual {
		t.Errorf("%s: expected %d, actual %d", name, expected, actual)
	}
}

func equalAngle(a, b float64) bool {
	const epsilon = 0.0000001
	return math.Abs(a-b) < epsilon
}

func equalDistance(a, b float64) bool {
	// check distances are equal to millimetre level
	const epsilon = 0.001
	return math.Abs(a-b) < epsilon
}
