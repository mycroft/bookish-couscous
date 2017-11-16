package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.mkz.me/mycroft/bookish-couscous/common"

	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"
)

//
// Remove older elements and returns sum of all values
//
func GetDuration(m *map[time.Time]uint64) uint64 {
	var total uint64

	weekago := time.Now().Add(time.Duration(-1 * time.Second * 86400 * 7))
	for k, v := range *m {
		if weekago.Before(k) {
			total += v
		}
	}

	return total
}

func LoadCachedState(rc redis.Conn, uid uint32) (*HelloReply, error) {
	v, err := rc.Do("EXISTS", fmt.Sprintf("state:%d", uid))
	if err != nil {
		return nil, err
	}

	b, err := redis.Bool(v, err)
	if err != nil {
		return nil, err
	}

	if !b {
		return nil, nil
	}

	v, err = rc.Do("GET", fmt.Sprintf("state:%d", uid))
	if err != nil {
		return nil, err
	}

	out, err := redis.Bytes(v, err)
	if err != nil {
		return nil, err
	}

	helloReply := new(HelloReply)

	err = proto.Unmarshal(out, helloReply)
	if err != nil {
		return nil, err
	}

	return helloReply, nil
}

func CacheState(rc redis.Conn, uid uint32, hr *HelloReply) error {
	out, err := proto.Marshal(hr)
	if err != nil {
		return err
	}

	_, err = rc.Do("SET", fmt.Sprintf("state:%d", uid), out)
	if err != nil {
		return err
	}

	_, err = rc.Do("SETEX", fmt.Sprintf("state:%d", uid), 3600*12)

	return err
}

func GetFriends(cql *gocql.Session, rc redis.Conn, uid uint32, recursive bool) *HelloReply {
	var rel_user_id uint32
	var duration uint64
	week_friends_list := make(map[time.Time]uint64)
	week_most_list := make(map[time.Time]uint64)
	nights := make([]time.Time, 0)

	if recursive {
		log.Println("Wake up Neo, we are looking at uid", uid)

		// Check in redis
		state, err := LoadCachedState(rc, uid)
		if err == nil && state != nil {
			log.Println("Using cache value...")
			return state
		} else if err != nil {
			panic(err)
		}
	}

	helloReply := new(HelloReply)

	iter := cql.Query(
		`SELECT rel_user_id, week_friends_list, week_most_list, nights, duration
		 FROM kyf
		 WHERE user_id = ?`,
		uid,
	).Iter()

	var max_friends_duration, max_most_duration, max_most_all_time_duration uint64
	var max_friends_uid, max_most_uid, max_most_all_time_uid uint32

	helloReply.BestFriend = uid
	helloReply.MostSeen = uid
	helloReply.Crush = uid
	helloReply.MutualLove = uid
	helloReply.MutualLoveAllTime = uid

	// Do not compute more things if this doesn't change.
	max_most_uid, max_friends_uid = uid, uid

	for iter.Scan(&rel_user_id, &week_friends_list, &week_most_list, &nights, &duration) {
		// For each friends, compute ... Best Friend
		bf_duration := GetDuration(&week_friends_list)
		if bf_duration > max_friends_duration {
			max_friends_uid = rel_user_id
			max_friends_duration = bf_duration
		}

		// ... then most seen ...
		mo_duration := GetDuration(&week_most_list)
		if mo_duration > max_most_duration {
			max_most_uid = rel_user_id
			max_most_duration = mo_duration
		}

		// ... then crush, if any (or return uid, because I'm lazy) ...
		nights = common.Last3Nights(nights)
		if len(nights) >= 3 {
			// we've got a winner.
			helloReply.Crush = rel_user_id
		}

		// ... and we get a final value for mutual love all time.
		if duration > max_most_all_time_duration {
			max_most_all_time_uid = rel_user_id
			max_most_all_time_duration = duration
		}
	}

	helloReply.BestFriend = max_friends_uid
	helloReply.MostSeen = max_most_uid

	if recursive == false {
		helloReply.MutualLove = max_most_uid
		helloReply.MutualLoveAllTime = max_most_all_time_uid

		return helloReply
	}

	log.Println("best friend uid:", max_friends_uid, "duration:", max_friends_duration)
	log.Println("most seen uid:", max_most_uid, "duration:", max_most_duration)
	log.Println("most seen all time", max_most_all_time_uid, "duration:", max_most_all_time_duration)

	// it remains to compte the following:
	// mutual love
	// - we take max_most_uid, and we check if this is true for this user as well
	// mutual love all time
	// - we take max_most_all_time_uid, and we check if this is true for this user as well

	if max_most_uid == uid {
		// Let's skip this part.
		return helloReply
	}

	// we all calling the same function over those 2 users

	rel_most := GetFriends(cql, rc, max_most_uid, false)
	if rel_most.GetMutualLove() == uid {
		log.Println(uid, "found its mutual love", max_most_uid)
		helloReply.MutualLove = max_most_uid
	}

	if max_most_uid != max_most_all_time_uid {
		// do not recompute only if user is different... we never know!
		rel_most = GetFriends(cql, rc, max_most_all_time_uid, false)
	}

	if rel_most.GetMutualLoveAllTime() == uid {
		log.Println(uid, "found its all time mutual love", max_most_all_time_uid)
		helloReply.MutualLoveAllTime = max_most_all_time_uid
	}

	CacheState(rc, uid, helloReply)

	return helloReply
}

func aggregate(cql *gocql.Session, rc redis.Conn, uid uint32) *HelloReply {
	// Get all records for given uid
	return GetFriends(cql, rc, uid, true)
}
