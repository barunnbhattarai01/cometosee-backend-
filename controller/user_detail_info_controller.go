package controller

import (
	"cometosee/common"
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

	authId := common.GetAuthid(r.Context())
	userDetailInfo.AuthId = int(authId)

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

	authId := common.GetAuthid(r.Context())

	userdetailId, err := c.service.GetUserDetailIDByAuthID(int(authId))
	if err != nil {
		common.WriteJSONError(w, "error in fetching user detail id", http.StatusInternalServerError)
		return
	}

	var location model.Location

	err = json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request payload",
		})
		return
	}

	location.User_Detail_Id = userdetailId

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

func (c *UserDetailInfoController) ProfileStatus(w http.ResponseWriter, r *http.Request) {

	authid := common.GetAuthid(r.Context())

	exists, err := c.service.IsProfileCompleted(int(authid))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"profile_completed": exists,
	})
}

func (c *UserDetailInfoController) UpdateUserDetailInfo(w http.ResponseWriter, r *http.Request) {

	var user model.UserDetailInfo

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.WriteJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authId := common.GetAuthid(r.Context())
	user.AuthId = int(authId)

	message, err := c.service.UpdateUserDetailInfo(&user)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
	})
}

func (c *UserDetailInfoController) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	authId := common.GetAuthid(r.Context())

	userdetailId, err := c.service.GetUserDetailIDByAuthID(int(authId))
	if err != nil {
		common.WriteJSONError(w, "error fetching user detail id", http.StatusInternalServerError)
		return
	}

	var location model.Location

	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		common.WriteJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	location.User_Detail_Id = userdetailId

	message, err := c.service.UpdateLocation(&location)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
	})
}
