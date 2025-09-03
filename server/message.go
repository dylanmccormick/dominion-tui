package server

import (
	"encoding/json"
	// "github.com/google/uuid"
)

type Message struct {
	Requester string `json:"requester"`
	Type      string `json:"message_type"`
	Body      map[string]any `json:"body"`
	Time      int64 
}

func decodeMessage(data []byte) Message {
	var msg Message

	if err := json.Unmarshal(data, &msg); err != nil {
		panic(err)
	}

	return msg
}
