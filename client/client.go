package main

import (
	"flag"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"gitlab.mkz.me/mycroft/bookish-couscous/common"
)

var (
	address = "fo:1980"
	uid     = 0
)

func init() {
	flag.StringVar(&address, "fo", "fo:1980", "fo host:port")
	flag.IntVar(&uid, "uid", 0, "uid to query")
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := common.NewGreeterClient(conn)

	// Contact the server and print out its response.
	log.Println("My uid is", uid)

	r, err := c.SayHello(context.Background(), &common.HelloRequest{Uid: uint32(uid)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Best friend:          %d\n", r.GetBestFriend())
	log.Printf("Crush:                %d\n", r.GetCrush())
	log.Printf("Most seen:            %d\n", r.GetMostSeen())
	if r.GetMutualLove() != uint32(uid) {
		log.Printf("Mutual love:          %d\n", r.GetMutualLove())
	}

	if r.GetMutualLoveAllTime() != uint32(uid) {
		log.Printf("Mutual love all time: %d\n", r.GetMutualLoveAllTime())
	}
}
