package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func Discoveryroutes(gor *mux.Router, discoveryController *controller.DiscoveryController) {
	gor.HandleFunc("/feed/dicovery", discoveryController.DiscoverUsers).Methods("GET")

}
