package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/alisa-vernigor/grpc-go-reverse/pkg/proto/chat"
	"google.golang.org/grpc"
)

type daytime int64

const (
	Day daytime = iota
	Night
)

type roles int64

const (
	Mafia roles = iota
	Commissioner
	Civilian
	Dead
)

func (r *roles) to_string() string {
	if *r == Mafia {
		return "Mafia"
	} else if *r == Commissioner {
		return "Commissioner"
	} else if *r == Civilian {
		return "Civilian"
	} else {
		return "Dead"
	}
}

type server struct {
	chat.UnimplementedChatRoomServer
	// key -- nickname, value -- instance for streaming
	clients_to_streams       map[string]chat.ChatRoom_ChatServer
	streams_to_clients       map[chat.ChatRoom_ChatServer]string
	clients_to_sessions      map[string]int
	sessions_to_clients      map[int]map[string]struct{}
	current_lobby_id         int
	session_to_daytime       map[int]daytime
	session_to_day_number    map[int]int
	clients_to_roles         map[string]roles
	sessions_to_finish_votes map[int]int
	sessions_to_alive        map[int]int
	session_to_votes         map[int]map[string]int
	session_to_voted         map[int]map[string]struct{}
	session_to_comm_vote     map[int]string
	session_to_mafia_vote    map[int]string
	session_to_game_over     map[int]bool
	mu                       sync.RWMutex
}

func (s *server) startGame() {
	s.session_to_daytime[s.current_lobby_id] = Day
	s.session_to_day_number[s.current_lobby_id] = 0
	s.sessions_to_alive[s.current_lobby_id] = 4
	s.sessions_to_finish_votes[s.current_lobby_id] = 0
	s.session_to_votes[s.current_lobby_id] = make(map[string]int)
	s.session_to_voted[s.current_lobby_id] = make(map[string]struct{})
	s.session_to_game_over[s.current_lobby_id] = false

	s.session_to_comm_vote[s.current_lobby_id] = ""
	s.session_to_mafia_vote[s.current_lobby_id] = ""

	role := []roles{Civilian, Civilian, Commissioner, Mafia}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(role), func(i, j int) { role[i], role[j] = role[j], role[i] })

	id := 0
	for c := range s.sessions_to_clients[s.current_lobby_id] {
		s.clients_to_roles[c] = role[id]

		if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_GAME_STARTED, Role: role[id].to_string()}); err != nil {
			log.Printf("broadcast err: %v", err)
		}

		id += 1
	}

	for c := range s.sessions_to_clients[s.current_lobby_id] {
		if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DAY, Day: 1}); err != nil {
			log.Printf("broadcast err: %v", err)
		}
	}

	s.current_lobby_id += 1
}

func (s *server) addClient(nick string, stream chat.ChatRoom_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients_to_streams[nick] = stream
	s.streams_to_clients[stream] = nick
	s.clients_to_sessions[nick] = s.current_lobby_id
	if len(s.sessions_to_clients[s.current_lobby_id]) == 0 {
		s.sessions_to_clients[s.current_lobby_id] = make(map[string]struct{})
	}
	s.sessions_to_clients[s.current_lobby_id][nick] = struct{}{}
	if len(s.sessions_to_clients[s.current_lobby_id]) == 4 {
		s.startGame()
	}
}

func (s *server) isGameOver(session int) (bool, string) {
	if session == s.current_lobby_id {
		return false, "Civilian"
	}

	if s.session_to_game_over[session] {
		return false, "Civilian"
	}

	counter := make(map[roles]int)
	for client := range s.sessions_to_clients[session] {
		counter[s.clients_to_roles[client]] += 1
	}

	if counter[Mafia] == 0 {
		return true, "Civilian"
	}
	if counter[Civilian]+counter[Commissioner] == counter[Mafia] {
		return true, "Mafia"
	}
	return false, "Civilian"
}

func (s *server) removeClient(stream chat.ChatRoom_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nickname := s.streams_to_clients[stream]
	delete(s.clients_to_streams, nickname)
	delete(s.streams_to_clients, stream)
	session := s.clients_to_sessions[nickname]
	delete(s.clients_to_sessions, nickname)
	delete(s.sessions_to_clients[session], nickname)

	is_win, winner := s.isGameOver(session)
	if is_win {
		s.session_to_game_over[session] = true
		for c := range s.sessions_to_clients[session] {
			if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_GAME_FINISHED, Role: winner}); err != nil {
				log.Printf("broadcast err: %v", err)
			}
		}
	}
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
		if _, ok := s.clients_to_sessions[req.Nickname]; ok {
			if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
				log.Printf("broadcast err: %v", err)
			}
			return
		} else {
			stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_OK})

			for c := range s.sessions_to_clients[s.current_lobby_id] {
				if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: req.Nickname}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
				if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: c}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}

			s.addClient(req.Nickname, stream)
			stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_CONNECTED, Nickname: req.Nickname})
		}
	} else {
		nick := s.streams_to_clients[stream]
		session := s.clients_to_sessions[nick]

		if s.session_to_game_over[session] {
			if req.Type == chat.ClientMessageType_MESSAGE_TO_CHAT {
				if s.session_to_game_over[session] {
					for c := range s.sessions_to_clients[session] {
						ss := s.clients_to_streams[c]
						if ss != stream {
							if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_MESSAGE_FROM_CHAT, Nickname: s.streams_to_clients[stream], Body: req.Body}); err != nil {
								log.Printf("broadcast err: %v", err)
							}
						}
					}
				}
			} else if req.Type == chat.ClientMessageType_DISCONNECT {
				nickname := s.streams_to_clients[stream]
				session := s.clients_to_sessions[nickname]
				s.removeClient(stream)
				for c := range s.sessions_to_clients[session] {
					ss := s.clients_to_streams[c]
					if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DISCONNECTED, Nickname: nickname}); err != nil {
						log.Printf("broadcast err: %v", err)
					}
				}
			} else {
				if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}
			return
		}

		if req.Type == chat.ClientMessageType_MESSAGE_TO_CHAT {
			nick := s.streams_to_clients[stream]
			session := s.clients_to_sessions[nick]

			if s.session_to_daytime[session] == Day {
				for c := range s.sessions_to_clients[session] {
					ss := s.clients_to_streams[c]
					if ss != stream {
						if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_MESSAGE_FROM_CHAT, Nickname: s.streams_to_clients[stream], Body: req.Body}); err != nil {
							log.Printf("broadcast err: %v", err)
						}
					}
				}
			} else {
				if s.clients_to_roles[nick] != Mafia {
					if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
						log.Printf("broadcast err: %v", err)
					}
				}
			}
		} else if req.Type == chat.ClientMessageType_DISCONNECT {
			nickname := s.streams_to_clients[stream]
			session := s.clients_to_sessions[nickname]
			s.removeClient(stream)
			for c := range s.sessions_to_clients[session] {
				ss := s.clients_to_streams[c]
				if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DISCONNECTED, Nickname: nickname}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}
		} else if req.Type == chat.ClientMessageType_FINISH_THE_DAY {
			nick := s.streams_to_clients[stream]
			session := s.clients_to_sessions[nick]
			if (s.session_to_daytime[session] != Day) || (s.clients_to_roles[nick] == Dead) {
				if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			} else {
				s.mu.Lock()
				defer s.mu.Unlock()
				s.sessions_to_finish_votes[session] += 1
				if s.sessions_to_finish_votes[session] == s.sessions_to_alive[session] {
					s.session_to_daytime[session] = Night
					s.sessions_to_finish_votes[session] = 0

					if s.session_to_day_number[session] != 0 {

						max_votes := 0
						var max_voted string
						all_votes := 0

						for c, v := range s.session_to_votes[session] {
							if v > max_votes {
								max_votes = v
								max_voted = c
							}
							all_votes += v
						}

						if max_votes == 0 || max_votes <= all_votes/2 {
							for c := range s.sessions_to_clients[session] {
								ss := s.clients_to_streams[c]
								if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DRAW}); err != nil {
									log.Printf("broadcast err: %v", err)
								}
							}
						} else {
							for c := range s.sessions_to_clients[session] {
								ss := s.clients_to_streams[c]
								if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_VOTED_OUT, Nickname: max_voted}); err != nil {
									log.Printf("broadcast err: %v", err)
								}
							}
							s.clients_to_roles[max_voted] = Dead
							s.sessions_to_alive[session] -= 1

							is_end, winner := s.isGameOver(session)
							if is_end {
								s.session_to_game_over[session] = true
								for c := range s.sessions_to_clients[session] {
									ss := s.clients_to_streams[c]
									if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_GAME_FINISHED, Role: winner}); err != nil {
										log.Printf("broadcast err: %v", err)
									}
								}
								return
							}
						}
					}

					s.session_to_voted[session] = make(map[string]struct{})
					s.session_to_votes[session] = make(map[string]int)
					s.sessions_to_finish_votes[session] = 0

					for c := range s.sessions_to_clients[session] {
						ss := s.clients_to_streams[c]
						if err := ss.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_NIGHT, Day: int64(s.session_to_day_number[session]) + 1}); err != nil {
							log.Printf("broadcast err: %v", err)
						}
					}
				}
			}
		} else if req.Type == chat.ClientMessageType_PUBLIC {
			nick := s.streams_to_clients[stream]
			session := s.clients_to_sessions[nick]

			if s.clients_to_roles[nick] != Commissioner {
				if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}

			if s.session_to_daytime[session] != Day {
				if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}

			for c := range s.sessions_to_clients[session] {
				if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_PUBLIC_FROM_COM, Nickname: req.Nickname}); err != nil {
					log.Printf("broadcast err: %v", err)
				}
			}
		} else if req.Type == chat.ClientMessageType_VOTE_FOR {
			nick := s.streams_to_clients[stream]
			session := s.clients_to_sessions[nick]

			if s.session_to_daytime[session] == Day {
				if _, ok := s.session_to_voted[session][nick]; ok {
					if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
						log.Printf("broadcast err: %v", err)
					}
				}

				if (s.session_to_daytime[session] != Day) || (s.clients_to_roles[nick] == Dead) || (s.session_to_day_number[session] == 0) {
					if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
						log.Printf("broadcast err: %v", err)
					}
				} else {
					s.mu.Lock()
					defer s.mu.Unlock()
					s.session_to_votes[session][req.Nickname] += 1
					s.session_to_voted[session][nick] = struct{}{}
				}
			} else {
				if (s.clients_to_roles[nick] != Mafia) && (s.clients_to_roles[nick] != Commissioner) {
					if err := stream.SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DENIED}); err != nil {
						log.Printf("broadcast err: %v", err)
					}
				}

				if s.clients_to_roles[nick] == Mafia {
					s.session_to_mafia_vote[session] = req.Nickname
				}
				if s.clients_to_roles[nick] == Commissioner {
					s.session_to_comm_vote[session] = req.Nickname
				}

				counter := make(map[roles]int)
				for client := range s.sessions_to_clients[session] {
					counter[s.clients_to_roles[client]] += 1
				}

				cnt := 0
				if (s.session_to_comm_vote[session] != "") || (counter[Commissioner] == 0) {
					cnt += 1
				}

				if s.session_to_mafia_vote[session] != "" {
					cnt += 1
				}

				if cnt == 2 {
					s.clients_to_roles[s.session_to_mafia_vote[session]] = Dead
					s.sessions_to_alive[session] -= 1

					for c := range s.sessions_to_clients[session] {
						if s.clients_to_roles[c] == Commissioner {
							role := s.clients_to_roles[s.session_to_comm_vote[session]]
							if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_ROLE, Nickname: s.session_to_comm_vote[session], Role: role.to_string()}); err != nil {
								log.Printf("broadcast err: %v", err)
							}
						}
					}

					s.session_to_day_number[session] += 1
					s.session_to_daytime[session] = Day
					for c := range s.sessions_to_clients[session] {
						if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_DAY, Day: int64(s.session_to_day_number[session]) + 1}); err != nil {
							log.Printf("broadcast err: %v", err)
						}
						if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_KILLED, Nickname: s.session_to_mafia_vote[session]}); err != nil {
							log.Printf("broadcast err: %v", err)
						}
					}

					is_win, winner := s.isGameOver(session)
					if is_win {
						s.session_to_game_over[session] = true
						for c := range s.sessions_to_clients[session] {
							if err := s.clients_to_streams[c].SendMsg(&chat.ServerMessage{Type: chat.ServerMessageType_GAME_FINISHED, Role: winner}); err != nil {
								log.Printf("broadcast err: %v", err)
							}
						}
					}
					s.session_to_mafia_vote[session] = ""
					s.session_to_comm_vote[session] = ""
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

func get_adress() string {
	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		fmt.Println(error)

	}

	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String()
}

func main() {
	var lis net.Listener
	var err error
	var addr string
	//reader := bufio.NewReader(os.Stdin)
	for {
		//	fmt.Println("Input port:")
		addr = "9090" //reader.ReadString('\n')
		//addr = strings.Replace(addr, "\n", "", -1)

		addr = ":" + addr
		lis, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		break
	}
	fmt.Println("Listening on:", get_adress()+addr)
	s := grpc.NewServer()
	chat.RegisterChatRoomServer(s, &server{
		clients_to_streams:       make(map[string]chat.ChatRoom_ChatServer),
		streams_to_clients:       make(map[chat.ChatRoom_ChatServer]string),
		clients_to_sessions:      make(map[string]int),
		sessions_to_clients:      make(map[int]map[string]struct{}),
		current_lobby_id:         0,
		session_to_daytime:       make(map[int]daytime),
		session_to_day_number:    make(map[int]int),
		clients_to_roles:         make(map[string]roles),
		sessions_to_finish_votes: make(map[int]int),
		sessions_to_alive:        make(map[int]int),
		session_to_votes:         make(map[int]map[string]int),
		session_to_voted:         make(map[int]map[string]struct{}),
		session_to_comm_vote:     make(map[int]string),
		session_to_mafia_vote:    make(map[int]string),
		session_to_game_over:     make(map[int]bool),
		mu:                       sync.RWMutex{},
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
