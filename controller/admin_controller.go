package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type AdminController struct {
	service service.AdminService
}

func NewAdminController(adminService service.AdminService) *AdminController {
	return &AdminController{service: adminService}
}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		common.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, admin, err := c.service.Login(request.Email, request.Password)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "admin login successful",
		"token":   token,
		"admin":   admin,
	})
}
