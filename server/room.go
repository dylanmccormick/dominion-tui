package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID               string
	Players          map[uuid.UUID]User // player id : TCP Connetion
	UpdateFunc       func(r *Room)
	BroadcastChannel chan Message
	CommandChannel   chan CommandBody
	Game             any // a room will have a game... with a state. Which we can point to
}

func (r *Room) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("ID: %s\n", r.ID))

	return out.String()
}

func (r *Room) Update() {
	select {
	case msg := <-r.BroadcastChannel:
		r.handleBroadcast(msg)
	case cmd := <-r.CommandChannel:
		r.handleCommand(cmd)
	default:
		return
	}
}

func (r *Room) handleBroadcast(msg Message) {
	for _, u := range r.Players {
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("got error sending message")
			return
		}
		fmt.Fprintf(u.Conn, "%s\r\n", string(data))

	}
}

func (r *Room) handleAction(msg Message) {
	fmt.Println("Handling Action...")
}

func (r *Room) handleCommand(cmd CommandBody) {
	fmt.Printf("Command: %s\n", cmd.Command)
}
