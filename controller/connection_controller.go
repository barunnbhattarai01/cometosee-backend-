package controller

import (
	"cometosee/common"
	"cometosee/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type ConnectionController struct {
	service service.ConnectionService
}

func NewConnectionController(s service.ConnectionService) *ConnectionController {
	return &ConnectionController{service: s}
}

type connectionRequest struct {
	User1 string `json:"user1"`
	User2 string `json:"user2"`
}

func (cc *ConnectionController) SendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	var req connectionRequest
	if err := common.ParseJSONBody(r, &req); err != nil {
		common.WriteJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.User1 == "" || req.User2 == "" {
		common.WriteJSONError(w, "missing users", http.StatusBadRequest)
		return
	}

	err := cc.service.SendRequest(req.User1, req.User2)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusConflict)
		return
	}

	common.WriteJSONMessage(w, "connection send sucessfully")
}

func (cc *ConnectionController) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req connectionRequest
	if err := common.ParseJSONBody(r, &req); err != nil {
		common.WriteJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := cc.service.AcceptRequest(req.User1, req.User2)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "request accepted")
}

func (cc *ConnectionController) BlockUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req connectionRequest
	if err := common.ParseJSONBody(r, &req); err != nil {
		common.WriteJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := cc.service.BlockUser(req.User1, req.User2)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "user blocked")
}

func (cc *ConnectionController) GetConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user1 := r.URL.Query().Get("user1")
	user2 := r.URL.Query().Get("user2")

	if user1 == "" || user2 == "" {
		common.WriteJSONError(w, "missing users", http.StatusBadRequest)
		return
	}

	_, err := cc.service.GetConnection(user1, user2)
	if err != nil {
		common.WriteJSONError(w, "connection not found", http.StatusNotFound)
		return
	}

	common.WriteJSONMessage(w, "connection found")
}

func (c *ConnectionController) GetReceivedRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"message":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	email := common.GetEmail(r.Context())
	fmt.Print(email)

	if email == "" {
		http.Error(w, `{"message":"unauthorized user"}`, http.StatusUnauthorized)
		return
	}

	data, err := c.service.GetReceivedRequests(email)
	if err != nil {
		http.Error(w, `{"message":"failed to fetch data"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "received requests fetched successfully",
		"data":    data,
	})
}

func (c *ConnectionController) GetSentRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"message":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	email := common.GetEmail(r.Context())
	fmt.Print(email)

	if email == "" {
		http.Error(w, `{"message":"unauthorized user"}`, http.StatusUnauthorized)
		return
	}

	data, err := c.service.GetSentRequests(email)
	if err != nil {
		http.Error(w, `{"message":"failed to fetch data"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sent requests fetched successfully",
		"data":    data,
	})
}
