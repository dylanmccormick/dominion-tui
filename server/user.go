package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

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
		Buffer   []byte
	}
)

const (
	UNAUTHENTICATED AuthState = iota
	AUTHENTICATED
	IN_GAME
)

func CreateNewUser(conn net.Conn) *User {
	pb, err := json.Marshal(PromptBody{
		PromptType: "authentication",
		Title:      "Login Required",
		Options:    []string{"login", "register"},
		Fields: []Field{
			{Name: "username", Type: "text", Required: true},
			{Name: "password", Type: "password", Required: true},
		},
	})

	if err != nil {
		fmt.Println(fmt.Errorf("Json marshaling went bad for prompt: %s", err))
		return nil
	}
	prompt := Message{
		Version:   "1",
		MessageId: "auth_001",
		Type:      "prompt",
		AckNeeded: true,
		Body: pb,
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

func (u *User) InputChannel() {
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
		u.Buffer = append(u.Buffer, utils.ClearZeros(buffer)...)
	}
}

func (u *User) HandleMessages() {
	go u.InputChannel()
	for {
		u.ProcessMessage()
		time.Sleep(1 * time.Second)
	}
	
}

// This method will read the user input and route it where it needs to go. May do some transformation or whatever
func (u *User) ProcessMessage() {
	// Find the first iteration of \r\n in u.Buffer
	// Process That message
	var message Message

	fmt.Println("gettting index of registered nurse")
	idx := bytes.Index(u.Buffer, []byte("\r\n"))
	fmt.Println(idx)
	if idx == -1 {
		return
	}
	data := u.Buffer[:idx]
	u.Buffer = u.Buffer[idx+2:]
	err := json.Unmarshal(data, &message)
	if err != nil {
		panic(err)
	}

	switch message.Type {
		case "prompt_response":
			u.handlePromptResponse(message)
		case "chat":
			u.handleChat(message)
		case "action":
			u.handleAction(message)
		case "command":
			u.handleCommand(message)
	}
}

func (u *User) handlePromptResponse(msg Message) {
}

func (u *User) handleChat(msg Message) {
	// Take the chat. Add the username to the message <username>: 
	// Push the message to the Broadcast channel of the room
	var cb ChatBody

	err := json.Unmarshal(msg.Body, &cb)
	if err != nil {
		e := fmt.Errorf("Invalid chat body: %s", err)
		panic(e)
	}

	cb.Message = fmt.Sprintf("%s: %s", u.Username, cb.Message)
	
	chatData, err := json.Marshal(cb)
	if err != nil {
		e := fmt.Errorf("Invalid chat alteration: %s", err)
		panic(e)
	}
	msg.Body = chatData

	u.Room.BroadcastChannel <- msg
}

func (u *User) handleAction(msg Message) {
}

func (u *User) handleCommand(msg Message) {
}
 
func (u *User) SendMessage(byts []byte) {
	prepend := fmt.Appendf(nil, "%s: ", u.Username)
	message := append(prepend, byts...)
	message = append(message, []byte("\r\n")...)
	fmt.Fprintf(u.Conn, "%s", message)
}

// TODO: Some sort of loop to constantly check for input from user and put it into the BUFFER
// TODO: Loop/ Goroutine to constantly process the BUFFER for the user
