package server

import (
	"encoding/json"
	// "github.com/google/uuid"
)

type Field struct {
	Name     string
	Type     string
	Required bool
}
type PromptBody struct {
	PromptType string   `json:"prompt_type"`
	Title      string   `json:"title"`
	Options    []string `json:"options"`
	Fields     []Field  `json:"fields"`
}
type (
	PromptResponseBody struct{}
	GameDiff           struct {
		Diffs []Diff `json:"diff"`
	}
	GameState struct{}
	ChatBody  struct {
		Message string `json:"message"`
	}
	Diff struct {
		Operation string `json:"op"`
		JsonPath  string `json:"path"`
		Value     any    `json:"value"`
	}
)

// Message type to prompt the user to send back input from a menu in the TUI
type Message struct {
	Version   string          `json:"version"`
	MessageId string          `json:"message_id"`
	Type      string          `json:"type"`
	AckNeeded bool            `json:"ack_needed"`
	Time      int64           `json:"time_sent"`
	Body      json.RawMessage `json:"body"`
}

func decodeMessage(data []byte) Message {
	var msg Message

	if err := json.Unmarshal(data, &msg); err != nil {
		panic(err)
	}

	return msg
}
