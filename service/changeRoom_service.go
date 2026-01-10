package service

import (
	"cometosee/model"
	"encoding/json"
	"fmt"
)

func (s *WebsocketService) ChangeRoom(event model.Event, c *model.Client) error {
	var change model.ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &change); err != nil {
		return fmt.Errorf("bad payload in change room: %v", err)
	}

	oldRoom := c.ChatRoom
	newRoom := change.Name

	s.Manager.Lock()
	defer s.Manager.Unlock()

	// remove from old room
	if oldRoom != "" {
		if room, ok := s.Manager.Rooms[oldRoom]; ok {
			delete(room, c)
			if len(room) == 0 {
				delete(s.Manager.Rooms, oldRoom)
			}
		}
	}

	// add to new room
	c.ChatRoom = newRoom
	if newRoom != "" {
		if _, ok := s.Manager.Rooms[newRoom]; !ok {
			s.Manager.Rooms[newRoom] = make(model.ClientList)
		}
		s.Manager.Rooms[newRoom][c] = true
	}

	return nil
}
