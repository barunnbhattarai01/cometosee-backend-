package controller

import (
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type UserFilterController struct {
	service service.UserFilterService
}

func NewUserFilterController(service service.UserFilterService) *UserFilterController {
	return &UserFilterController{service: service}
}

func (c *UserFilterController) FilterUsersByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "innnnvalid method",
		})
		return
	}

	users, err := c.service.FilterUsersByName()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucess",
		"users":   users,
	})
}
