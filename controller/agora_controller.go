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

type req struct {
	SessionID    int64 `json:"session_id"`
	UserID       int64 `json:"user_id"`
	ConnectionID int64 `json:"connection_id"`
}

func (c *AgoraController) CreateCall(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request req
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session, token, err := c.service.Create(request.ConnectionID, request.UserID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": session.Id,
		"token":      token,
		"channel":    session.AgoraChannelName,
	})
}

func (c *AgoraController) StartCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request req
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session, token, err := c.service.Start(request.SessionID, request.UserID, request.ConnectionID)
	if err != nil {
		//fmt.Print("errorror hererere")
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"data":  session,
	})
}

func (c *AgoraController) EndCall(w http.ResponseWriter, r *http.Request) {
	var request req
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session, err := c.service.End(request.SessionID, request.ConnectionID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(session)
}
