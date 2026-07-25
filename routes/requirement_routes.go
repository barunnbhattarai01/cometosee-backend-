package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func Requirementroutes(router *mux.Router, requirementController *controller.Requirementcontroller) {
	router.HandleFunc("/requirements", requirementController.CreateRequirement).Methods("POST")
	router.HandleFunc("/requirements/{postID}", requirementController.GetRequirement).Methods("GET")
	router.HandleFunc("/update/requirements", requirementController.UpdateRequirement).Methods("PUT")
	router.HandleFunc("/requirements/{postID}", requirementController.DeleteRequirement).Methods("DELETE")
}
