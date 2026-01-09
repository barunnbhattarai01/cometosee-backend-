package di

import (
	"cometosee/controller"
	"cometosee/repository"
	"cometosee/service"
)

func SetupPostController() *controller.PostController {
	repo := repository.NewPostRepository()
	srv := service.NewPostService(repo)
	ctrl := controller.NewPostController(srv)
	return ctrl
}
