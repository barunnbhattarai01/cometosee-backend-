package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func UserfilterDI() *controller.UserFilterController {
	//repo
	userfilterrepo := repository.NewUserFilterRepository()

	//service
	userfilterservice := service.NewUserFilterService(userfilterrepo)

	//controller
	userfilterctrl := controller.NewUserFilterController(userfilterservice)

	return userfilterctrl
}
