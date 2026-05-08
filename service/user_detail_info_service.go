package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
)

type UserDetailInfoService interface {
	TakeUserDetailInfo(user *model.UserDetailInfo) (string, error)
	TakeUserLocation(user *model.Location) (string, error)
	IsProfileCompleted(auth_id int) (bool, error)
}

type userDetailInfoService struct {
	repo repository.UserDetailInfo
}

func NewUserDetailInfoService(repo repository.UserDetailInfo) UserDetailInfoService {

	return &userDetailInfoService{repo: repo}
}

func (s *userDetailInfoService) TakeUserDetailInfo(user *model.UserDetailInfo) (string, error) {

	if user.Calling_name == "" || user.Sport == "" || user.Skill == "" || user.Bio == "" {
		return "", errors.New("All fields are required")
	}

	return s.repo.TakeUserDetailInfo(user)
}

func (s *userDetailInfoService) TakeUserLocation(user *model.Location) (string, error) {
	if user.Country == "" || user.City == "" || user.Latitude == 0 || user.Longitude == 0 {
		return "", errors.New("All fields are required")
	}
	return s.repo.TakeUserLocation(user)
}

func (s *userDetailInfoService) IsProfileCompleted(authId int) (bool, error) {
	return s.repo.IsProfileExists(authId)
}
