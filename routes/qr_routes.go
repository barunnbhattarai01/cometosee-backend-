package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func QRROUTES(router *mux.Router, qrcontroller *controller.QRController) {
	router.HandleFunc("/qr/joined", middleware.JwtMiddlware(qrcontroller.GetJoinedEvents)).Methods("GET")
	router.HandleFunc("/qr/my-posts", middleware.JwtMiddlware(qrcontroller.GetOwnerPosts)).Methods("GET")
	router.HandleFunc("/qr/verify", middleware.JwtMiddlware(qrcontroller.VerifyQR)).Methods("POST")
}
