package routes

import (
	"cometosee/controller"
	"cometosee/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func PaymentRoutes(gor *mux.Router, pc *controller.PaymentController) {
	gor.Handle("/esewa/initiate", middleware.JwtMiddlware(
		http.HandlerFunc(pc.InitiateHandler),
	)).Methods("POST")

	gor.HandleFunc("/esewa/verify", pc.VerifyHandler).Methods("GET")
	gor.HandleFunc("/esewa/failure", pc.FailureHandler).Methods("GET")
}
