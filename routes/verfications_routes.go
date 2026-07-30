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

	gor.HandleFunc("/admin/verification/pending", middleware.AdminJwtMiddleware(verificationcontroller.GetPendingVerifications)).Methods("GET")
	gor.HandleFunc("/admin/verification/approve", middleware.AdminJwtMiddleware(verificationcontroller.ApproveVerification)).Methods("POST")
	gor.HandleFunc("/admin/verification/reject", middleware.AdminJwtMiddleware(verificationcontroller.RejectVerification)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/approve", middleware.AdminJwtMiddleware(verificationcontroller.ApprovePlayerDocument)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/reject", middleware.AdminJwtMiddleware(verificationcontroller.RejectPlayerDocument)).Methods("POST")
	gor.HandleFunc("/admin/verification/documents/pending", middleware.AdminJwtMiddleware(verificationcontroller.GetPendingPlayerDocuments)).Methods("GET")
}
