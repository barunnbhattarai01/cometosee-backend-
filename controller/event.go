package controller

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send message"
	EventNewMessage  = "new message"
	EventChatRoom    = "change room"
	EventRegister    = "register"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
}

type NewMessage struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeroomEvent struct {
	Name string `json:"sent"`
}

type RegisterEvent struct {
	Name string `json:"name"`
	Room string `json:"room,omitempty"`
}
