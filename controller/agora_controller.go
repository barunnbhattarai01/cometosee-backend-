package controller

import (
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type AgoraController struct {
	service *service.VideoCallService
}

func NewAgoraController(s *service.VideoCallService) *AgoraController {
	return &AgoraController{service: s}
}

func (c *AgoraController) CreateCall(w http.ResponseWriter, r *http.Request) {
	userID := int64(1)       // from context
	connectionID := int64(1) // from context

	session, token, err := c.service.Create(connectionID, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"channel": session.AgoraChannelName,
	})
}

func (c *AgoraController) StartCall(w http.ResponseWriter, r *http.Request) {
	sessionID := int64(1)
	userID := int64(2)
	connectionID := int64(1)

	session, token, err := c.service.Start(sessionID, userID, connectionID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"data":  session,
	})
}

func (c *AgoraController) EndCall(w http.ResponseWriter, r *http.Request) {
	sessionID := int64(1)
	connectionID := int64(1)

	session, err := c.service.End(sessionID, connectionID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(session)
}
