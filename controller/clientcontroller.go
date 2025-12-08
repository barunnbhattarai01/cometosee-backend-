package controller

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 60 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	chatroom   string
	username   string
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		//cleanup conection
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.connection.SetReadLimit(512)
	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error in reading messages:%v", err)
			}
			break
		}
		//sklip the invalid json
		if !json.Valid(payload) {
			log.Println("invalid json received:", string(payload))
			continue
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Print("error handling event", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Print("errror in handling message:", err)
		}

	}

}

func (c *Client) WriteMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				//chaneel closeed and cleose the connection
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Print(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("failed to send message:", err)
				return
			}
			log.Println("message sent")

		case <-ticker.C:

			//send ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ping error:", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Print("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
