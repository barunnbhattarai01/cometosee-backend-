package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service: service}
}

// signup
func (c *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := c.service.Signup(body); err != nil {
		common.WriteJSONError(w, "error creating user", http.StatusInternalServerError)
		return
	}

	common.WriteJSONMessage(w, "user created sucessfully")

}

// login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, username, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		common.WriteJSONError(w, "login failed", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "login sucessfully",
		"token":    token,
		"username": username,
	})
}
