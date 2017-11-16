package main

import (
	// "log"
	"testing"
	"time"
)

func TestIsNight(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Error("Can't load UTC.")
	}

	v := IsNight(
		uint64(time.Date(2017, 11, 15, 23, 0, 0, 0, loc).Unix()),
		uint64(time.Date(2017, 11, 16, 12, 0, 0, 0, loc).Unix()),
	)
	if !v {
		t.Error("Should be a night.")
	}
}

func TestNights(t *testing.T) {
	var t1, t2, t3 time.Time
	nights := make([]time.Time, 0)

	// outdated
	nights = append(nights, time.Now().Add(-1*time.Hour*24*14))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*13))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*12))

	nights = Last3Nights(nights)
	if len(nights) != 0 {
		t.Error("Invalid nights len")
	}

	nights = append(nights, time.Now().Add(-1*time.Hour*24*6))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*5))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*4))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*3))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*2))

	nights = Last3Nights(nights)
	if len(nights) != 3 {
		t.Error("Invalid nights len")
	}

	nights = make([]time.Time, 0)

	nights = append(nights, time.Now().Add(-1*time.Hour*24*12))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*3))
	nights = append(nights, time.Now().Add(-1*time.Hour*24*2))

	nights = Last3Nights(nights)
	if len(nights) != 2 {
		t.Error("Invalid nights len")
	}

	// order
	nights = make([]time.Time, 0)

	t1 = time.Now().Add(-1 * time.Hour * 24 * 2)
	t2 = time.Now().Add(-1 * time.Hour * 24 * 3)
	t3 = time.Now().Add(-1 * time.Hour * 24 * 12)
	nights = append(nights, t1, t2, t3)

	nights = Last3Nights(nights)
	if len(nights) != 2 {
		t.Error("Invalid nights len")
	}

	if nights[0] != t2 || nights[1] != t1 {
		t.Error("Invalid dates")
	}
}

func TestIsNear(t *testing.T) {
	ref := &SignPlace{Latitude: 48.901741, Longitude: 2.261124}
	t1 := &SignPlace{Latitude: 48.901743, Longitude: 2.261128}
	t2 := &SignPlace{Latitude: 48.816811, Longitude: 2.405319}

	if !IsNear(ref, t1) {
		t.Error("t1 is not near from ref")
	}

	if IsNear(ref, t2) {
		t.Error("t2 near from ref but should not")
	}
}
