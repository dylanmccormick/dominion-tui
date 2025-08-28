package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"

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

func (u *User) GetConnection() *net.Conn {
	return u.Conn
}

func createUser(conn *net.Conn) *User {
	message := "Please enter a username"
	fmt.Fprintf(*conn, "%s\r\n", message)
	u :=  &User{
		Conn:     conn,
		ID: uuid.New(),
	}
	name := string(u.GetUserInput())
	u.Username = name
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
