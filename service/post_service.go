package service

import (
	"cometosee/model"
	"cometosee/repository"
	"context"
	"errors"
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

func (s *PostService) UploadPost(authID int, caption, image, venue string) (int, error) {
	return s.repo.CreatePOST(authID, caption, image, venue)
}

func (s *PostService) LikePost(postId, authId int) (bool, error) {
	return s.repo.ToggleLike(postId, authId)
}

func (s *PostService) AddComment(postId int, authId int, comment string) error {
	return s.repo.AddComment(postId, authId, comment)
}

func (s *PostService) FetchFeed(
	ctx context.Context,
	authId int,
	lat float64,
	lon float64,
	radius int,

) ([]map[string]interface{}, error) {

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	sport, err := s.repo.GetUserSport(authId)
	if err != nil {
		return nil, err
	}

	posts, err := s.repo.FetchFeed(ctx, lat, lon, radius, sport)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, 10)

	for _, post := range posts {
		post := post
		id := post["id"].(int)

		wg.Add(5)

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			likes := s.repo.LikeCount(id)

			mu.Lock()
			post["likes"] = likes
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			comments := s.repo.CommentCount(id)

			mu.Lock()
			post["comments"] = comments
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			commentsList, err := s.repo.GetComments(id)
			if err != nil {
				return
			}

			mu.Lock()
			post["comments_list"] = commentsList
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			isLiked, err := s.repo.Islike(id, authId)
			if err != nil {
				return
			}

			mu.Lock()
			post["is_liked"] = isLiked
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			isJoinedSlot, err := s.repo.IsJoinedSlot(id, authId)
			if err != nil {
				return
			}

			mu.Lock()
			post["is_joined"] = isJoinedSlot
			mu.Unlock()
		}()
	}

	wg.Wait()
	return posts, nil
}

func (s *PostService) SharePost(postId, authId int) error {
	return s.repo.SharePost(postId, authId)
}

// slot
func (s *PostService) CreateSlot(slot model.PostSlot) (int, error) {
	if slot.EndTime.Before(slot.StartTime) {
		return 0, errors.New("end time must be after start time")
	}
	slotid, err := s.repo.CreateSlot(slot)
	if err != nil {
		return 0, err
	}

	return slotid, nil
}

func (s *PostService) JoinSlot(slotID, authID int) error {
	return s.repo.JoinSlotTx(slotID, authID)
}

// fetch joined user from slot
func (s *PostService) GetSlotParticipants(slotID int, authId int) ([]map[string]interface{}, error) {
	return s.repo.GetSlotParticipants(slotID, authId)
}
