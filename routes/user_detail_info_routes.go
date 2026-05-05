package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func InitializeUserDetailInfoRoutes(gor *mux.Router, controller *controller.UserDetailInfoController) {
	gor.HandleFunc("/userdetailinfo", controller.TakeUserDetailInfo).Methods("POST")
	gor.HandleFunc("/userlocation", controller.TakeUserLocation).Methods("POST")
}
