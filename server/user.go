package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/dylanmccormick/dominion-tui/internal/utils"
	"github.com/google/uuid"
)

type (
	AuthState int
	User      struct {
		Conn     net.Conn
		Username string
		ID       uuid.UUID
		RoomId   string
		Room     *Room
		State    AuthState
	}
)

const (
	UNAUTHENTICATED AuthState = iota
	AUTHENTICATED
	IN_GAME
)

func CreateNewUser(conn net.Conn) *User {
	// TODO: Add in a prompt here
	prompt := Prompt{
		Version:   "1",
		MessageId: "auth_001",
		Type:      "prompt",
		AckNeeded: true,
		Body:      map[string]any{},
	}

	msg, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(conn, "%s\r\n", msg)

	u := &User{
		Conn: conn,
		ID:   uuid.New(),
	}
	name := string(u.GetUserInput())
	u.Username = strings.Trim(name, "\n")
	return u
}

func (u *User) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("Username: %s\n", u.Username))
	out.WriteString(fmt.Sprintf("RoomId: %s\n", u.RoomId))
	out.WriteString(u.Room.String())
	out.WriteString("\n")

	return out.String()
}

func (u *User) GetConnection() net.Conn {
	return u.Conn
}

func (u *User) GetUserInput() []byte {
	scanner := bufio.NewReader(u.Conn)
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
		scanner := bufio.NewReader(u.Conn)
		buffer := make([]byte, 10240)
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
		for _, dat := range data[:len(data)-1] {
			clean := utils.ClearZeros(dat)
			c <- clean
		}
	}
}

func (u *User) SendMessage(byts []byte) {
	prepend := fmt.Appendf(nil, "%s: ", u.Username)
	message := append(prepend, byts...)
	message = append(message, []byte("\r\n")...)
	fmt.Fprintf(u.Conn, "%s", message)
}

func (u *User) HandleChat() {
	fmt.Printf("Handling chat for user: %s\n", u.Username)
	if u.Room != nil {
		go u.InputChannel(u.Room.InputChannel)
	}
}
