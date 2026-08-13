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
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.AuthTotal.Inc()
	common.WriteJSONMessage(w, "User registered successfully. Please check your email for the OTP to verify your account.")

}

// login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid method",
		})
		return
	}

	var body model.Auth
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, username, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
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

func (c *AuthController) Getprofile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email, ok := r.Context().Value("email").(string)
	if !ok {
		common.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.service.GetProfile(email)
	if err != nil {
		common.WriteJSONError(w, "error fetching profile", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "profile fetched sucessfully",
		"profile": user,
	})

}

func (c *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ok, err := c.service.VerifyEmail(body.Email, body.OTP)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ok {
		common.WriteJSONError(w, "verification failed", http.StatusBadRequest)
		return
	}

	common.WriteJSONMessage(w, "User registered successfully")
}
