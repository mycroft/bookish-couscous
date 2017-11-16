//go:generate protoc -I ../ --go_out=plugins=grpc:. ../client.proto

package main

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "fo:1980"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewGreeterClient(conn)

	// Contact the server and print out its response.
	uid := 0
	if len(os.Args) <= 1 {
		log.Fatalf("Please give an uid.")
	} else {
		uid, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("My uid is", uid)

	r, err := c.SayHello(context.Background(), &HelloRequest{Uid: uint32(uid)})
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
