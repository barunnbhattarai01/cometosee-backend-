package service

import (
	"cometosee/model"
	"cometosee/repository"
	"fmt"
)

type WebsocketService struct {
	Manager *model.Manager
	Repo    repository.MessageRepository
}

func NewWebscoketService(m *model.Manager, r repository.MessageRepository) *WebsocketService {
	s := &WebsocketService{
		Manager: m,
		Repo:    r,
	}
	s.registerHandler()
	return s
}

func (s *WebsocketService) registerHandler() {
	s.Manager.Handlers[model.EventRegister] = s.Register
	s.Manager.Handlers[model.EventChatRoom] = s.ChangeRoom
	s.Manager.Handlers[model.EventSendMessage] = s.SendMessage
	s.Manager.Handlers[model.EventVideoCallInvite] = s.VideoCallInvite

}

func (s *WebsocketService) RouteEvent(event model.Event, c *model.Client) error {
	if handler, ok := s.Manager.Handlers[event.Type]; ok {
		return handler(event, c)
	}
	return fmt.Errorf("unknown event type: %s", event.Type)
}
