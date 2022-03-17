package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/alisa-vernigor/grpc-go-reverse/pkg/proto/chat"
	"google.golang.org/grpc"
)

type User struct {
	chat chat.ChatRoom_ChatServer
	move chat.ChatRoom_ChatServer
}

type server struct {
	chat.UnimplementedChatRoomServer
	// key -- nickname, value -- instance for streaming
	clients_to_streams map[string]chat.ChatRoom_ChatServer
	streams_to_clients map[chat.ChatRoom_ChatServer]string
	mu                 sync.RWMutex
}

func (s *server) addClient(nick string, stream chat.ChatRoom_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Print("here44")
	for _, ss := range s.clients_to_streams {
		if err := ss.Send(&chat.ChatResponse{Name: nick}); err != nil {
			log.Printf("broadcast err: %v", err)
		}
	}
	fmt.Print("here55")
	s.clients_to_streams[nick] = stream
	fmt.Print("here66")
	s.streams_to_clients[stream] = nick
	fmt.Print("here77")
}

func (s *server) removeClient(stream chat.ChatRoom_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients_to_streams, s.streams_to_clients[stream])
	delete(s.streams_to_clients, stream)
}

func (s *server) getClients() []chat.ChatRoom_ChatServer {
	var cs []chat.ChatRoom_ChatServer

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients_to_streams {
		cs = append(cs, c)
	}
	return cs
}

func (s *server) handle_request(stream chat.ChatRoom_ChatServer, req *chat.ChatRequest) {
	s.addClient(req.GetName(), stream)
}

func (s *server) Chat(stream chat.ChatRoom_ChatServer) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()

	waitc := make(chan struct{})
	go func() {
		for {
			fmt.Print("here1")
			in, err := stream.Recv()
			fmt.Print("here2")
			if err == io.EOF {
				// read done.
				close(waitc)
				s.removeClient(stream)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message : %v", err)
			}
			fmt.Print("here3")
			s.handle_request(stream, in)
			fmt.Print("here4")
		}
	}()
	return nil
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chat.RegisterChatRoomServer(s, &server{
		clients_to_streams: make(map[string]chat.ChatRoom_ChatServer),
		streams_to_clients: make(map[chat.ChatRoom_ChatServer]string),
		mu:                 sync.RWMutex{},
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
