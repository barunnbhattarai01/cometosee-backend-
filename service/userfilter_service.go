package service

import "cometosee/repository"

type UserFilterService interface {
	FilterUsersByName() ([]string, error)
}

type userFilterService struct {
	repo repository.UserFilterRepository
}

func NewUserFilterService(repo repository.UserFilterRepository) UserFilterService {
	return &userFilterService{repo: repo}
}

func (s *userFilterService) FilterUsersByName() ([]string, error) {
	return s.repo.FilterUsersByName()
}
