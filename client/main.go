package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/alisa-vernigor/grpc-go-reverse/pkg/proto/chat"
	"google.golang.org/grpc"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewChatRoomClient(conn)

	var clientName string
	// Taking input from user
	fmt.Scanln(&clientName)

	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var f bool = false
	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				print("EOF")
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message : %v", err)
			}
			log.Printf("Got message: %s", in.Name)
		}
	}()

	for {
		if !f {
			if err := stream.Send(&chat.ChatRequest{Name: clientName}); err != nil {
				log.Fatal(err)
			}
			f = true
		}
	}
}
