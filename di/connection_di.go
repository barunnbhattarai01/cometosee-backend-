package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupConnectionController() *controller.ConnectionController {
	repo := repository.NewConnectionRepo()
	svc := service.NewConnectionService(repo)
	ctrl := controller.NewConnectionController(svc)
	return ctrl
}
