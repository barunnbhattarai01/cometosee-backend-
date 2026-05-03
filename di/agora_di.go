package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
	"os"
)

func SetupAgoraDI() (*controller.AgoraController, error) {
	//agora token
	agoraid := service.NewAgoraService(
		os.Getenv("Agora_APP_ID"),
		os.Getenv("Agora_APP_CERTIFICATE"),
	)

	//repo
	repo := repository.NewVideoCallRepo()
	//service
	service := service.NewVideoCallService(repo, agoraid)
	//controller
	controller := controller.NewAgoraController(service)
	return controller, nil
}
