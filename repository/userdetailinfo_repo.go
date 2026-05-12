package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"fmt"
)

type UserDetailInfo interface {
	TakeUserDetailInfo(user *model.UserDetailInfo) (string, error)
	TakeUserLocation(user *model.Location) (string, error)
	IsProfileExists(auth_id int) (bool, error)
	GetUserDetailIDByAuthID(authID int) (int, error)
	UpdateUserDetailInfo(user *model.UserDetailInfo) (string, error)
	UpdateLocation(user *model.Location) (string, error)
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

	return "sucessfully inserted", nil
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

func (r *userDetailInfoRepo) GetUserDetailIDByAuthID(authID int) (int, error) {
	query := `SELECT user_detail_id FROM userdetailinfo WHERE auth_id = $1`

	var id int
	err := intailizer.DB.QueryRow(query, authID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userDetailInfoRepo) UpdateUserDetailInfo(user *model.UserDetailInfo) (string, error) {

	query := `
	UPDATE userdetailinfo
	SET 
		calling_name = $1,
		sport = $2,
		skill = $3,
		avatar = $4,
		bio = $5
	WHERE auth_id = $6
	`

	result, err := intailizer.DB.Exec(
		query,
		user.Calling_name,
		user.Sport,
		user.Skill,
		user.Avatar,
		user.Bio,
		user.AuthId,
	)

	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("no user found with given auth_id")
	}

	return "User detail updated successfully", nil
}

func (r *userDetailInfoRepo) UpdateLocation(user *model.Location) (string, error) {
	query := `
	UPDATE location
	SET 
		country = $1,
		city = $2,
		latitude = $3,
		longitude = $4,
		geom = ST_SetSRID(ST_MakePoint($4, $3), 4326)::geography
	WHERE user_detail_id = $5
	`

	result, err := intailizer.DB.Exec(
		query,
		user.Country,
		user.City,
		user.Latitude,
		user.Longitude,
		user.User_Detail_Id,
	)

	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("no location found for given user_detail_id")
	}

	return "User location updated successfully", nil
}
