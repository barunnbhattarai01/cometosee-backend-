package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type UserDetailInfo interface {
	TakeUserDetailInfo(user *model.UserDetailInfo) (string, error)
	TakeUserLocation(user *model.Location) (string, error)
	IsProfileExists(auth_id int) (bool, error)
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

	query := `
	INSERT INTO location(
		user_detail_id,
		country,
		city,
		latitude,
		longitude,
		geom
	)
	VALUES(
		$1, $2, $3, $4, $5,
		ST_SetSRID(ST_MakePoint($5, $4), 4326)::geography
	)
	`

	_, err := intailizer.DB.Exec(
		query,
		user.User_Detail_Id,
		user.Country,
		user.City,
		user.Latitude,
		user.Longitude,
	)

	if err != nil {
		return "", err
	}

	return "User Location Information inserted successfully", nil
}

func (r *userDetailInfoRepo) IsProfileExists(authId int) (bool, error) {

	query := `SELECT EXISTS (
		SELECT 1 FROM userdetailinfo WHERE auth_id = $1
	)`

	var exists bool

	err := intailizer.DB.QueryRow(query, authId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
