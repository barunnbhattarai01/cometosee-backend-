package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func Maproutes(gor *mux.Router, mapController *controller.MapController) {
	gor.HandleFunc("/map/pins", mapController.GetMapEventPins).Methods("GET")
}
