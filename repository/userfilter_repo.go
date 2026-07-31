package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type UserFilterRepository interface {
	FilterUsersByName() ([]model.Auth, error)
}

type userFilterRepo struct{}

func NewUserFilterRepository() UserFilterRepository {
	return &userFilterRepo{}
}

func (r *userFilterRepo) FilterUsersByName() ([]model.Auth, error) {
	query := `
SELECT
    a.auth_id,
    a.username,
    COALESCE(u.avatar, '') AS avatar
FROM cometoseeauth a
LEFT JOIN userdetailinfo u
    ON a.auth_id = u.auth_id
`
	rows, err := intailizer.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.Auth
	for rows.Next() {
		var u model.Auth
		if err := rows.Scan(&u.Auth_id, &u.Username, &u.Avatar); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
