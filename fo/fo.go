//go:generate protoc -I ../ --go_out=plugins=grpc:. ../client.proto

package main

import (
	"flag"
	"log"
	"net"

	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	scyllaHP = "scylla"
	redisHP  = "redis:6379"
)

const (
	port = ":1980"
)

func init() {
	flag.StringVar(&scyllaHP, "scylla", scyllaHP, "scylla host port")
	flag.StringVar(&redisHP, "redis", redisHP, "redis host port")
}

type server struct {
	cql       *gocql.Session
	redisConn redis.Conn
}

func InitCqlCluster() *gocql.Session {
	cluster := gocql.NewCluster(scyllaHP)

	cluster.Keyspace = "zenly"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()

	return session
}

func GetRedis() redis.Conn {
	rc, err := redis.Dial("tcp", redisHP)
	if err != nil {
		panic(err)
	}

	return rc
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	helloReply := aggregate(s.cql, s.redisConn, in.GetUid())

	return helloReply, nil
}

func main() {
	flag.Parse()

	serv := &server{}

	cqlSession := InitCqlCluster()
	defer cqlSession.Close()

	rc := GetRedis()
	defer rc.Close()

	serv.cql = cqlSession
	serv.redisConn = rc

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterGreeterServer(s, serv)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
