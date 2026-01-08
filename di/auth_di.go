package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupAuthcontroller() *controller.AuthController {
	repo := repository.NewAuthRepository()
	svc := service.NewAuthService(repo)
	ctrl := controller.NewAuthController(svc)
	return ctrl
}
