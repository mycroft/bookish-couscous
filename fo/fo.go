//go:generate protoc -I ../ --go_out=plugins=grpc:. ../client.proto

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gocql/gocql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	scyllaHP = "scylla"
)

const (
	port = ":1980"
)

func init() {
	flag.StringVar(&scyllaHP, "scylla", scyllaHP, "scylla host port")
}

type server struct {
	cql *gocql.Session
}

func InitCqlCluster() *gocql.Session {
	cluster := gocql.NewCluster(scyllaHP)

	cluster.Keyspace = "zenly"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()

	return session
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	helloReply := aggregate(s.cql, in.GetUid())

	return helloReply, nil
}

func main() {
	flag.Parse()

	serv := &server{}
	fmt.Println("Hello world")

	cqlSession := InitCqlCluster()
	defer cqlSession.Close()

	serv.cql = cqlSession

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
