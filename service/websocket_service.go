package service

import (
	"cometosee/model"
	"cometosee/repository"
	"encoding/json"
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

	//webrtrc
	s.Manager.Handlers[model.Web_rtcoffer] = s.HandleWebRTCOffer
	s.Manager.Handlers[model.Web_rtcanswer] = s.HandleWebRTCAnswer
	s.Manager.Handlers[model.Web_rtccandidate] = s.HandleWebRTCCandidate
}

func (s *WebsocketService) RouteEvent(event model.Event, c *model.Client) error {
	if handler, ok := s.Manager.Handlers[event.Type]; ok {
		return handler(event, c)
	}
	return fmt.Errorf("unknown event type: %s", event.Type)
}

// webrtc
// forward offer to target peer
func (s *WebsocketService) HandleWebRTCOffer(event model.Event, c *model.Client) error {
	var signal model.SignalingMessage
	if err := json.Unmarshal(event.Payload, &signal); err != nil {
		return err
	}
	return s.forwardSignal(model.Web_rtcoffer, signal, c)
}

// forward answer to target peer
func (s *WebsocketService) HandleWebRTCAnswer(event model.Event, c *model.Client) error {
	var signal model.SignalingMessage
	if err := json.Unmarshal(event.Payload, &signal); err != nil {
		return err
	}
	return s.forwardSignal(model.Web_rtcanswer, signal, c)
}

// forward ICE candidate to target peer
func (s *WebsocketService) HandleWebRTCCandidate(event model.Event, c *model.Client) error {
	var signal model.SignalingMessage
	if err := json.Unmarshal(event.Payload, &signal); err != nil {
		return err
	}
	return s.forwardSignal(model.Web_rtccandidate, signal, c)
}

// shared forward logic — finds target client by username and sends event
func (s *WebsocketService) forwardSignal(eventType string, signal model.SignalingMessage, from *model.Client) error {
	s.Manager.RLock()
	target, ok := s.Manager.Users[signal.ToID]
	s.Manager.RUnlock()

	if !ok {
		return fmt.Errorf("target peer %s not found", signal.ToID)
	}

	payload, err := json.Marshal(signal)
	if err != nil {
		return err
	}

	target.Egress <- model.Event{
		Type:    eventType,
		Payload: payload,
	}
	return nil
}
