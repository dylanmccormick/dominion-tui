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

type User struct {
	Conn     *net.Conn
	Username string
	ID       uuid.UUID
	RoomId   string
	Room     *Room
}

func (u *User) String() string{
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("Username: %s\n", u.Username))
	out.WriteString(fmt.Sprintf("RoomId: %s\n", u.RoomId))
	out.WriteString(u.Room.String())
	out.WriteString("\n")

	return out.String()

}

func (u *User) GetConnection() *net.Conn {
	return u.Conn
}

func createUser(conn *net.Conn) *User {
	message := "Please enter a username"
	fmt.Fprintf(*conn, "%s\n", message)
	u := &User{
		Conn: conn,
		ID:   uuid.New(),
	}
	name := string(u.GetUserInput())
	u.Username = strings.Trim(name, "\n")
	return u
}

func (u *User) GetUserInput() []byte {
	scanner := bufio.NewReader(*u.Conn)
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
	return clean
}

func (u *User) InputChannel(c chan []byte) {
	for {
		scanner := bufio.NewReader(*u.Conn)
		buffer := make([]byte, 4096)
		_, err := scanner.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of connection closed gracefully")
				break
			} else {
				fmt.Printf("Unknown error from TCP request: %s", err)
				break
			}
		}
		data := bytes.Split(buffer, []byte("\r\n"))
		clean := utils.ClearZeros(data[0])
		prepend := fmt.Appendf(nil, "%s: ", u.Username)
		message := append(prepend, clean...)
		c <- message
	}
}

func (u *User) SendMessage(message string) {
	fmt.Fprintf(*u.Conn, "%s" , message)
}

func (u *User) HandleChat() {
	go u.InputChannel(u.Room.Chat)
	for {
	}
}
