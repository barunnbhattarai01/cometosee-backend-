package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpdgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	Clients ClientList
	Users   map[string]*Client
	Rooms   map[string]ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManger(ctx context.Context) *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		Users:    make(map[string]*Client),
		Rooms:    make(map[string]ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandler()
	return m
}

func (m *Manager) setupEventHandler() {
	m.handlers[EventSendMessage] = m.SendMessage
	m.handlers[EventChatRoom] = m.ChatRoomhandler
	m.handlers[EventRegister] = m.RegisterHandler
}

func (m *Manager) RegisterHandler(event Event, c *Client) error {
	var re RegisterEvent
	if err := json.Unmarshal(event.Payload, &re); err != nil {
		return fmt.Errorf("bad payload in register: %v", err)
	}
	if re.Name == "" {
		return fmt.Errorf("regsiter :name required")
	}
	m.Lock()
	defer m.Unlock()

	//if username already taken ,retunr errrrrir
	if _, exists := m.Users[re.Name]; exists {
		return fmt.Errorf("username already taken")
	}
	c.username = re.Name
	m.Users[re.Name] = c

	if re.Room != "" {
		c.chatroom = re.Room
		if _, ok := m.Rooms[re.Room]; !ok {
			m.Rooms[re.Room] = make(ClientList)
		}
		m.Rooms[re.Room][c] = true
	}
	log.Printf("registered user %s in room %s", c.username, c.chatroom)
	return nil
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("new connection")
	websocketUpdgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := websocketUpdgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading websocket: %v", err)
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.WriteMessage()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Clients[client]; ok {
		//close connection
		_ = client.connection.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		client.connection.Close()

		//remove from global clients
		delete(m.Clients, client)

		//remove from User map
		if client.username != "" {
			if cur, exists := m.Users[client.username]; exists && cur == client {
				delete(m.Users, client.username)
			}
		}

		//remove from room
		if client.chatroom != "" {
			if cl, ok := m.Rooms[client.chatroom]; ok {
				delete(cl, client)
				if len(cl) == 0 {
					delete(m.Rooms, client.chatroom)
				}

			}

			//close the egress cahnnel to signal writer to stop
			close(client.egress)
		}
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	//check event type
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) ChatRoomhandler(event Event, c *Client) error {
	var change ChangeroomEvent
	if err := json.Unmarshal(event.Payload, &change); err != nil {
		return fmt.Errorf("bad payload in change room: %v", err)
	}
	oldRoom := c.chatroom
	newRoom := change.Name

	m.Lock()
	defer m.Unlock()

	// remove from old room
	if oldRoom != "" {
		if cl, ok := m.Rooms[oldRoom]; ok {
			delete(cl, c)
			// if empty, delete the room map entry
			if len(cl) == 0 {
				delete(m.Rooms, oldRoom)
			}
		}
	}

	// add to new room
	c.chatroom = newRoom
	if newRoom != "" {
		if _, ok := m.Rooms[newRoom]; !ok {
			m.Rooms[newRoom] = make(ClientList)
		}
		m.Rooms[newRoom][c] = true
	}
	return nil
}

func (m *Manager) SendMessage(event Event, c *Client) error {
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in send message: %v", err)
	}

	broadMessage := NewMessage{
		SendMessageEvent: chatevent,
		Sent:             time.Now(),
	}

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal outgoing message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventNewMessage,
	}

	m.RLock()
	defer m.RUnlock()

	if chatevent.To != "" {
		if dest, ok := m.Users[chatevent.To]; ok {
			select {
			case dest.egress <- outgoingEvent:
			default:
				log.Printf("drop message to %s (egress full)", chatevent.To)
			}
			// send to sender
			if c != nil && c != dest {
				select {
				case c.egress <- outgoingEvent:
				default:
				}
			}
			return nil
		}
		return fmt.Errorf("user %s not found", chatevent.To)
	}

	room := c.chatroom
	if room == "" {
		return fmt.Errorf("you are not in any room")
	}
	if members, ok := m.Rooms[room]; ok {
		for client := range members {
			select {
			case client.egress <- outgoingEvent:
			default:
				log.Printf("drop message to user (egress full)")
			}
		}
		return nil
	}

	return fmt.Errorf("room %s not found", room)
}
