package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func SubcribtionRoute(gor *mux.Router, sc *controller.SubscriptionController) {

	gor.HandleFunc("/delete", middleware.JwtMiddlware(sc.UnsubscribeUser)).Methods("POST")
	gor.HandleFunc("/status", middleware.JwtMiddlware(sc.GetSubscriptionStatusHandler)).Methods("GET")
}
