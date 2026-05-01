package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func ConnectionRoutes(gor *mux.Router, connectionController *controller.ConnectionController) {
	gor.HandleFunc("/connection/send", connectionController.SendRequest).Methods("POST")
	gor.HandleFunc("/connection/accept", connectionController.AcceptRequest).Methods("POST")
	gor.HandleFunc("/connection/block", connectionController.BlockUser).Methods("POST")
	gor.HandleFunc("/connection/get", connectionController.GetConnection).Methods("GET")

}
