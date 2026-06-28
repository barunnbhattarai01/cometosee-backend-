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

	if err := s.Repo.SaveMessage(msg.From, msg.To, c.ChatRoom, msg.Message, newMsg.Sent); err != nil {
		log.Println("db error:", err)
	}
	// fmt.Print(msg.From)
	// fmt.Print(msg.To, msg.Message, newMsg.Sent, c.ChatRoom)

	data, _ := json.Marshal(newMsg)
	out := model.Event{
		Type:    model.EventNewMessage,
		Payload: data,
	}

	s.Manager.RLock()
	defer s.Manager.RUnlock()

	//send to recipient
	if msg.To != "" {
		if target, ok := s.Manager.Users[msg.To]; ok {
			select {
			case target.Egress <- out:
			default:
				log.Printf("target peer %s is not ready to receive\n", msg.To)
			}
		}
	}

	//echo back to sender too(so tgheir msg appears in their own chat)
	if sender, ok := s.Manager.Users[msg.From]; ok {
		select {
		case sender.Egress <- out:
		default:
			log.Printf("sender peer %s is not ready to receive\n", msg.From)
		}
	}

	return nil
}
