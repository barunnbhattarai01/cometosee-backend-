package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func InitializeUserDetailInfo() *controller.UserDetailInfoController {

	repo := repository.NewUserDetailInfoRepository()
	service := service.NewUserDetailInfoService(repo)
	controller := controller.NewUserDetailInfoController(service)
	return controller
}
