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
	EventGetHistory      = "get history"
	EventHistoryResponse = "history response"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Room    string `json:"room,omitempty"`
	Type    string `json:"type"`
}

type NewMessage struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Room string `json:"room"`
	Type string `json:"type"`
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

type GetHistoryEvent struct {
	Room     string `json:"room"`
	Limit    int    `json:"limit"`
	BeforeID int64  `json:"before_id,omitempty"`
	AfterID  int64  `json:"after_id,omitempty"`
}

type HistoryResponseEvent struct {
	Messages []Message `json:"messages"`
}
