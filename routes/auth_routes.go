package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func AuthRoutes(gor *mux.Router, authcontroller *controller.AuthController) {
	gor.HandleFunc("/signup", authcontroller.Signup).Methods("POST")
	gor.HandleFunc("/login", authcontroller.Login).Methods("POST")
	gor.HandleFunc("/forgetpassword", authcontroller.ForgetPassword).Methods("POST")
	gor.HandleFunc("/resetpassword", authcontroller.ResetPassword).Methods("POST")
	gor.HandleFunc("/getprofile", middleware.JwtMiddlware(authcontroller.Getprofile)).Methods("GET")
	gor.HandleFunc("/verifyemail", authcontroller.VerifyEmail).Methods("POST")
}
