package service

import (
	"cometosee/model"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (s *WebsocketService) SendMessage(event model.Event, c *model.Client) error {
	var msg model.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &msg); err != nil {
		return err
	}

	newMsg := model.NewMessage{SendMessageEvent: msg, Sent: time.Now()}

	room := msg.Room
	if room == "" {
		room = c.ChatRoom
	}

	if msg.Type == "group" {
		allowed, err := s.PostRepo.CanAccessRoomByRoomID(msg.Room, c.AuthID)
		if err != nil || !allowed {
			return fmt.Errorf("unauthorized: not a participant of this chat")
		}
	}

	if err := s.Repo.SaveMessage(msg.From, msg.To, room, msg.Message, newMsg.Sent); err != nil {
		log.Println("db error:", err)
	}

	data, _ := json.Marshal(newMsg)
	out := model.Event{Type: model.EventNewMessage, Payload: data}

	s.Manager.RLock()
	defer s.Manager.RUnlock()

	if msg.Type == "group" {
		log.Printf("SendMessage: client=%p user=%s room=%q members=%d\n", c, msg.From, room, len(s.Manager.Rooms[room]))
		// broadcast to everyone currently in this room
		for member := range s.Manager.Rooms[room] {
			select {
			case member.Egress <- out:
			default:
				log.Printf("member %s not ready to receive\n", member.Username)
			}
		}
		return nil
	}

	if msg.Type == "private" {
		if msg.To != "" {
			if target, ok := s.Manager.Users[msg.To]; ok {
				select {
				case target.Egress <- out:
				default:
				}
			}
		}
		if sender, ok := s.Manager.Users[msg.From]; ok {
			select {
			case sender.Egress <- out:
			default:
			}
		}
	}
	return nil
}
