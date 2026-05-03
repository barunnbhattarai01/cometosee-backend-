package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func SetupAgoraRoutes(gor *mux.Router, agoraController *controller.AgoraController) {
	gor.HandleFunc("/createcall", agoraController.CreateCall).Methods("POST")
	gor.HandleFunc("/startcall", agoraController.StartCall).Methods("POST")
	gor.HandleFunc("/endcall", agoraController.EndCall).Methods("POST")
}
