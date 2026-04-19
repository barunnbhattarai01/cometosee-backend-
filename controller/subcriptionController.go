package controller

import (
	"cometosee/common"
	"cometosee/service"
	"encoding/json"
	"net/http"
)

type SubscriptionController struct {
	service service.SubscriptionService
}

func NewSubscriptionController(service service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{service: service}
}

func (c *SubscriptionController) SubscribeUser(email string) error {
	serviceErr := c.service.SubscribeUser(email)

	if serviceErr != nil {
		return serviceErr
	}
	return nil
}

func (c *SubscriptionController) UnsubscribeUser(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email string `json:"email"`
	}

	var reqBody requestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	serviceErr := c.service.UnsubscribeUser(reqBody.Email)

	if serviceErr != nil {
		common.WriteJSONError(w, "failed to unsubscribe user", http.StatusInternalServerError)
		return
	}

	common.WriteJSONMessage(w, "sucessfully deleted the subcription")
}

func (c *SubscriptionController) GetSubscriptionStatus(email string) bool {
	return c.service.GetSubscriptionStatus(email)
}

func (c *SubscriptionController) UpdateSubscriptionEndDate(email string) error {
	service := c.service.UpdateSubscriptionEndDate(email)
	return service
}
