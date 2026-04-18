package controller

import (
	"cometosee/service"
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

func (c *SubscriptionController) UnsubscribeUser(email string) error {
	serviceErr := c.service.UnsubscribeUser(email)

	if serviceErr != nil {
		return serviceErr
	}
	return nil
}

func (c *SubscriptionController) GetSubscriptionStatus(email string) bool {
	return c.service.GetSubscriptionStatus(email)
}

func (c *SubscriptionController) UpdateSubscriptionEndDate(email string) error {
	service := c.service.UpdateSubscriptionEndDate(email)
	return service
}
