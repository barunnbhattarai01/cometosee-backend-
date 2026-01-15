package service

import (
	"cometosee/model"
	"cometosee/repository"
	"context"
	"sync"
	"time"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) UploadPost(caption, imageUrl, community, username string) (string, error) {
	post := model.POSTDATA{
		Caption:   caption,
		ImageUrl:  imageUrl,
		Community: community,
		Username:  username,
		Created:   time.Now(),
	}
	return s.repo.CreatePOST(post)
}

func (s *PostService) LikePost(postId, userId string) (bool, error) {
	return s.repo.ToggleLike(postId, userId)
}

func (s *PostService) AddComment(postId, userId, comment string) error {
	return s.repo.Addcomment(postId, userId, comment)
}

func (s *PostService) FetchPost(ctx context.Context) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	posts, err := s.repo.FetchPost()

	if err != nil {
		return nil, err
	}

	//waitgroup is a counter thata lets one goroutine waits for other
	var wg sync.WaitGroup

	for _, post := range posts {
		post := post
		id := post["id"].(string)

		wg.Add(3)

		go func() {
			defer wg.Done()
			post["like_count"] = s.repo.LikeCount(id)
		}()
		go func() {
			defer wg.Done()
			post["comment_count"] = s.repo.CommentCount(id)
			post["comments"] = s.repo.FetchComment(id)
		}()
		go func() {
			defer wg.Done()
			post["share_count"] = s.repo.ShareCount(id)
		}()

	}
	wg.Wait() //if counter never reached 0,it deadlock
	return posts, nil
}

func (s *PostService) SharePost(postId, userId string) error {
	return s.repo.SharePost(postId, userId)
}
