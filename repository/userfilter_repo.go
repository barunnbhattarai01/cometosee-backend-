package repository

import "cometosee/intailizer"

type UserFilterRepository interface {
	FilterUsersByName() ([]string, error)
}

type userFilterRepo struct{}

func NewUserFilterRepository() UserFilterRepository {
	return &userFilterRepo{}
}

func (r *userFilterRepo) FilterUsersByName() ([]string, error) {
	query := `SELECT username FROM cometoseeauth`
	rows, err := intailizer.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}
