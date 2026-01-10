package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func WebscoketRoutes(gor *mux.Router, wscontroller *controller.WSController) {
	gor.HandleFunc("/ws/manager", wscontroller.ServeWS).Methods("GET")

}
