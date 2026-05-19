package service

import (
	"cometosee/features"
	"cometosee/model"
	"cometosee/repository"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type VideoCallService struct {
	repo  repository.VideoCallRepository
	agora features.AgoraService
}

func NewVideoCallService(r repository.VideoCallRepository, a features.AgoraService) *VideoCallService {
	return &VideoCallService{repo: r, agora: a}
}

func (s *VideoCallService) Create(connectionID, userID int64) (*model.VideoCallSession, string, error) {

	channel := fmt.Sprintf("connection-%d-%d", connectionID, time.Now().UnixNano())

	session := &model.VideoCallSession{
		ConnectionID:      connectionID,
		InitiatedByUserID: userID,
		AgoraChannelName:  channel,
		Status:            model.VideoCallSessionInitiated,
	}

	if err := s.repo.Insert(session); err != nil {
		return nil, "", err
	}

	token, err := s.agora.GenerateToken(channel, uint32(userID), 36000)
	if err != nil {
		return nil, "", err
	}

	return session, token, nil
}

func (s *VideoCallService) Start(sessionID, userID, connectionID int64) (*model.VideoCallSession, string, error) {

	session, err := s.repo.GetByID(sessionID)
	if err != nil {
		return nil, "", err
	}

	//fmt.Println("InitiatedBy:", session.InitiatedByUserID)
	//fmt.Println("UserID:", userID)
	//fmt.Println("ConnectionID:", connectionID)
	//fmt.Println("SessionConnection:", session.ConnectionID)

	if session.ConnectionID != connectionID {
		return nil, "", fmt.Errorf("forbidden")
	}

	if session.InitiatedByUserID == userID {
		return nil, "", fmt.Errorf("only receiver can start")
	}

	started, err := s.repo.Start(sessionID)
	if err != nil {
		return nil, "", err
	}

	token, err := s.agora.GenerateToken(started.AgoraChannelName, uint32(userID), 36000)
	if err != nil {
		return nil, "", err
	}

	return started, token, nil
}

func (s *VideoCallService) End(sessionID, connectionID int64) (*model.VideoCallSession, error) {

	session, err := s.repo.GetByID(sessionID)
	if err != nil {
		return nil, err
	}

	if session.ConnectionID != connectionID {
		return nil, fmt.Errorf("forbidden")
	}

	return s.repo.End(sessionID)
}

// handler for video call
func (s *WebsocketService) VideoCallInvite(event model.Event, c *model.Client) error {
	var invite model.VideoCallInviteEvent
	if err := json.Unmarshal(event.Payload, &invite); err != nil {
		return err
	}

	channelName := fmt.Sprintf("call_%d_%d", invite.SessionID, invite.ConnectionID)

	payload := map[string]interface{}{
		"from":          invite.From,
		"to":            invite.To,
		"session_id":    invite.SessionID,
		"connection_id": invite.ConnectionID,
		"channel_name":  channelName,
	}

	updatedEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	event.Payload = updatedEvent

	// forward the raw event as-is to the target user
	s.Manager.RLock()
	target, ok := s.Manager.Users[invite.To]
	s.Manager.RUnlock()

	if ok {
		select {
		case target.Egress <- event:
		default:
			log.Printf("target %s not ready to receive video call invite\n", invite.To)
		}
	} else {
		log.Printf("video call invite: user %s not found\n", invite.To)
	}

	return nil
}
