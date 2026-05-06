package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

// we need to track a intereset and location table so we need auth_id as foreign key to link it to the user table
// and we need to track a location of user so we need a user_id as foreign key to link it to the user table
type InterestAndLocationDiscoveryRepository interface {
	DiscoverInterests(auth_id int, user_detail_id int) ([]*model.UserDetailInfo, error)
	DiscoverLocations(auth_id int, user_detail_id int) ([]*model.Location, error)
}

type interestAndLocationDiscoveryRepository struct{}

func NewInterestAndLocationDiscoveryRepository() InterestAndLocationDiscoveryRepository {
	return &interestAndLocationDiscoveryRepository{}
}

func (r *interestAndLocationDiscoveryRepository) DiscoverInterests(auth_id int, user_detail_id int) ([]*model.UserDetailInfo, error) {

	query := `SELECT * FROM user_interests WHERE auth_id = $1 AND user_detail_id = $2`

	rows, err := intailizer.DB.Query(query, auth_id, user_detail_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userinterest []*model.UserDetailInfo
	for rows.Next() {
		var interest *model.UserDetailInfo
		if err := rows.Scan(&interest.AuthId, &interest.User_Detail_Id, &interest.Sport, &interest.Skill); err != nil {
			return nil, err
		}
		userinterest = append(userinterest, interest)
	}

	return userinterest, nil
}

func (r *interestAndLocationDiscoveryRepository) DiscoverLocations(auth_id int, user_detail_id int) ([]*model.Location, error) {

	query := `SELECT * FROM user_locations WHERE auth_id = $1 AND user_detail_id = $2`

	rows, err := intailizer.DB.Query(query, auth_id, user_detail_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userlocation []*model.Location
	for rows.Next() {
		var location *model.Location
		if err := rows.Scan(&location.User_Detail_Id, &location.City, &location.Country); err != nil {
			return nil, err
		}
		userlocation = append(userlocation, location)
	}

	return userlocation, nil

}

//what we need to do is simpple awhh it not simple not let say it it simple
//we need to take the user interest like cricket so like skill and location then first prioery is our is location
//so we look up into location first and whoever is neaer taht we will match interest and skill then we display it
//does it simple yeah but it is not good for real application and we are buliding real app so what??
//firest we will do in db but later we will do it using our own ai model to detect user interest and location
