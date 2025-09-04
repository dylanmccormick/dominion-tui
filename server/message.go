package server

import (
	"encoding/json"
	// "github.com/google/uuid"
)

type Message struct {
	Requester string         `json:"requester"`
	Typ       string         `json:"message_type"`
	Body      map[string]any `json:"body"`
	Time      int64
}

// {
//     "version": 1,
//     "message_id": "auth_001",
//     "type": "prompt",
//     "ack_needed": true,
//     "body": {
//         "prompt_type": "authentication",
//         "title": "Login Required",
//         "options": ["login", "register"]
//         "fields": [
//             {"name": "username", "type": "text", "required": true},
//             {"name": "password", "type": "password", "required": true}
//         ]
//     }
// }

type Prompt struct {
	Version   string `json:"version"`
	MessageId string `json:"message_id"`
	Type      string `json:"type"`
	AckNeeded bool `json:"ack_needed"`
	Body map[string]any `json:"body"`
}

func decodeMessage(data []byte) Message {
	var msg Message

	if err := json.Unmarshal(data, &msg); err != nil {
		panic(err)
	}

	return msg
}
