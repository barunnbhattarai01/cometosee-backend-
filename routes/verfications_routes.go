package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func VerficationsRoutes(gor *mux.Router, verificationcontroller *controller.VerificationController) {

	gor.HandleFunc("/verification/upload", middleware.JwtMiddlware(verificationcontroller.UploadVerification)).Methods("POST")
	gor.HandleFunc("/verification/documents", middleware.JwtMiddlware(verificationcontroller.UploadPlayerDocument)).Methods("POST")
	gor.HandleFunc("/verification", middleware.JwtMiddlware(verificationcontroller.GetVerification)).Methods("GET")
	gor.HandleFunc("/verification/player/documents", middleware.JwtMiddlware(verificationcontroller.GetPlayerDocuments)).Methods("GET")

	//admin  TODO need admin middleware

	gor.HandleFunc("/admin/verification/pending", middleware.JwtMiddlware(verificationcontroller.GetPendingVerifications)).Methods("GET")
	gor.HandleFunc("/admin/verification/approve", middleware.JwtMiddlware(verificationcontroller.ApproveVerification)).Methods("POST")
	gor.HandleFunc("/admin/verification/reject", middleware.JwtMiddlware(verificationcontroller.RejectVerification)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/approve", middleware.JwtMiddlware(verificationcontroller.ApprovePlayerDocument)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/reject", middleware.JwtMiddlware(verificationcontroller.RejectPlayerDocument)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/pending", middleware.JwtMiddlware(verificationcontroller.GetPendingPlayerDocuments)).Methods("GET")
}
