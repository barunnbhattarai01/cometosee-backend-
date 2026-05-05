package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type UserDetailInfo interface {
	TakeUserDetailInfo(user *model.UserDetailInfo) (string, error)
	TakeUserLocation(user *model.Location) (string, error)
}

type userDetailInfoRepo struct{}

func NewUserDetailInfoRepository() UserDetailInfo {
	return &userDetailInfoRepo{}
}

func (r *userDetailInfoRepo) TakeUserDetailInfo(user *model.UserDetailInfo) (string, error) {

	query := `insert into userdetailinfo(auth_id, calling_name, sport, skill, avatar, bio) values($1, $2, $3, $4, $5, $6)`
	_, err := intailizer.DB.Exec(query, user.AuthId, user.Calling_name, user.Sport, user.Skill, user.Avatar, user.Bio)
	if err != nil {
		return "", err
	}

	return "User Detail Information insert sucessfully", nil
}

func (r *userDetailInfoRepo) TakeUserLocation(user *model.Location) (string, error) {

	query := `insert into location(user_detail_id, country, city, latitude, longitude) values($1, $2, $3, $4, $5)`
	_, err := intailizer.DB.Exec(query, user.User_Detail_Id, user.Country, user.City, user.Latitude, user.Longitude)
	if err != nil {
		return "", err
	}

	return "User Location Information insert sucessfully", nil
}
