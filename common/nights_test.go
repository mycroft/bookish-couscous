package common

import (
	"testing"
	"time"
)

func TestNoNight(t *testing.T) {
	nights := make([]time.Time, 0)
	nights = Last3Nights(nights)
	if len(nights) != 0 {
		t.Error("Invalid nights len")
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

func TestNextNight(t *testing.T) {
	var cts, nts, ts time.Time

	ts = time.Date(2017, time.November, 17, 21, 31, 0, 0, time.UTC)
	nts = NextNight(ts)
	cts = time.Date(2017, time.November, 17, 22, 0, 0, 0, time.UTC)
	if cts.String() != nts.String() {
		t.Error("Invalid next night time")
	}

	ts = time.Date(2017, time.November, 17, 22, 31, 0, 0, time.UTC)
	nts = NextNight(ts)
	cts = time.Date(2017, time.November, 18, 22, 0, 0, 0, time.UTC)
	if cts.String() != nts.String() {
		t.Error("Invalid next night time")
	}

	ts = time.Date(2017, time.December, 31, 22, 31, 0, 0, time.UTC)
	nts = NextNight(ts)
	cts = time.Date(2018, time.January, 1, 22, 0, 0, 0, time.UTC)
	if cts.String() != nts.String() {
		t.Error("Invalid next night time")
	}
}

func TestIsNight(t *testing.T) {
	var t1, t2 time.Time

	// valid nights
	t1 = time.Date(2017, time.November, 17, 22, 1, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 4, 2, 0, 0, time.UTC)
	if true != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	t1 = time.Date(2017, time.November, 18, 0, 30, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 7, 0, 0, 0, time.UTC)
	if true != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	// night start before 22pm
	t1 = time.Date(2017, time.November, 17, 17, 30, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 4, 1, 0, 0, time.UTC)
	if true != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	// night stops before 8am
	t1 = time.Date(2017, time.November, 18, 1, 59, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 9, 0, 0, 0, time.UTC)
	if true != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	// invalid nights
	t1 = time.Date(2017, time.November, 17, 22, 1, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 4, 0, 0, 0, time.UTC)
	if false != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	t1 = time.Date(2017, time.November, 18, 0, 1, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 18, 5, 59, 0, 0, time.UTC)
	if false != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}

	t1 = time.Date(2017, time.November, 17, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2017, time.November, 17, 11, 0, 0, 0, time.UTC)
	if false != IsNight(t1, t2) {
		t.Error("Invalid IsNight result")
	}
}
