package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type VerificationController struct {
	service service.VerificationService
}

func NewVerificationController(s service.VerificationService) *VerificationController {
	return &VerificationController{
		service: s,
	}
}

type uploadVerificationRequest struct {
	FrontURL string `json:"front_url"`
	BackURL  string `json:"back_url"`
}

type rejectVerificationRequest struct {
	AuthID int    `json:"auth_id"`
	Reason string `json:"reason"`
}

type approveVerificationRequest struct {
	AuthID int `json:"auth_id"`
}
type approvePlayerDocumentRequest struct {
	DocumentID int `json:"document_id"`
}

type rejectPlayerDocumentRequest struct {
	DocumentID int    `json:"document_id"`
	Reason     string `json:"reason"`
}

func (c *VerificationController) UploadVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authID := int(common.GetAuthid(r.Context()))
	if authID == 0 {
		common.WriteJSONError(w, "missing or invalid auth id", http.StatusUnauthorized)
		return
	}

	var req uploadVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UploadVerification(authID, req.FrontURL, req.BackURL); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "verification submitted")
}

func (c *VerificationController) UploadPlayerDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authID := int(common.GetAuthid(r.Context()))
	if authID == 0 {
		common.WriteJSONError(w, "missing or invalid auth id", http.StatusUnauthorized)
		return
	}

	var doc model.PlayerDocument
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	doc.AuthID = authID

	if err := c.service.UploadPlayerDocument(doc); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "document uploaded")
}

func (c *VerificationController) GetVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authID := int(common.GetAuthid(r.Context()))
	if authID == 0 {
		common.WriteJSONError(w, "missing or invalid auth id", http.StatusUnauthorized)
		return
	}

	v, err := c.service.GetVerification(authID)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if v == nil {
		common.WriteJSONError(w, "verification not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "verification retrieved",
		"verification": v,
	})
}

func (c *VerificationController) GetPlayerDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authID := int(common.GetAuthid(r.Context()))
	if authID == 0 {
		common.WriteJSONError(w, "missing or invalid auth id", http.StatusUnauthorized)
		return
	}

	docs, err := c.service.GetPlayerDocuments(authID)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "documents retrieved",
		"documents": docs,
	})
}

func (c *VerificationController) GetPendingVerifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	list, err := c.service.GetPendingVerifications()
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "pending verifications retrieved",
		"verifications": list,
	})
}

func (c *VerificationController) ApproveVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req approveVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.ApproveVerification(req.AuthID); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "verification approved")
}

func (c *VerificationController) RejectVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req rejectVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.RejectVerification(req.AuthID, req.Reason); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "verification rejected")
}

func (c *VerificationController) ApprovePlayerDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req approvePlayerDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.ApprovePlayerDocument(req.DocumentID); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "document approved")
}

func (c *VerificationController) RejectPlayerDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req rejectPlayerDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.RejectPlayerDocument(req.DocumentID, req.Reason); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "document rejected")
}

// get pending player documents(admin)
func (c *VerificationController) GetPendingPlayerDocuments(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	docs, err := c.service.GetPendingPlayerDocuments()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to fetch pending player documents",
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Pending player documents fetched successfully",
		"data":    docs,
	})
}
