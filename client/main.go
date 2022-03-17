package main

import (
	"context"
	"fmt"
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

	ctx := context.Background()

	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var f bool = false
	go func() {
		for {
			mes := clientName
			if !f {
				if err := stream.SendMsg(&chat.ChatRequest{Name: mes}); err != nil {
					log.Fatal(err)
				}
				f = true
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("recv: %s", resp.Name)
	}
}
