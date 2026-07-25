package controller

import (
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Requirementcontroller struct {
	sev service.RequirementService
}

func NewRequirementController(sev service.RequirementService) *Requirementcontroller {
	return &Requirementcontroller{sev: sev}
}

func (c *Requirementcontroller) CreateRequirement(w http.ResponseWriter, r *http.Request) {

	var req model.Requirement

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.sev.CreateRequirement(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Requirement created successfully",
	})
}

func (c *Requirementcontroller) GetRequirement(w http.ResponseWriter, r *http.Request) {

	postID, _ := strconv.Atoi(mux.Vars(r)["postID"])

	req, err := c.sev.GetRequirement(postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(req)
}

func (c *Requirementcontroller) UpdateRequirement(w http.ResponseWriter, r *http.Request) {

	var req model.Requirement

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.sev.UpdateRequirement(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Requirement updated",
	})
}

func (c *Requirementcontroller) DeleteRequirement(w http.ResponseWriter, r *http.Request) {

	postID, _ := strconv.Atoi(mux.Vars(r)["postID"])

	if err := c.sev.DeleteRequirement(postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Requirement deleted",
	})
}
