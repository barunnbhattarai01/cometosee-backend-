package routes

import (
	"cometosee/controller"

	"github.com/gorilla/mux"
)

func PostRoutes(gor *mux.Router, postcontroller *controller.PostController) {
	gor.HandleFunc("/post", postcontroller.UploadPost).Methods("POST")
	gor.HandleFunc("/post/like", postcontroller.LikePost).Methods("POST")
	gor.HandleFunc("/post/comment", postcontroller.CommentPost).Methods("POST")
	gor.HandleFunc("/post/share", postcontroller.SharePost).Methods("POST")
	gor.HandleFunc("/getpost", postcontroller.FetchPost).Methods("GET")
}
