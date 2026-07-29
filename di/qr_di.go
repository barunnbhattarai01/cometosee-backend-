package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func QrDI() *controller.QRController {
	repo := repository.NewQRRepository()
	srv := service.NewQRService(repo)
	ctrl := controller.NewQRController(srv)
	return ctrl
}
