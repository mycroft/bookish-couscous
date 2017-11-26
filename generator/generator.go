package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"gitlab.mkz.me/mycroft/bookish-couscous/common"
)

var (
	kafkaHP = "kafka:9092"
	redisHP = "redis:6379"

	topic  = "sessions"
	topics = []string{topic}

	max_events = 5000
	max_days   = 15

	users   = 100
	friends = 10
)

func init() {
	flag.StringVar(&kafkaHP, "kafka", kafkaHP, "kafka host port")
	flag.StringVar(&redisHP, "redis", redisHP, "redis host port")

	flag.IntVar(&max_events, "event", max_events, "Number of event to inject per day")
	flag.IntVar(&max_days, "days", max_days, "Number of days")

	flag.IntVar(&users, "users", users, "Number of users")
	flag.IntVar(&friends, "friends", friends, "Number of friends per user")
}

func newKafkaConfiguration() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.ChannelBufferSize = 1
	conf.Version = sarama.V0_10_1_0

	return conf
}

func newKafkaSyncProducer() sarama.SyncProducer {
	brokers := []string{kafkaHP}
	kafka, err := sarama.NewSyncProducer(brokers, newKafkaConfiguration())

	if err != nil {
		fmt.Printf("Kafka error: %s\n", err)
		os.Exit(-1)
	}

	return kafka
}

func sendMsg(kafka sarama.SyncProducer, event []byte) error {
	msgLog := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(event)),
	}

	_, _, err := kafka.SendMessage(msgLog)
	if err != nil {
		fmt.Printf("Kafka error: %s\n", err)
	}

	return nil
}

/*
Paris:
lat: 48.901741 / long: 2.261124 (saint ouen)
lat: 48.816811 / long: 2.405319 (ivry)

Coeur de Paris:
lat: 48.873522 / long: 2.326269 (bvd haussman)
lat: 48.846079 / long: 2.360687 (institut monde arabe)
*/
func generateRandomLocation() *common.SignPlace {
	sp := new(common.SignPlace)
	sp.Latitude = 48.816811 + rand.Float64()*(48.901741-48.816811)
	sp.Longitude = 2.261124 + rand.Float64()*(2.405319-2.261124)

	return sp
}

func main() {
	flag.Parse()

	kafka := newKafkaSyncProducer()

	redis, err := redis.Dial("tcp", redisHP)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

	// Always sure to have the same data sample (!)
	rand.Seed(42)

	num_users := users
	num_friends := friends

	rel := make([][]int, users)
	spt := make([]common.SignPlace, users)

	// Initialize users
	// (Create them an SP)
	for i := 0; i < num_users; i++ {
		// For each friends, we have 50 friends.
		friends := make([]int, num_friends)
		for j := 0; j < num_friends; j++ {
			new_friend := rand.Int() % num_users
			if new_friend == i {
				continue
			}
			friends[j] = new_friend

			redis.Do("SADD", fmt.Sprintf("friends:%d", i), friends[j])
			redis.Do("SADD", fmt.Sprintf("friends:%d", friends[j]), i)
		}
		rel[i] = friends

		// Generate a few of locations for this user (up to 3)
		num_location := 1 + (rand.Int() % 3)

		for j := 0; j < num_location; j++ {
			sp := generateRandomLocation()
			spt[i] = *sp
			out, err := proto.Marshal(sp)
			if err != nil {
				panic(err)
			}

			fmt.Println(out)

			redis.Do("SADD", fmt.Sprintf("loc:%d", i), out)
		}
	}

	// Choosing a Time between now and X days before.
	starting_ts := time.Now().Add(time.Hour * time.Duration(24*max_days))

	// starting_ts := uint64(time.Now().Unix() - int64(86400*max_days))
	injected_events := 0
	days := 0

	for {
		u1 := rand.Uint32() % uint32(num_users)

		start_ts := starting_ts.Add(time.Second * time.Duration(rand.Uint64()%86400))
		duration := time.Duration(time.Hour + time.Second*time.Duration((rand.Uint64()%(8*60))*60))

		loc := generateRandomLocation()
		if rand.Uint32()%10 == 0 {
			loc = &spt[u1]
		}

		pStartTs, err := ptypes.TimestampProto(start_ts)
		if err != nil {
			panic(err)
		}

		pEndTs, err := ptypes.TimestampProto(start_ts.Add(duration))
		if err != nil {
			panic(err)
		}

		p := common.Session{
			User1Id:   u1,
			User2Id:   uint32(rel[u1][rand.Int()%num_friends]),
			StartTs:   pStartTs,
			EndTs:     pEndTs,
			Latitude:  loc.GetLatitude(),
			Longitude: loc.GetLongitude(),
		}

		out, err := proto.Marshal(&p)
		if err != nil {
			panic(err)
		}

		// Write in kafka topic
		err = sendMsg(kafka, out)
		if err != nil {
			panic(err)
		}

		log.Println("Event sent.")

		injected_events++

		if injected_events == max_events {
			starting_ts = starting_ts.Add(time.Second * 86400)
			injected_events = 0
			days++
			if days == max_days {
				break
			}
		}
	}
}
