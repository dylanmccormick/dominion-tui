package server

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID             string
	Players        map[uuid.UUID]User // player id : TCP Connetion
	UpdateFunc     func(r *Room)
	InputChannel   chan []byte
	ChatChannel    chan Message
	ActionChannel  chan Message
	CommandChannel chan Message
}

func (r *Room) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("ID: %s\n", r.ID))
	out.WriteString(fmt.Sprintf("Chat %T", r.InputChannel))

	return out.String()
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
		fmt.Fprintf(u.Conn, "%s: %s\r\n", u.Username, m)
	}
}

func (r *Room) handleAction(msg Message) {
	fmt.Println("Handling Action...")
	
}

func (r *Room) handleCommand(msg Message) {
	fmt.Println("Handling Command...")
}
