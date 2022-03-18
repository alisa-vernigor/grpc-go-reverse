package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/alisa-vernigor/grpc-go-reverse/pkg/proto/chat"
	"google.golang.org/grpc"
)

type client struct {
	clients_online map[string]struct{}
}

func (c *client) server_handler(to_sender chan *chat.ClientMessage, stream chat.ChatRoom_ChatClient, wg *sync.WaitGroup) {
	defer wg.Done()
	var in chat.ServerMessage
	for {
		err := stream.RecvMsg(&in)
		if err == io.EOF {
			print("EOF")
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a message : %v", err)
		}

		if in.Type == chat.ServerMessageType_CONNECTED {
			var out []string
			c.clients_online[in.Nickname] = struct{}{}
			for nickname := range c.clients_online {
				out = append(out, nickname)
			}
			fmt.Println("Online players: ", out)
		} else if in.Type == chat.ServerMessageType_MESSAGE_FROM_CHAT {
			fmt.Println(in.Nickname, ": ", in.Body)
		}
	}

}

func (c *client) sender(ch chan *chat.ClientMessage, stream chat.ChatRoom_ChatClient, wg *sync.WaitGroup) {
	defer wg.Done()
	var m *chat.ClientMessage
	for {
		m = <-ch
		stream.SendMsg(m)
	}
}

func (c *client) handle_command(command string) (chat.ClientMessage, error) {
	return chat.ClientMessage{}, nil
}

func (с *client) handle_chat_message(message string) chat.ClientMessage {
	return chat.ClientMessage{Type: chat.ClientMessageType_MESSAGE_TO_CHAT, Body: message}
}

func (c *client) input_handler(to_sender chan *chat.ClientMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)
	var m chat.ClientMessage
	for {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input[0] == '!' {
			var err error
			m, err = c.handle_command(input[1:])
			if err != nil {
				fmt.Println("No such command. Type \"help\" to see all possible commands.")
			}
		} else {
			m = c.handle_chat_message(input)
		}
		to_sender <- &m
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	client := client{clients_online: make(map[string]struct{})}

	to_sender := make(chan *chat.ClientMessage)

	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewChatRoomClient(conn)

	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// sending messages
	go client.sender(to_sender, stream, &wg)

	reader := bufio.NewReader(os.Stdin)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.Replace(clientName, "\n", "", -1)
	to_sender <- &chat.ClientMessage{Type: chat.ClientMessageType_CONNECT, Nickname: clientName}

	go client.input_handler(to_sender, &wg)
	go client.server_handler(to_sender, stream, &wg)

	wg.Wait()
}
