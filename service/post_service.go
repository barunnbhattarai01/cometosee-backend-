package service

import (
	"cometosee/repository"
	"context"
	"strconv"
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

func (s *PostService) UploadPost(authID int, caption, image string) (string, error) {
	return s.repo.CreatePOST(authID, caption, image)
}

func (s *PostService) LikePost(postId, authId string) (bool, error) {
	return s.repo.ToggleLike(postId, authId)
}

func (s *PostService) AddComment(postId, authId, comment string) error {
	return s.repo.AddComment(postId, authId, comment)
}

func (s *PostService) FetchFeed(
	ctx context.Context,
	lat float64,
	lon float64,
	radius int,
	skill string,
) ([]map[string]interface{}, error) {

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	posts, err := s.repo.FetchFeed(ctx, lat, lon, radius, skill)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, 10)

	for _, post := range posts {
		post := post
		id := post["id"].(int)

		wg.Add(2)

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			likes := s.repo.LikeCount(strconv.Itoa(id))

			mu.Lock()
			post["likes"] = likes
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			comments := s.repo.CommentCount(strconv.Itoa(id))

			mu.Lock()
			post["comments"] = comments
			mu.Unlock()
		}()
	}

	wg.Wait()
	return posts, nil
}

func (s *PostService) SharePost(postId, authId string) error {
	return s.repo.SharePost(postId, authId)
}
