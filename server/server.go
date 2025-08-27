package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type User struct {
	Conn net.Conn
	Username string
	ID string
}
type Server struct {
	Rooms map[string]*Room  // id: Room
	Port string
	Users map[string]User // id: User
}

type Room struct {
	ID string
	Players map[string]User // player id : TCP Connetion
}

func Init(port string) *Server {
	return &Server{
		Rooms: make(map[string]*Room),
		Port: port,
	}
}

func (s Server) Serve() error{
	listener, err := net.Listen("tcp", ":"+s.Port)
	fmt.Printf("Server is listening on port %s\n", s.Port)
	if err != nil {
		return fmt.Errorf("error occurred in server.Serve: %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("a listener error occurred: %s", err)
		}

		go s.handleRequest(conn)


	}
}

func (s *Server) handleRequest(conn net.Conn) {
	scanner := bufio.NewReader(conn)
	for {
		fmt.Fprintf(conn, "%s\n", "Welcome to the server. This is the main menu")
		user := s.createUser(conn)
		fmt.Fprintf(conn, "%s\r\n", "Please select a room to join")
		buffer := make([]byte, 4096)
		i, err := scanner.Read(buffer)
		fmt.Printf("Read %d bytes\n", i)
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
		name := string(data[0])
		s.assignRoom(user, name)

	}
}



func (s *Server) createUser(conn net.Conn) *User {
	message := "Please enter a username"
	fmt.Fprintf(conn, "%s\r\n", message)
	scanner := bufio.NewReader(conn)
	buffer := make([]byte, 4096)
	i, err := scanner.Read(buffer)
	fmt.Printf("Read %d bytes\n", i)
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
	name := string(data[0])

	return &User {
		Conn: conn,
		Username: name,
	}
}
func (s *Server) assignRoom(user *User, name string) string{
	room, ok := s.Rooms[name]
	if !ok {
		return "That room does not exist"
	}

	id := "one"
	room.Players[id] = *user
	user.ID = id

	return fmt.Sprintf("User added to room %s", room.ID)
}



// func (c *ServerConfig) handleRequest(conn net.Conn) {
// 	scanner := bufio.NewReader(conn)
// 	defer conn.Close()
// 	log.Printf("Handling connection")
//
// 	for {
// 		buffer := make([]byte, 4096)
// 		i, err := scanner.Read(buffer)
// 		fmt.Printf("Read %d bytes\n", i)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("End of connection closed gracefully")
// 				return
// 			} else {
// 				fmt.Printf("Unknown error from TCP request: %s", err)
// 			}
// 			return
// 		}
// 		buffer = util.ClearZeros(buffer)
// 		log.Printf("clean read into buffer %#v\n", string(buffer))
// 		response, err := cmd.HandleMessage(c.Database, buffer)
// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 		fmt.Println(response)
// 		fmt.Fprintf(conn, "+%s\r\n", response)
// 	}
// }
//
// func (c *ServerConfig) Shell() {
// 	reader := bufio.NewReader(os.Stdin)
//
// 	for {
// 		fmt.Printf("localhost:%s> ", strconv.Itoa(c.Port))
// 		input, err := reader.ReadString('\n')
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("read input: %s\n", input)
// 		strArr := strings.Split(strings.Trim(input, "\n"), " ")
// 		resp, err := cmd.HandleCommand(c.Database, strArr)
// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 			continue
// 		}
// 		fmt.Println(resp)
// 	}
// }
