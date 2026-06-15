package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func Map() *controller.MapController {
	repo := repository.NewMap()
	svc := service.NewMapService(repo)
	ctrl := controller.NewMapController(svc)
	return ctrl
}
