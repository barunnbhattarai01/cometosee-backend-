package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func SubcribtionRoute(gor *mux.Router, subscriptionController *controller.SubscriptionController) {
	gor.HandleFunc("/sub/delete", subscriptionController.UnsubscribeUser).Methods("POST")
}
