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
	GetUserDetailIDByAuthID(authID int) (int, error)
	UpdateUserDetailInfo(user *model.UserDetailInfo) (string, error)
	UpdateLocation(user *model.Location) (string, error)
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

	repomessage, err := s.repo.TakeUserDetailInfo(user)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return repomessage, nil
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

func (s *userDetailInfoService) GetUserDetailIDByAuthID(authID int) (int, error) {
	return s.repo.GetUserDetailIDByAuthID(authID)
}

func (s *userDetailInfoService) UpdateUserDetailInfo(user *model.UserDetailInfo) (string, error) {

	if user.AuthId == 0 {
		return "", errors.New("auth_id is required")
	}

	if user.Calling_name == "" || user.Sport == "" || user.Skill == "" || user.Bio == "" {
		return "", errors.New("All fields are required")
	}

	message, err := s.repo.UpdateUserDetailInfo(user)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (s *userDetailInfoService) UpdateLocation(user *model.Location) (string, error) {

	if user.User_Detail_Id == 0 {
		return "", errors.New("user_detail_id is required")
	}

	if user.Country == "" || user.City == "" || user.Latitude == 0 || user.Longitude == 0 {
		return "", errors.New("All fields are required")
	}

	message, err := s.repo.UpdateLocation(user)
	if err != nil {
		return "", err
	}

	return message, nil
}
