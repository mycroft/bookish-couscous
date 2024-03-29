//go:generate protoc -I ./ --go_out=plugins=grpc:. signplace.proto session.proto client.proto

package common

import (
	"time"
)

//
// Remove older elements and returns sum of all values
//
func CleanMap(m *map[time.Time]time.Duration) time.Duration {
	var total time.Duration

	weekago := time.Now().Add(time.Duration(-1 * time.Second * 86400 * 7))
	for k, v := range *m {
		if !weekago.Before(k) {
			delete(*m, k)
		} else {
			total += v
		}
	}

	return total
}

//
// Add time for given day (most seen)
// It will clean up obsolete data (> 7 days)
//
func AddTimeTogether(m map[time.Time]time.Duration, date time.Time, duration time.Duration) time.Duration {
	if _, ok := m[date]; ok {
		m[date] += duration
	} else {
		m[date] = duration
	}

	// Remove older elements & return total time passed together.
	return CleanMap(&m)
}
