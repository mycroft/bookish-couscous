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
