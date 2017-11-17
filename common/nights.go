package common

import (
	// "log"
	"sort"
	"time"
)

type timeSlice []time.Time

func (p timeSlice) Len() int           { return len(p) }
func (p timeSlice) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p timeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//
// Store last 3 nights in array.
// Also, it removes very old nights (> 7 days)
//
func Last3Nights(nights timeSlice) []time.Time {
	if len(nights) < 1 {
		return nights
	}

	weekago := time.Now().Add(time.Duration(-1 * time.Second * 86400 * 7))
	sort.Sort(nights)

	for {
		if nights[0].Before(weekago) || len(nights) > 3 {
			nights = nights[1:]
		} else {
			break
		}

		if len(nights) < 1 {
			break
		}
	}

	return nights
}

//
// Compute next night starting time (22pm)
//
func NextNight(t time.Time) time.Time {
	var nt time.Time
	if t.Hour() < 22 {
		nt = time.Date(
			t.Year(),
			t.Month(),
			t.Day(),
			22, 0, 0, 0,
			t.Location(),
		)
	} else {
		nt = time.Date(
			t.Year(),
			t.Month(),
			t.Day()+1,
			22, 0, 0, 0,
			t.Location(),
		)
	}

	return nt
}

//
// Determine if there is a night inside 2 time.Time.
// Night criteria:
// Must include a 6h duration between 22pm & 8am.
//
func IsNight(start time.Time, end time.Time) bool {
	if end.Before(start) {
		return false
	}

	duration := end.Sub(start)
	if duration.Hours() <= 6 {
		return false
	}

	// if start E [22:02[, then, ok, ...
	if (start.Hour() >= 22 && start.Hour() < 24) || (start.Hour() >= 0 && start.Hour() < 2) {
		return true
	}

	// ... else, we pick next night, and try again.
	return IsNight(NextNight(start), end)
}

/**
Old implementation:
//
// Is [start_ts, end_ts] a night ?
// if < 6h, not a night
// if start_ts > 22h && end_ts < 8h, ok
// if night start before 22h, it must stop after 4h of the morning (and start_hour > end_hour!)
// if night stop after 8h, it must start before 2am (and start_hour < end_hour!)
//
// Also, I know that doesn't cover most cases as I only manage session < 8 hours.
// It would be a little more complicated if we wanted to manage longer session (multiple days...)
//
// XXX TO REWRITE & TO TEST
func IsNight(start_ts uint64, end_ts uint64) bool {
	res := false
	duration := end_ts - start_ts
	if duration < 6*3600 {
		return res
	}

	start_hour := time.Unix(int64(start_ts), 0).Hour()
	end_hour := time.Unix(int64(end_ts), 0).Hour()

	if start_hour >= 22 && end_hour <= 8 {
		res = true
	}

	if start_hour < 22 && end_hour > 4 && start_hour > end_hour {
		res = true
	}

	if end_hour > 8 && start_hour < 2 && start_hour < end_hour {
		res = true
	}

	return res
}
**/
