package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func WEbrtcroutes(gor *mux.Router, webrtccontroller *controller.WebRTCController) {
	gor.HandleFunc("/join", webrtccontroller.JoinRoom).Methods("POST")
	gor.HandleFunc("/leave", webrtccontroller.LeaveRoom).Methods("POST")
	gor.HandleFunc("/addtrack", webrtccontroller.AddTrack).Methods("POST")
	gor.HandleFunc("/removetrack", webrtccontroller.RemoveTrack).Methods("POST")
	gor.HandleFunc("/roomtracks/{roomId}", webrtccontroller.GetTracksInRoom).Methods("GET")
	gor.HandleFunc("/roompeers/{roomId}", webrtccontroller.GetPeersInRoom).Methods("GET")
}
