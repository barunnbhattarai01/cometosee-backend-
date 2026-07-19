package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func VerificationDi() *controller.VerificationController {
	//repo
	verificationrepo := repository.NewVerificationRepository()

	//service
	verificationservice := service.NewVerificationService(verificationrepo)

	//controller
	verificationctrl := controller.NewVerificationController(verificationservice)

	return verificationctrl
}
