package model

import "github.com/gorilla/websocket"

type Client struct {
	Connection *websocket.Conn
	Egress     chan Event
	Username   string
	ChatRoom   string
}
