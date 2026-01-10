package service

import (
	"cometosee/model"
	"encoding/json"
	"log"
	"time"
)

func (s *WebsocketService) SendMessage(event model.Event, c *model.Client) error {
	var msg model.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &msg); err != nil {
		return err
	}

	newMsg := model.NewMessage{
		SendMessageEvent: msg,
		Sent:             time.Now(),
	}

	if err := s.Repo.SaveMessage(msg.From, c.ChatRoom, msg.Message, newMsg.Sent); err != nil {
		log.Println("db error:", err)
	}

	data, _ := json.Marshal(newMsg)
	out := model.Event{
		Type:    model.EventNewMessage,
		Payload: data,
	}

	s.Manager.RLock()
	defer s.Manager.RUnlock()

	for client := range s.Manager.Rooms[c.ChatRoom] {
		select {
		case client.Egress <- out:
		default:
		}
	}
	return nil
}
