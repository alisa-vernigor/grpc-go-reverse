package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	s.clients_to_streams[nick] = stream
	s.streams_to_clients[stream] = nick
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

func (s *server) handle_request(stream chat.ChatRoom_ChatServer, req *chat.ClientMessage) {
	if req.Type == chat.ClientMessageType_CONNECT {
		for _, ss := range s.clients_to_streams {
			if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: req.Nickname}); err != nil {
				log.Printf("broadcast err: %v", err)
			}

			if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: s.streams_to_clients[ss]}); err != nil {
				log.Printf("broadcast err: %v", err)
			}
		}

		s.addClient(req.Nickname, stream)
		stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: req.Nickname})
	} else if req.Type == chat.ClientMessageType_MESSAGE_TO_CHAT {
		for _, ss := range s.clients_to_streams {
			if ss != stream {
				if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_MESSAGE_FROM_CHAT, Nickname: s.streams_to_clients[stream], Body: req.Body}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}
		}
	}
}

func (s *server) Chat(stream chat.ChatRoom_ChatServer) error {
	for {
		var in chat.ClientMessage
		err := stream.RecvMsg(&in)

		if err == io.EOF {
			s.removeClient(stream)
			return nil
		}

		if err != nil {
			s.removeClient(stream)
			return nil
		}

		fmt.Println("Got message: ", in.Type.String())

		s.handle_request(stream, &in)
	}
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
