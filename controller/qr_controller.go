package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type QRController struct {
	service *service.QRService
}

func NewQRController(service *service.QRService) *QRController {
	return &QRController{
		service: service,
	}
}

func (c *QRController) GetJoinedEvents(w http.ResponseWriter, r *http.Request) {

	authID := common.GetAuthid(r.Context())

	events, err := c.service.GetJoinedEvents(int(authID))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (c *QRController) GetOwnerPosts(w http.ResponseWriter, r *http.Request) {

	authID := common.GetAuthid(r.Context())

	posts, err := c.service.GetOwnerPosts(int(authID))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (c *QRController) VerifyQR(w http.ResponseWriter, r *http.Request) {

	authID := common.GetAuthid(r.Context())

	var req model.VerifyParticipantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusInternalServerError)

		return
	}

	participant, err := c.service.VerifyQR(req, int(authID))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participant)
}
