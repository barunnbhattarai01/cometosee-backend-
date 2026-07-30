package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func AdminRoutes(gor *mux.Router, adminController *controller.AdminController) {
	gor.HandleFunc("/admin/login", adminController.Login).Methods("POST")
}
