package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func SetupUserfilterroutes(gor *mux.Router, userfiltercontroller *controller.UserFilterController) {
	gor.HandleFunc("/userfilter", userfiltercontroller.FilterUsersByName).Methods("GET")
}
