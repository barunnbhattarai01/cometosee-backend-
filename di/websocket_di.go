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
	//create servicce
	wssrv := service.NewWebscoketService(manager, messagerepo)

	//create contrroler
	wsctrl := controller.NewWSController(wssrv)

	return wsctrl
}
