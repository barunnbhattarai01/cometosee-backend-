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
	posts, err := s.repo.FetchPost(ctx)

	if err != nil {
		return nil, err
	}

	//waitgroup is a counter thata lets one goroutine waits for other
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, 10) //limit concurrency to 10

	for _, post := range posts {
		post := post
		id := post["id"].(string)

		wg.Add(3)

		go func() {
			defer wg.Done()
			sem <- struct{}{}        //acquire semaphore
			defer func() { <-sem }() //release semaphore
			count := s.repo.LikeCount(ctx, id)
			mu.Lock()
			post["like_count"] = count
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			sem <- struct{}{}        //acquire semaphore
			defer func() { <-sem }() //release semaphore

			comments := s.repo.FetchComment(ctx, id)
			commentCount := s.repo.CommentCount(ctx, id)

			mu.Lock()
			post["comment_count"] = commentCount
			post["comments"] = comments
			mu.Unlock()

		}()
		go func() {
			defer wg.Done()
			//semaphore is like traffic controller and limits number of goroutine accessing a resource
			sem <- struct{}{}        //acquire semaphore
			defer func() { <-sem }() //release semaphore

			count := s.repo.ShareCount(ctx, id)
			mu.Lock()
			post["share_count"] = count
			mu.Unlock()
		}()

	}
	wg.Wait() //if counter never reached 0,it deadlock
	return posts, nil
}

func (s *PostService) SharePost(postId, userId string) error {
	return s.repo.SharePost(postId, userId)
}

func (s *PostService) Latestlikes(ctx context.Context, postId string) ([]string, error) {
	return s.repo.Latest10Likers(ctx, postId)
}
