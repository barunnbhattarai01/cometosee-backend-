package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cometosee/model"
	"cometosee/service"

	"github.com/gorilla/websocket"
)

type WSController struct {
	Service *service.WebsocketService
}

func NewWSController(s *service.WebsocketService) *WSController {
	return &WSController{Service: s}
}

var (
	pongWait     = 60 * time.Second
	pingInterval = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (c *WSController) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &model.Client{
		Connection: conn,
		Egress:     make(chan model.Event),
	}

	c.Service.Manager.Clients[client] = true

	go c.read(client)
	go c.write(client)
}

func (c *WSController) read(client *model.Client) {
	defer func() {
		c.Service.Manager.RemoveClient(client)
	}()

	conn := client.Connection

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("read error:", err)
			}
			return
		}

		if !json.Valid(payload) {
			log.Println("invalid json:", string(payload))
			continue
		}

		var event model.Event
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		if err := c.Service.RouteEvent(event, client); err != nil {
			log.Println("route error:", err)
		}
	}
}

func (c *WSController) write(client *model.Client) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Service.Manager.RemoveClient(client)
	}()

	for {
		select {
		case msg, ok := <-client.Egress:
			if !ok {
				_ = client.Connection.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				return
			}

			if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("write error:", err)
				return
			}

		case <-ticker.C:
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
