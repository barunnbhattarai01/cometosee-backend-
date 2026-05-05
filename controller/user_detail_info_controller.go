package controller

import (
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type UserDetailInfoController struct {
	service service.UserDetailInfoService
}

func NewUserDetailInfoController(service service.UserDetailInfoService) *UserDetailInfoController {
	return &UserDetailInfoController{service: service}
}

func (c *UserDetailInfoController) TakeUserDetailInfo(w http.ResponseWriter, r *http.Request) {

	var userDetailInfo model.UserDetailInfo
	err := json.NewDecoder(r.Body).Decode(&userDetailInfo)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request payload",
		})
		return

	}

	user, err := c.service.TakeUserDetailInfo(&userDetailInfo)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": user,
	})

}

func (c *UserDetailInfoController) TakeUserLocation(w http.ResponseWriter, r *http.Request) {

	var location model.Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request payload",
		})
		return
	}

	user, err := c.service.TakeUserLocation(&location)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return

	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": user,
	})
}
