package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func WebscoketRoutes(gor *mux.Router, wscontroller *controller.WSController) {
	gor.HandleFunc("/ws/manager", middleware.JwtMiddlware(wscontroller.ServeWS)).Methods("GET")

}
