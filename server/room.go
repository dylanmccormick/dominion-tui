package server

import (
	"bytes"
	"encoding/json"
	// "encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID               string
	Players          map[uuid.UUID]User // player id : TCP Connetion
	UpdateFunc       func(r *Room)
	BroadcastChannel chan Message
}

func (r *Room) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("ID: %s\n", r.ID))

	return out.String()
}

func (r *Room) Update() {
	select {
	case msg := <- r.BroadcastChannel:
		r.handleChat(msg)
	default:
		return
	}
}

// func (r *Room) Update() {
// 	select {
// 	case msg := <-r.ChatChannel:
// 		r.handleChat(msg)
// 	case msg := <-r.ActionChannel:
// 		r.handleAction(msg)
// 	case msg := <-r.CommandChannel:
// 		r.handleCommand(msg)
// 	default:
// 		return
// 	}
// }
//
// func (r *Room) handleInput(msg []byte) {
// 	decMsg := decodeMessage(msg)
// 	switch decMsg.Type {
// 	case "command":
// 		r.CommandChannel <- decMsg
// 	case "action":
// 		r.ActionChannel <- decMsg
// 	case "chat":
// 		r.ChatChannel <- decMsg
// 	default:
// 		return
// 	}
// }

func (r *Room) handleChat(msg Message) {
	for _, u := range r.Players {
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("got error sending message")
			return
		}
		fmt.Fprintf(u.Conn, "%s\r\n", string(data))

	}
}

// TODO: Rethink. Chat body will be just passed from server to client. Is there a need to translate it?
// func (r *Room) handleChat(msg Message) {
// 	for _, u := range r.Players {
// 		var cb ChatBody
// 		json.Unmarshal(msg, &cb)
// 		if !ok {
// 			fmt.Println("Bad message. not a string")
// 			return
// 		}
// 		m := cb.Message
// 		fmt.Fprintf(u.Conn, "%s: %s\r\n", u.Username, m)
// 	}
// }

func (r *Room) handleAction(msg Message) {
	fmt.Println("Handling Action...")
}

func (r *Room) handleCommand(msg Message) {
	fmt.Println("Handling Command...")
}
