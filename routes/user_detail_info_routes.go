package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func InitializeUserDetailInfoRoutes(gor *mux.Router, controller *controller.UserDetailInfoController) {
	gor.HandleFunc("/userdetailinfo", middleware.JwtMiddlware(controller.TakeUserDetailInfo)).Methods("POST")
	gor.HandleFunc("/userlocation", middleware.JwtMiddlware(controller.TakeUserLocation)).Methods("POST")
	gor.HandleFunc("/profilestatus", middleware.JwtMiddlware(controller.ProfileStatus)).Methods("GET")
}
