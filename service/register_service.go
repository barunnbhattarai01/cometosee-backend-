package service

import (
	"cometosee/model"
	"encoding/json"
	"fmt"
)

func (s *WebsocketService) Register(event model.Event, c *model.Client) error {
	var re model.RegisterEvent
	if err := json.Unmarshal(event.Payload, &re); err != nil {
		return err
	}

	if re.Name == "" {
		return fmt.Errorf("username required")
	}

	s.Manager.Lock()
	defer s.Manager.Unlock()

	if old, ok := s.Manager.Users[re.Name]; ok && old != c {
		delete(s.Manager.Clients, old)
		return fmt.Errorf("username already exists")
	}

	c.Username = re.Name
	s.Manager.Users[re.Name] = c

	if re.Room != "" {
		c.ChatRoom = re.Room
		if _, ok := s.Manager.Rooms[re.Room]; !ok {
			s.Manager.Rooms[re.Room] = make(model.ClientList)
		}
		s.Manager.Rooms[re.Room][c] = true
	}
	return nil
}
