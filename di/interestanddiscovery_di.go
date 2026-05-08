package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func Setupdiscoveryfeed() *controller.DiscoveryController {
	repo := repository.NewInterestAndLocationDiscoveryRepository()
	svc := service.NewDiscoveryService(repo)
	ctrl := controller.NewDiscoveryController(svc)
	return ctrl
}
