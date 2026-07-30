package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupAdminController() *controller.AdminController {
	repo := repository.NewAdminRepository()
	svc := service.NewAdminService(repo)
	return controller.NewAdminController(svc)
}
