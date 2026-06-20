package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupSubscriptionController() *controller.SubscriptionController {
	subRepo := repository.NewSubscriptionRepository()
	subService := service.NewSubscriptionService(subRepo)

	return controller.NewSubscriptionController(subService)
}
