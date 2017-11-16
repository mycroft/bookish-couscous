package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.mkz.me/mycroft/bookish-couscous/common"

	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"github.com/golang/geo/s2"
	"github.com/golang/protobuf/proto"
)

type RelMetadata struct {
	uid1              uint32
	uid2              uint32
	duration          uint64
	nights            []time.Time
	week_most_list    map[time.Time]uint64
	week_most         uint64
	week_friends_list map[time.Time]uint64
	week_friends      uint64
}

//
// Load data from DB
//
func GetRelMetadata(cql *gocql.Session, u1 uint32, u2 uint32) *RelMetadata {
	relMetadata := new(RelMetadata)
	relMetadata.uid1 = u1
	relMetadata.uid2 = u2

	_ = cql.Query(
		`SELECT duration, nights, week_most_list, week_friends_list
		 FROM kyf WHERE user_id = ? AND rel_user_id = ?`,
		u1,
		u2,
	).Consistency(gocql.One).Scan(
		&relMetadata.duration,
		&relMetadata.nights,
		&relMetadata.week_most_list,
		&relMetadata.week_friends_list,
	)

	if relMetadata.week_most_list == nil {
		relMetadata.week_most_list = make(map[time.Time]uint64)
	}

	if relMetadata.week_friends_list == nil {
		relMetadata.week_friends_list = make(map[time.Time]uint64)
	}

	return relMetadata
}

//
// Save data to DB
//
func SaveRelMetadata(cql *gocql.Session, rm *RelMetadata) error {
	err := cql.Query(
		`INSERT INTO kyf(user_id, rel_user_id, duration, nights, week_most_list, week_friends_list, week_most, week_friends)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		rm.uid1,
		rm.uid2,
		rm.duration,
		rm.nights,
		rm.week_most_list,
		rm.week_friends_list,
		rm.week_most,
		rm.week_friends,
	).Exec()

	if err != nil {
		log.Println(err)
	}

	return err
}

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

//
// Is current position (e) is in same cell that given place (sp) ?
// As advised, it uses level 16 cells.
//
func IsNear(sp *SignPlace, e *SignPlace) bool {

	sp_latlon := s2.LatLngFromDegrees(sp.GetLatitude(), sp.GetLongitude())
	sp_cell := s2.CellFromLatLng(sp_latlon)

	parent_sp_cell_id := sp_cell.ID().Parent(16)
	parent_sp_cell := s2.CellFromCellID(parent_sp_cell_id)

	latlon := s2.LatLngFromDegrees(e.GetLatitude(), e.GetLongitude())
	p1 := s2.PointFromLatLng(latlon)

	return parent_sp_cell.ContainsPoint(p1)
}

//
// Remove older elements and returns sum of all values
//
func CleanMap(m *map[time.Time]uint64) uint64 {
	var total uint64

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
func AddTimeTogether(rm *RelMetadata, date time.Time, duration uint64) {
	if _, ok := rm.week_most_list[date]; ok {
		rm.week_most_list[date] += duration
	} else {
		rm.week_most_list[date] = duration
	}

	// Remove older elements.
	rm.week_most = CleanMap(&rm.week_most_list)

	return
}

//
// Add time for given day (friends)
// It will clean up obsolete data (> 7 days)
//
func AddTimeTogetherFriends(rm *RelMetadata, date time.Time, duration uint64) {
	if _, ok := rm.week_friends_list[date]; ok {
		rm.week_friends_list[date] += duration
	} else {
		rm.week_friends_list[date] = duration
	}

	// Remove older elements.
	rm.week_friends = CleanMap(&rm.week_friends_list)
}

//
// Retrieve in redis SP for a user
//
func LoadSP(rc redis.Conn, uid uint32) (*SignPlace, error) {
	c, err := rc.Do("GET", fmt.Sprintf("loc:%d", uid))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	v, err := redis.String(c, err)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	loc := &SignPlace{}
	if err := proto.Unmarshal([]byte(v), loc); err != nil {
		return nil, err
	}

	return loc, err
}

//
// Process a session
// - check users are friends
// - fetch metadatadata from db
// - fill metadata for "most seen" feature
// - load SP for both users
// - for both user,  fill metadata for "best friends" feature
// - fill "mutual love" metadata
// - if night & location requirements validated, fill metadata for "crush" feature
// - save metadatas to db
//
// Note at this point we did not compute best friend, crush, etc. We just prepared
// data needed to do it according to given session
//
func Process(cql *gocql.Session, rc redis.Conn, session Session) error {
	// check if they are friends.
	c, err := rc.Do("SISMEMBER", fmt.Sprintf("friends:%d", session.GetUser1Id()), session.GetUser2Id())
	if v, err := redis.Bool(c, err); !v {
		// If not, we drop the data as it is irrelevant.
		return err
	}

	// Retrieve users data in DB
	relMetadata1 := GetRelMetadata(cql, session.GetUser1Id(), session.GetUser2Id())
	relMetadata2 := GetRelMetadata(cql, session.GetUser2Id(), session.GetUser1Id())

	// Compute session duration in minutes.
	session_duration := (session.GetEndTs() - session.GetStartTs()) / 60

	// In DB, we are storing using days data using midnight time for each day as a key,
	// so lets compute it...
	today := time.Unix(int64(session.EndTs-(session.EndTs%86400)), 0).UTC()

	// Storing in struct information to compute "Most seen" in last 7 days
	AddTimeTogether(relMetadata1, today, session_duration)
	AddTimeTogether(relMetadata2, today, session_duration)

	// Load SP (significant place)
	session_loc := &SignPlace{session.GetLatitude(), session.GetLongitude()}
	sp1, err := LoadSP(rc, session.GetUser1Id())
	sp2, err := LoadSP(rc, session.GetUser2Id())

	// Adding spend together out of our SP: "Best friend" feature
	if !IsNear(sp1, session_loc) {
		AddTimeTogetherFriends(relMetadata1, today, session_duration)
	}

	// ... and we are doing this for both users.
	if !IsNear(sp2, session_loc) {
		AddTimeTogetherFriends(relMetadata2, today, session_duration)
	}

	// Adding all time duration between those friends.
	// Storing in DB information to compute "mutual love"
	relMetadata1.duration += session_duration
	relMetadata2.duration = relMetadata1.duration

	// For crush, we make sure we are staying in the night,
	// SPs must NOT be the same
	// location must be either sp1 or sp2
	if true == IsNight(session.GetStartTs(), session.GetEndTs()) && !IsNear(sp1, sp2) && (IsNear(sp1, session_loc) || IsNear(sp2, session_loc)) {
		// Add night to night list
		t := time.Unix(int64(session.GetEndTs()), 0)

		relMetadata1.nights = common.Last3Nights(append(relMetadata1.nights, t))
		relMetadata2.nights = relMetadata1.nights
	}

	SaveRelMetadata(cql, relMetadata1)
	SaveRelMetadata(cql, relMetadata2)

	return nil
}
