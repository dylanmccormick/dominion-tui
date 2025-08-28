package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/dylanmccormick/dominion-tui/internal/utils"
)

type User struct {
	Conn     *net.Conn
	Username string
	ID       string
	RoomId   string
	Room     *Room
}

type Server struct {
	Rooms map[string]*Room // id: Room
	Port  string
	Users map[string]User // id: User
}

type Room struct {
	ID         string
	Players    map[string]User // player id : TCP Connetion
	UpdateFunc func(r *Room)
}

func Init(port string) *Server {
	return &Server{
		Rooms: make(map[string]*Room),
		Port:  port,
	}
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
			fmt.Println(room.ID)
			room.UpdateFunc(room)
		}
	}
}

func (s *Server) handleRequest(conn *net.Conn) {
	scanner := bufio.NewReader(*conn)
	s.Rooms["lobby"] = createLobby()
	fmt.Fprintf(*conn, "%s\n", "Welcome to the server. This is the main menu")
	user := s.createUser(conn)
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
	fmt.Println(name)
	fmt.Println(s.assignRoom(user, name))
}

func (u *User) GetConnection() *net.Conn {
	return u.Conn
}

func (s *Server) createUser(conn *net.Conn) *User {
	message := "Please enter a username"
	fmt.Fprintf(*conn, "%s\r\n", message)
	scanner := bufio.NewReader(*conn)
	buffer := make([]byte, 4096)
	_, err := scanner.Read(buffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("End of connection closed gracefully")
			return nil
		} else {
			fmt.Printf("Unknown error from TCP request: %s", err)
		}
		return nil
	}
	data := bytes.Split(buffer, []byte("\r\n"))
	clean := utils.ClearZeros(data[0])
	name := string(clean)

	return &User{
		Conn:     conn,
		Username: name,
	}
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

	id := "one"
	room.Players[id] = *user
	user.ID = id

	return fmt.Sprintf("User added to room %s", room.ID)
}

func createLobby() *Room {
	f := func(r *Room) {
		for _, v := range r.Players {
			fmt.Println(v.Username)
			fmt.Fprintf(*v.Conn, "Hello from the lobby, %s", v.Username)
		}
	}
	return &Room{
		ID:         "lobby",
		Players:    map[string]User{},
		UpdateFunc: f,
	}
}
