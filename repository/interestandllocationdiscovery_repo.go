package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

// we need to track a intereset and location table so we need auth_id as foreign key to link it to the user table
// and we need to track a location of user so we need a user_id as foreign key to link it to the user table
type InterestAndLocationDiscoveryRepository interface {
	FindNearbyUsers(
		lat float64,
		lon float64,
		radius int,
		sport string,
		skill string,
		currentUserId int,
	) ([]model.UserDetailInfo, error)
}

type interestAndLocationDiscoveryRepository struct{}

func NewInterestAndLocationDiscoveryRepository() InterestAndLocationDiscoveryRepository {
	return &interestAndLocationDiscoveryRepository{}
}

func (r *interestAndLocationDiscoveryRepository) FindNearbyUsers(
	lat float64,
	lon float64,
	radius int,
	sport string,
	skill string,
	currentUserId int,
) ([]model.UserDetailInfo, error) {

	query := `
	SELECT u.user_detail_id, u.auth_id, u.calling_name, u.sport, u.skill,u.bio,u.avatar,u.created_at
FROM location l
JOIN userdetailinfo u 
  ON l.user_detail_id = u.user_detail_id
WHERE l.geom IS NOT NULL
AND ST_DWithin(
  l.geom,
  ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
  $3::double precision
)
AND u.user_detail_id != $4
AND LOWER(u.sport) = LOWER($5)
AND LOWER(u.skill) = LOWER($6)
ORDER BY ST_Distance(
  l.geom,
  ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
)
LIMIT 20;
	`

	rows, err := intailizer.DB.Query(query, lon, lat, radius, currentUserId, sport, skill)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserDetailInfo

	for rows.Next() {
		var u model.UserDetailInfo
		err := rows.Scan(
			&u.User_Detail_Id,
			&u.AuthId,
			&u.Calling_name,
			&u.Sport,
			&u.Skill,
			&u.Bio,
			&u.Avatar,
			&u.Created_at,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

//what we need to do is simpple awhh it not simple not let say it it simple
//we need to take the user interest like cricket so like skill and location then first prioery is our is location
//so we look up into location first and whoever is neaer taht we will match interest and skill then we display it
//does it simple yeah but it is not good for real application and we are buliding real app so what??
//firest we will do in db but later we will do it using our own ai model to detect user interest and location
