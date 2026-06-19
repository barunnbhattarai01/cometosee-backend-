package routes

import (
	"cometosee/controller"
	"cometosee/middleware"

	"github.com/gorilla/mux"
)

func PostRoutes(gor *mux.Router, postcontroller *controller.PostController) {
	gor.HandleFunc("/post", middleware.JwtMiddlware(postcontroller.UploadPost)).Methods("POST")
	gor.HandleFunc("/post/like", middleware.JwtMiddlware(postcontroller.LikePost)).Methods("POST")
	gor.HandleFunc("/post/comment", middleware.JwtMiddlware(postcontroller.CommentPost)).Methods("POST")
	gor.HandleFunc("/post/share", middleware.JwtMiddlware(postcontroller.SharePost)).Methods("POST")
	gor.HandleFunc("/getpost", middleware.JwtMiddlware(postcontroller.FetchFeed)).Methods("POST")
	gor.HandleFunc("/latestlike", middleware.JwtMiddlware(postcontroller.LatestLike)).Methods("GET")
	gor.HandleFunc("/createslot", middleware.JwtMiddlware(postcontroller.CreateSlot)).Methods("POST")
	gor.HandleFunc("/joinslot", middleware.JwtMiddlware(postcontroller.JoinSlot)).Methods("POST")
	gor.HandleFunc("/slot/participant", middleware.JwtMiddlware(postcontroller.GetparticipantsFromslot)).Methods("GET")
}
