package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Manager struct {
	Clients  ClientList
	Users    map[string]*Client
	Rooms    map[string]ClientList
	Handlers map[string]EventHandler
	sync.RWMutex
}

type EventHandler func(Event, *Client) error

func NewManager() *Manager {
	return &Manager{
		Clients:  make(ClientList),
		Users:    make(map[string]*Client),
		Rooms:    make(map[string]ClientList),
		Handlers: make(map[string]EventHandler),
	}
}

func (m *Manager) AddClient(c *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[c] = true
}

func (m *Manager) RemoveClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Clients[c]; !ok {
		return
	}

	// close websocket
	_ = c.Connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	c.Connection.Close()

	// remove from global clients
	delete(m.Clients, c)

	// remove form users
	if c.Username != "" {
		if cur, ok := m.Users[c.Username]; ok && cur == c {
			delete(m.Users, c.Username)
		}
	}

	// remove from rooms
	if c.ChatRoom != "" {
		if room, ok := m.Rooms[c.ChatRoom]; ok {
			delete(room, c)
			if len(room) == 0 {
				delete(m.Rooms, c.ChatRoom)
			}
		}
	}

	// stop writer goroutine
	close(c.Egress)
}
