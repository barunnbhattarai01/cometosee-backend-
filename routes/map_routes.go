package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func Maproutes(gor *mux.Router, mapController *controller.MapController) {
	gor.HandleFunc("/map/pins", middleware.JwtMiddlware(mapController.GetMapEventPins)).Methods("GET")
}
