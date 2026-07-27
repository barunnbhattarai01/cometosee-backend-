package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn
	Egress     chan Event
	AuthID     int
	Username   string
	ChatRoom   string
	ClosedOnce sync.Once
}
