package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupWebRTC() *controller.WebRTCController {
	//create repo
	repo := repository.NewWebRTCRepository()
	//create service
	srv := service.NewWebRTCService(repo)
	//create controller
	ctrl := controller.NewWebRTCController(srv)
	return ctrl
}
