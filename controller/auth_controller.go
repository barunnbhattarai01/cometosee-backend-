package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type AuthController struct {
	service   service.AuthService
	AuthTotal prometheus.Counter
}

func NewAuthController(service service.AuthService) *AuthController {

	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_total",
			Help: "Total number of user registered",
		},
	)
	//registerr
	if err := prometheus.Register(counter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			counter = are.ExistingCollector.(prometheus.Counter)
		} else {
			panic(err)
		}
	}

	return &AuthController{
		service:   service,
		AuthTotal: counter,
	}
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
	c.AuthTotal.Inc()
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

func (c *AuthController) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	model.Email = body.Email
	err := c.service.ForgetPassword(body.Email)

	if err != nil {
		common.WriteJSONError(w, "error in changing password", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "sucessfully  verify the email",
	})
}

func (c *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "error request body", http.StatusBadRequest)
		return
	}

	err := c.service.ResetPassword(body.Email, body.Otp.Code, body.Password)
	if err != nil {
		common.WriteJSONError(w, "error in reseting password", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "sucessfully change the password",
	})
}
