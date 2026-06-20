package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
	"os"
)

func SetupPaymentController() *controller.PaymentController {
	paymentRepo := repository.NewPaymentRepository()
	paymentService := service.NewPaymentService(paymentRepo)

	subRepo := repository.NewSubscriptionRepository()
	subService := service.NewSubscriptionService(subRepo)

	esewaSecret := os.Getenv("ESEWA_SECRET")
	successURL := os.Getenv("ESEWA_SUCCESS_URL")
	failureURL := os.Getenv("ESEWA_FAILURE_URL")

	return controller.NewPaymentController(paymentService, subService, esewaSecret, successURL, failureURL)
}
