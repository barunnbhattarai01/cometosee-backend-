package service

import (
	"cometosee/model"
	"encoding/json"
	"fmt"
	"log"
)

func (s *WebsocketService) ChangeRoom(event model.Event, c *model.Client) error {
	var change model.ChangeRoomEvent

	if err := json.Unmarshal(event.Payload, &change); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}

	if change.Type == "group" {
		allowed, err := s.PostRepo.CanAccessRoomByRoomID(change.Room, c.AuthID)
		if err != nil {
			return err
		}
		if !allowed {
			return fmt.Errorf("unauthorized: not a participant of this chat")
		}
	}

	oldRoom := c.ChatRoom
	newRoom := change.Room

	s.Manager.Lock()
	defer s.Manager.Unlock()

	if oldRoom != "" {
		if room, ok := s.Manager.Rooms[oldRoom]; ok {
			delete(room, c)
			if len(room) == 0 {
				delete(s.Manager.Rooms, oldRoom)
			}
		}
	}

	c.ChatRoom = newRoom
	if newRoom != "" {
		if _, ok := s.Manager.Rooms[newRoom]; !ok {
			s.Manager.Rooms[newRoom] = make(model.ClientList)
		}
		s.Manager.Rooms[newRoom][c] = true
		log.Printf("ChangeRoom: client=%p user=%s room=%q members=%d\n", c, c.Username, newRoom, len(s.Manager.Rooms[newRoom]))
	}

	return nil
}
