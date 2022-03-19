package main

import (
	"bufio"
	"context"
	"errors"
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
	clients_alive  map[string]struct{}
	known_players  map[string]string
	stream         chat.ChatRoom_ChatClient
	nickname       string
	role           string
	daytime        string
	is_game_over   bool
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
			c.clients_alive[in.Nickname] = struct{}{}
			fmt.Println("Online players: ", out)
		} else if in.Type == chat.ServerMessageType_MESSAGE_FROM_CHAT {
			fmt.Println(in.Nickname, ": ", in.Body)
		} else if in.Type == chat.ServerMessageType_DISCONNECTED {
			var out []string
			delete(c.clients_online, in.Nickname)
			for nickname := range c.clients_online {
				out = append(out, nickname)
			}
			fmt.Println("Online players: ", out)
		} else if in.Type == chat.ServerMessageType_GAME_STARTED {
			fmt.Println("Game started, your role:", in.Role)
			c.role = in.Role
			c.is_game_over = false
		} else if in.Type == chat.ServerMessageType_DAY {
			fmt.Println("It's a new day, day", int64(in.Day))
			c.daytime = "Day"
		} else if in.Type == chat.ServerMessageType_NIGHT {
			fmt.Println("It's a night time, com and mafia are able to do their work now, night", int64(in.Day))
			c.daytime = "Night"
		} else if in.Type == chat.ServerMessageType_VOTED_OUT {
			if in.Nickname == c.nickname {
				fmt.Println("I am sorry, you've been voted out")
				c.role = "Dead"
			} else {
				fmt.Println(in.Nickname, "has been voted out")
			}
			delete(c.clients_alive, in.Nickname)
		} else if in.Type == chat.ServerMessageType_DRAW {
			fmt.Println("City did not choose anyone to kill today")
		} else if in.Type == chat.ServerMessageType_GAME_FINISHED {
			fmt.Println("Game is over. Winner:", in.Role)
			fmt.Println("You can disconnect or stay in lobby talking")
			c.is_game_over = true
		} else if in.Type == chat.ServerMessageType_DENIED {
			fmt.Println("The action is prohibited")
		} else if in.Type == chat.ServerMessageType_ROLE {
			fmt.Println("Now you know of player", in.Nickname, "he is", in.Role)
			c.known_players[in.Nickname] = in.Role
		} else if in.Type == chat.ServerMessageType_KILLED {
			if in.Nickname == c.nickname {
				fmt.Println("I am sorry, you've been killed")
				c.role = "Dead"
			} else {
				fmt.Println(in.Nickname, "has been killed")
			}
			delete(c.clients_alive, in.Nickname)
		} else if in.Type == chat.ServerMessageType_PUBLIC_FROM_COM {
			fmt.Println("The player", in.Nickname, "is mafia")
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
	commands := strings.Split(command, " ")

	if c.is_game_over {
		if commands[0] == "disconnect" {
			return chat.ClientMessage{Type: chat.ClientMessageType_DISCONNECT}, nil
		} else {
			return chat.ClientMessage{}, errors.New("No commands allowed after game")
		}
	} else {
		if commands[0] == "disconnect" {
			return chat.ClientMessage{Type: chat.ClientMessageType_DISCONNECT}, nil
		} else if commands[0] == "finishday" {
			return chat.ClientMessage{Type: chat.ClientMessageType_FINISH_THE_DAY}, nil
		} else if commands[0] == "votefor" {
			if c.role == "Dead" {
				return chat.ClientMessage{}, errors.New("You can't vote, you are dead")
			}
			if c.daytime != "Day" {
				return chat.ClientMessage{}, errors.New("You can't vote at night")
			}

			if commands[1] == c.nickname {
				return chat.ClientMessage{}, errors.New("You can't vote for yourself")
			}

			if _, ok := c.clients_online[commands[1]]; ok {
				if _, ok := c.clients_alive[commands[1]]; ok {
					return chat.ClientMessage{Type: chat.ClientMessageType_VOTE_FOR, Nickname: commands[1]}, nil
				} else {
					return chat.ClientMessage{}, errors.New("This player is not online")
				}
			} else {
				return chat.ClientMessage{}, errors.New("No such player")
			}
		} else if commands[0] == "check" {
			if c.role != "Commissioner" {
				return chat.ClientMessage{}, errors.New("You are not Commissioner")
			}

			if c.daytime != "Night" {
				return chat.ClientMessage{}, errors.New("You can't check at day")
			}

			if commands[1] == c.nickname {
				return chat.ClientMessage{}, errors.New("You can't check yourself")
			}

			if _, ok := c.known_players[commands[1]]; ok {
				return chat.ClientMessage{}, errors.New("You already know the role of this player")
			}

			if _, ok := c.clients_online[commands[1]]; ok {
				if _, ok := c.clients_alive[commands[1]]; ok {
					return chat.ClientMessage{Type: chat.ClientMessageType_VOTE_FOR, Nickname: commands[1]}, nil
				} else {
					return chat.ClientMessage{}, errors.New("This player is not alive")
				}
			} else {
				return chat.ClientMessage{}, errors.New("No such player")
			}
		} else if commands[0] == "public" {
			if c.role != "Commissioner" {
				return chat.ClientMessage{}, errors.New("You are not Commissioner")
			}

			if c.daytime != "Day" {
				return chat.ClientMessage{}, errors.New("You can't public at night")
			}

			for c, r := range c.known_players {
				if r == "Mafia" {
					return chat.ClientMessage{Type: chat.ClientMessageType_PUBLIC, Nickname: c}, nil
				}
			}
			return chat.ClientMessage{}, errors.New("You don't know the Mafia yet")
		} else if commands[0] == "kill" {
			if c.role != "Mafia" {
				return chat.ClientMessage{}, errors.New("You are not Mafia")
			}
			if c.daytime != "Night" {
				return chat.ClientMessage{}, errors.New("You can't kill at day")
			}

			if commands[1] == c.nickname {
				return chat.ClientMessage{}, errors.New("You can't kill yourself")
			}

			if _, ok := c.clients_online[commands[1]]; ok {
				if _, ok := c.clients_alive[commands[1]]; ok {
					return chat.ClientMessage{Type: chat.ClientMessageType_VOTE_FOR, Nickname: commands[1]}, nil
				} else {
					return chat.ClientMessage{}, errors.New("This player is dead")
				}
			} else {
				return chat.ClientMessage{}, errors.New("No such player")
			}
		}
	}

	return chat.ClientMessage{}, errors.New("No such commands")
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

		if input == "" {
			continue
		}

		if input[0] == '!' {
			var err error
			m, err = c.handle_command(input[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			if c.is_game_over {
				m = c.handle_chat_message(input)
				to_sender <- &m
				continue
			}

			if c.role == "Dead" {
				fmt.Println("You are dead, you can't speak")
				continue
			}

			if (c.daytime == "Night") && (c.role != "Mafia") {
				fmt.Println("You are not Mafia, you can't speak at night")
				continue
			}
			m = c.handle_chat_message(input)
		}
		to_sender <- &m
	}
}

func main() {
	var wg sync.WaitGroup
	reader := bufio.NewReader(os.Stdin)
	wg.Add(3)
	client := client{clients_online: make(map[string]struct{}), clients_alive: make(map[string]struct{}), known_players: make(map[string]string)}

	to_sender := make(chan *chat.ClientMessage)

	var conn *grpc.ClientConn
	var err error

	for {
		fmt.Println("Input addr:")
		addr, _ := reader.ReadString('\n')
		addr = strings.Replace(addr, "\n", "", -1)
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			fmt.Println("Wrong adress, try again")
		}
		break
	}
	defer conn.Close()
	c := chat.NewChatRoomClient(conn)

	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	client.stream = stream

	// sending messages
	go client.sender(to_sender, stream, &wg)

	for {
		fmt.Println("Input nickname:")

		clientName, _ := reader.ReadString('\n')
		clientName = strings.Replace(clientName, "\n", "", -1)

		client.nickname = clientName

		to_sender <- &chat.ClientMessage{Type: chat.ClientMessageType_CONNECT, Nickname: clientName}

		var in chat.ServerMessage
		err := stream.RecvMsg(&in)
		if err != nil {
			log.Fatalf("Failed to receive a message : %v", err)
		}

		if in.Type == chat.ServerMessageType_DENIED {
			fmt.Println("Nickname is in use, try again")
		} else {
			break
		}
	}

	go client.input_handler(to_sender, &wg)
	go client.server_handler(to_sender, stream, &wg)

	wg.Wait()
}
