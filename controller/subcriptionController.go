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

func (c *SubscriptionController) SubscribeUser(authID int, plan string) error {
	return c.service.SubscribeUser(authID, plan)
}

func (c *SubscriptionController) UnsubscribeUser(w http.ResponseWriter, r *http.Request) {
	authID := common.GetAuthid(r.Context())

	if err := c.service.UnsubscribeUser(int(authID)); err != nil {
		common.WriteJSONError(w, "failed to unsubscribe user", http.StatusInternalServerError)
		return
	}

	common.WriteJSONMessage(w, "successfully cancelled the subscription")
}

func (c *SubscriptionController) GetSubscriptionStatusHandler(w http.ResponseWriter, r *http.Request) {
	authID := common.GetAuthid(r.Context())

	active, err := c.service.GetSubscriptionStatus(int(authID))
	if err != nil {
		common.WriteJSONError(w, "failed to fetch subscription status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"active": active,
	})
}

func (c *SubscriptionController) ExtendSubscription(authID int, plan string) error {
	return c.service.ExtendSubscription(authID, plan)
}
