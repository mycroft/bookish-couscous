//go:generate protoc -I ./ --go_out=plugins=grpc:. signplace.proto session.proto client.proto

package common

import (
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
