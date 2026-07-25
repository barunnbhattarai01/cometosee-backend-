package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func RequirementDi() *controller.Requirementcontroller {
	repo := repository.NewRequirementRepository()
	srv := service.NewRequirementService(repo)
	ctrl := controller.NewRequirementController(srv)
	return ctrl
}
