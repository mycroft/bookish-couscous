package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bsm/sarama-cluster"
	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"
)

var (
	kafkaHP  = "kafka:9092"
	redisHP  = "redis:6379"
	scyllaHP = "scylla"
)

func init() {
	flag.StringVar(&kafkaHP, "kafka", kafkaHP, "kafka host port")
	flag.StringVar(&redisHP, "redis", redisHP, "redis host port")
	flag.StringVar(&scyllaHP, "scylla", scyllaHP, "scylla host port")
}

func WaitForClusterSession(cluster *gocql.ClusterConfig, wait int) (*gocql.Session, error) {
	var session *gocql.Session
	var err error

	total_wait := 0

	for {
		session, err = cluster.CreateSession()
		if err != nil {
			if wait <= total_wait {
				time.Sleep(5)
				total_wait += 5
			}
		} else {
			break
		}
	}

	return session, err
}

func InitCqlCluster() *gocql.Session {
	cluster := gocql.NewCluster(scyllaHP)
	cluster.Keyspace = "system"

	session, err := WaitForClusterSession(cluster, 30)
	if err != nil {
		panic(err)
	}

	// Create keyspace.
	err = session.Query(
		`CREATE KEYSPACE IF NOT EXISTS zenly
		 WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };`).Exec()
	if err != nil {
		panic(err)
	}

	// Create table
	err = session.Query(
		`CREATE TABLE IF NOT EXISTS zenly.kyf (
		 user_id bigint,
		 rel_user_id bigint,
		 PRIMARY KEY(user_id, rel_user_id),
		 duration bigint,
		 week_most bigint,
		 week_friends bigint,
		 nights list<timestamp>,
		 week_most_list map<timestamp, int>,
		 week_friends_list map<timestamp, int>,
		);
		`,
	).Exec()

	if err != nil {
		panic(err)
	}

	session.Close()

	cluster.Keyspace = "zenly"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()

	return session
}

func main() {
	flag.Parse()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init cql
	cqlSession := InitCqlCluster()
	defer cqlSession.Close()

	// init consumer
	brokers := []string{kafkaHP}
	topics := []string{"sessions"}

	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	redis, err := redis.Dial("tcp", redisHP)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	processed := 0

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				session := &Session{}
				if err := proto.Unmarshal(msg.Value, session); err != nil {
					panic(err)
				}

				err := Process(cqlSession, redis, *session)
				if err != nil {
					panic(err)
				}

				consumer.MarkOffset(msg, "")

				processed++
				if processed%1000 == 0 {
					log.Printf("Processed %d messages.\n", processed)
				}
			}
		case <-signals:
			return
		}
	}

	return
}
