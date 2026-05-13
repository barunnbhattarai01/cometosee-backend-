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
	gor.HandleFunc("/updateuserdetailinfo", middleware.JwtMiddlware(controller.UpdateUserDetailInfo)).Methods("PATCH")
	gor.HandleFunc("/updatelocation", middleware.JwtMiddlware(controller.UpdateLocation)).Methods("PATCH")
	gor.HandleFunc("/fetchuserprofile", middleware.JwtMiddlware(controller.GetProfile)).Methods("GET")
}
