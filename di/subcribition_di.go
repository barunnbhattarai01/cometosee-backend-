package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SubcribntionDi() *controller.SubscriptionController {
	//repo
	subscriptionRepo := repository.NewSubscriptionRepository()
	//service
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	//controller
	subscriptionController := controller.NewSubscriptionController(subscriptionService)
	return subscriptionController
}
