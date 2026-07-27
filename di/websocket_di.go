package di

import (
	"cometosee/controller"
	"cometosee/model"
	"cometosee/repository"
	"cometosee/service"
)

func SetupMessage() *controller.WSController {
	//create manager
	manager := model.NewManager()
	//create repo
	messagerepo := repository.NewMessageRepo()

	postrepo := repository.NewPostRepository()
	//create servicce

	wssrv := service.NewWebscoketService(manager, messagerepo, postrepo)

	//create contrroler
	wsctrl := controller.NewWSController(wssrv)

	return wsctrl
}
