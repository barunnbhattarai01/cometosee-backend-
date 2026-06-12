package service

import (
	"cometosee/model"
	"cometosee/repository"
)

type UserFilterService interface {
	FilterUsersByName() ([]model.Auth, error)
}

type userFilterService struct {
	repo repository.UserFilterRepository
}

func NewUserFilterService(repo repository.UserFilterRepository) UserFilterService {
	return &userFilterService{repo: repo}
}

func (s *userFilterService) FilterUsersByName() ([]model.Auth, error) {
	return s.repo.FilterUsersByName()
}
