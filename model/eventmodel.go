package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventSendMessage     = "send message"
	EventNewMessage      = "new message"
	EventChatRoom        = "change room"
	EventRegister        = "register"
	EventVideoCallInvite = "video call invite"
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

type ChangeRoomEvent struct {
	Name string `json:"name"`
}

type RegisterEvent struct {
	Name string `json:"name"`
	Room string `json:"room,omitempty"`
}
type VideoCallInviteEvent struct {
	From         string `json:"from"`
	To           string `json:"to"`
	SessionID    int64  `json:"session_id"`
	ConnectionID int64  `json:"connection_id"`
}
