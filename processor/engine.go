package main

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"gitlab.mkz.me/mycroft/bookish-couscous/common"
)

type RelMetadata struct {
	uid1              uint32
	uid2              uint32
	duration          time.Duration
	nights            []time.Time
	week_most_list    map[time.Time]time.Duration
	week_most         time.Duration
	week_friends_list map[time.Time]time.Duration
	week_friends      time.Duration
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
		relMetadata.week_most_list = make(map[time.Time]time.Duration)
	}

	if relMetadata.week_friends_list == nil {
		relMetadata.week_friends_list = make(map[time.Time]time.Duration)
	}

	return relMetadata
}

//
// Save data to DB
//
func SaveRelMetadata(cql *gocql.Session, rm *RelMetadata) error {
	err := cql.Query(
		`INSERT INTO kyf(user_id, rel_user_id, duration, nights, week_most_list, week_friends_list)
		 VALUES(?, ?, ?, ?, ?, ?)`,
		rm.uid1,
		rm.uid2,
		rm.duration,
		rm.nights,
		rm.week_most_list,
		rm.week_friends_list,
	).Exec()

	if err != nil {
		log.Println(err)
	}

	return err
}

//
// Retrieve in redis SPs for a user
//
func LoadSPs(rc redis.Conn, uid uint32) ([]*common.SignPlace, error) {
	c, err := rc.Do("SMEMBERS", fmt.Sprintf("loc:%d", uid))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	vs, err := redis.Strings(c, err)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	locations := make([]*common.SignPlace, 0)

	for _, v := range vs {
		loc := new(common.SignPlace)
		if err := proto.Unmarshal([]byte(v), loc); err != nil {
			return nil, err
		}

		locations = append(locations, loc)
	}

	return locations, nil
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
func Process(cql *gocql.Session, rc redis.Conn, session common.Session) error {
	log.Printf("Got message...")

	// check if they are friends.
	c, err := rc.Do("SISMEMBER", fmt.Sprintf("friends:%d", session.GetUser1Id()), session.GetUser2Id())
	if v, err := redis.Bool(c, err); !v {
		// If not, we drop the data as it is irrelevant.
		return err
	}

	// Retrieve users data in DB
	relMetadata1 := GetRelMetadata(cql, session.GetUser1Id(), session.GetUser2Id())
	relMetadata2 := GetRelMetadata(cql, session.GetUser2Id(), session.GetUser1Id())

	// Compute session duration in minutes (we store minutes in db).
	end_ts, err := ptypes.Timestamp(session.GetEndTs())
	if err != nil {
		return err
	}

	start_ts, err := ptypes.Timestamp(session.GetStartTs())
	if err != nil {
		return err
	}

	session_duration := end_ts.Sub(start_ts) / time.Minute

	// In DB, we are storing using days data using midnight time for each day as a key,
	// so lets compute it...
	today := time.Date(
		end_ts.Year(),
		end_ts.Month(),
		end_ts.Day(),
		0, 0, 0, 0,
		end_ts.Location(),
	).UTC()

	// Storing in struct information to compute "Most seen" in last 7 days
	common.AddTimeTogether(relMetadata1.week_most_list, today, session_duration)
	common.AddTimeTogether(relMetadata2.week_most_list, today, session_duration)

	// Load SP (significant place)
	session_loc := &common.SignPlace{session.GetLatitude(), session.GetLongitude()}
	sps1, err := LoadSPs(rc, session.GetUser1Id())
	sps2, err := LoadSPs(rc, session.GetUser2Id())

	// Adding spend together out of our SP: "Best friend" feature
	if !common.IsNear(sps1, session_loc) {
		common.AddTimeTogether(relMetadata1.week_friends_list, today, session_duration)
	}

	// ... and we are doing this for both users.
	if !common.IsNear(sps2, session_loc) {
		common.AddTimeTogether(relMetadata2.week_friends_list, today, session_duration)
	}

	// Adding all time duration between those friends.
	// Storing in DB information to compute "mutual love"
	relMetadata1.duration += session_duration
	relMetadata2.duration = relMetadata1.duration

	// For crush, we make sure we are staying in the night,
	// SPs must NOT be the same
	// location must be either sp1 or sp2
	if true == common.IsNight(start_ts, end_ts) && !common.IsNearMultiple(sps1, sps2) && (common.IsNear(sps1, session_loc) || common.IsNear(sps1, session_loc)) {
		// Add night to night list
		relMetadata1.nights = common.Last3Nights(append(relMetadata1.nights, end_ts))
		relMetadata2.nights = relMetadata1.nights
	}

	SaveRelMetadata(cql, relMetadata1)
	SaveRelMetadata(cql, relMetadata2)

	// Drop cache for state:uid, if any
	_, err = rc.Do("DEL", fmt.Sprintf("state:%d", session.GetUser1Id()))
	if err != nil {
		return err
	}
	_, err = rc.Do("DEL", fmt.Sprintf("state:%d", session.GetUser2Id()))
	if err != nil {
		return err
	}

	return nil
}
