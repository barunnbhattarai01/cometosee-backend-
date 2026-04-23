package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type AuthRepositry interface {
	CreateUser(user model.Auth) error
	GetUserByEmail(email string) (*model.Auth, error)
	ForgetPassword(email, password string) (string, error)
}

type authrepo struct{}

func NewAuthRepository() AuthRepositry {
	return &authrepo{}
}
func (r *authrepo) CreateUser(user model.Auth) error {
	query := `INSERT INTO cometoseeauth (username,email,password) VALUES ($1,$2,$3)`
	_, err := intailizer.DB.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func (r *authrepo) GetUserByEmail(email string) (*model.Auth, error) {
	var user model.Auth
	err := intailizer.DB.QueryRow(
		`SELECT username,password,email FROM cometoseeauth WHERE email=$1`, email,
	).Scan(&user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authrepo) ForgetPassword(email, password string) (string, error) {

	query := `update cometoseeauth set password=$1 where email=$2`
	_, err := intailizer.DB.Exec(query, password, email)
	if err != nil {
		return "", err
	}

	return "password updated successfully", nil
}
