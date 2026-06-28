package model

import "time"

type Message struct {
	ID       int64     `json:"id"`
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Room     string    `json:"room"`
	Message  string    `json:"message"`
	SentAt   time.Time `json:"sent_at"`
}
