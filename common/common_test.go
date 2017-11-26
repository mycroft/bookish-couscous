package common

import (
	"testing"
	"time"
)

func TestAddTimeTogether(t *testing.T) {
	var t1, t2, t3, t4 time.Time
	var v time.Duration

	m := make(map[time.Time]time.Duration)

	t1 = time.Now().Add(-1 * time.Hour * 24 * 14)
	t2 = time.Now().Add(-1 * time.Hour * 24 * 13)
	t3 = time.Now().Add(-1 * time.Hour * 24 * 13)

	v = AddTimeTogether(m, t1, 1)
	if v != 0 || len(m) != 0 {
		t.Error("Map is not empty (1)")
	}

	_ = AddTimeTogether(m, t1, 1)
	_ = AddTimeTogether(m, t2, 2)
	v = AddTimeTogether(m, t3, 3)
	if v != 0 || len(m) != 0 {
		t.Error("Map is not empty (1)")
	}

	t3 = time.Now().Add(-1 * time.Hour * 24 * 2)
	t4 = time.Now().Add(-1 * time.Hour * 24 * 1)

	_ = AddTimeTogether(m, t1, 1)
	_ = AddTimeTogether(m, t2, 2)
	_ = AddTimeTogether(m, t3, 3)
	v = AddTimeTogether(m, t4, 4)
	if v != (3+4) || len(m) != 2 {
		t.Error("V != 7 || len(m) != 2")
	}

	m = make(map[time.Time]time.Duration)

	_ = AddTimeTogether(m, t1, 1)
	_ = AddTimeTogether(m, t2, 2)
	_ = AddTimeTogether(m, t3, 3)
	_ = AddTimeTogether(m, t4, 4)
	v = AddTimeTogether(m, t4, 4)
	if v != (3+4+4) || len(m) != 2 {
		t.Error("V != 10 || len(m) != 2")
	}
}
