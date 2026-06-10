package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func ConnectionRoutes(gor *mux.Router, connectionController *controller.ConnectionController) {
	gor.HandleFunc("/connection/send", connectionController.SendRequest).Methods("POST")
	gor.HandleFunc("/connection/accept", connectionController.AcceptRequest).Methods("POST")
	gor.HandleFunc("/connection/block", connectionController.BlockUser).Methods("POST")
	gor.HandleFunc("/connection/get", connectionController.GetConnection).Methods("GET")
	gor.HandleFunc("/connection/reject", connectionController.RejectRequest).Methods("POST")
	gor.HandleFunc("/connection/received", middleware.JwtMiddlware(connectionController.GetReceivedRequests)).Methods("GET")
	gor.HandleFunc("/connection/sended", middleware.JwtMiddlware(connectionController.GetSentRequests)).Methods("GET")
	gor.HandleFunc("/connection/unsend", connectionController.UnsendRequest).Methods("POST")
	gor.HandleFunc("/connection/filter", connectionController.UserFilteraftersentandblock).Methods("POST")
	gor.HandleFunc("/connection/connectedpeople", middleware.JwtMiddlware(connectionController.ConnectedPeople)).Methods("GET")
	gor.HandleFunc("/connection/discoveredpeople", middleware.JwtMiddlware(connectionController.DiscoveredPeople)).Methods("GET")

}
