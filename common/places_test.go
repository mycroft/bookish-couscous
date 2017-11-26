package common

import "testing"

func TestIsNear(t *testing.T) {
	ref := &SignPlace{Latitude: 48.901741, Longitude: 2.261124}
	t1 := &SignPlace{Latitude: 48.901743, Longitude: 2.261128}
	t2 := &SignPlace{Latitude: 48.816811, Longitude: 2.405319}

	if !IsNearSingle(ref, t1) {
		t.Error("t1 is not near from ref")
	}

	if IsNearSingle(ref, t2) {
		t.Error("t2 near from ref but should not")
	}
}
