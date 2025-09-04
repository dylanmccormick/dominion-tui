package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/dylanmccormick/dominion-tui/internal/utils"
	"github.com/google/uuid"
)

type Server struct {
	Rooms map[string]*Room // id: Room
	Port  string
	Users map[string]User // id: User
}

type Room struct {
	ID             string
	Players        map[uuid.UUID]User // player id : TCP Connetion
	UpdateFunc     func(r *Room)
	InputChannel   chan []byte
	ChatChannel    chan Message
	ActionChannel  chan Message
	CommandChannel chan Message
}

func Init(port string) *Server {
	return &Server{
		Rooms: make(map[string]*Room),
		Port:  port,
	}
}

func (r *Room) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("ID: %s\n", r.ID))
	out.WriteString(fmt.Sprintf("Chat %T", r.InputChannel))

	return out.String()
}

func (s Server) String() string {
	return fmt.Sprintf(
		`Port: %s,
		Users: %v,
		Rooms: %v,
		`,
		s.Port,
		s.Users,
		s.Rooms,
	)
}

func (s Server) Serve() error {
	listener, err := net.Listen("tcp", ":"+s.Port)
	fmt.Printf("Server is listening on port %s\n", s.Port)
	if err != nil {
		return fmt.Errorf("error occurred in server.Serve: %s", err)
	}

	defer listener.Close()
	// Initialize Server?
	s.Rooms["lobby"] = createLobby()

	go s.updateRooms()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("a listener error occurred: %s", err)
		}

		go s.handleRequest(&conn)
	}
}

func (s *Server) updateRooms() {
	for {
		for _, room := range s.Rooms {
			room.GetInputs()
			room.Update()
		}
	}
}

// change this to create a new client and assign to a room
func (s *Server) handleRequest(conn *net.Conn) {
	scanner := bufio.NewReader(*conn) // we shouldn't need to read input yet. Let's break this up logically
	fmt.Fprintf(*conn, "%s\r\n", "Welcome to the server. This is the main menu")
	user := createUser(conn)
	fmt.Fprintf(*conn, "%s\r\n", "Please select a room to join")
	buffer := make([]byte, 4096)
	_, err := scanner.Read(buffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("End of connection closed gracefully")
			return
		} else {
			fmt.Printf("Unknown error from TCP request: %s", err)
		}
		return
	}
	data := bytes.Split(buffer, []byte("\r\n"))
	clean := utils.ClearZeros(data[0])
	name := strings.Trim(string(clean), " \n")
	fmt.Println(s.assignRoom(user, name))
	go user.HandleChat()
}

func (s *Server) assignRoom(user *User, name string) string {
	fmt.Printf("Attempting to join room: %#v", name)
	for k := range s.Rooms {
		fmt.Printf("Room name: %s\n", k)
	}
	room, ok := s.Rooms[name]
	if !ok {
		return "That room does not exist"
	}

	room.Players[user.ID] = *user
	user.Room = room

	return fmt.Sprintf("User added to room %s", room.ID)
}

func createLobby() *Room {
	ic := make(chan []byte)
	chc := make(chan Message, 10)
	cc := make(chan Message, 10)
	ac := make(chan Message, 10)


	return &Room{
		ID:           "lobby",
		Players:      map[uuid.UUID]User{},
		InputChannel: ic,
		ChatChannel: chc,
		CommandChannel: cc,
		ActionChannel: ac,
	}
}

func (r *Room) GetInputs() {
		select {
		case msg := <-r.InputChannel:
			r.handleInput(msg)
		default: 
			return
		}
}

func (r *Room) Update() {
	select {
	case msg := <-r.ChatChannel:
		r.handleChat(msg)
	case msg := <-r.ActionChannel:
		r.handleAction(msg)
	case msg := <-r.CommandChannel:
		r.handleCommand(msg)
	default:
		return
	}
}

func (r *Room) handleInput(msg []byte) {
	decMsg := decodeMessage(msg)
	switch decMsg.Typ {
	case "command":
		r.CommandChannel <- decMsg
	case "action":
		r.ActionChannel <- decMsg
	case "chat":
		r.ChatChannel <- decMsg
	default:
		return
	}
}

func (r *Room) handleChat(msg Message) {
	for _, u := range r.Players {
		m, ok := msg.Body["message"].(string)
		if !ok {
			fmt.Println("Bad message. not a string")
			return
		}
		fmt.Fprintf(*u.Conn, "%s: %s\r\n", u.Username, m)
	}
}

func (r *Room) handleAction(msg Message) {
	fmt.Println("Handling Action...")
}

func (r *Room) handleCommand(msg Message) {
	fmt.Println("Handling Command...")
}
