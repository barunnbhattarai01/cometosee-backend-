package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(user model.Auth) error
	Login(email, password string) (string, string, error)
}

type authService struct {
	repo repository.AuthRepositry
}

func NewAuthService(repo repository.AuthRepositry) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Signup(user model.Auth) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	user.Email = strings.ToLower(user.Email)

	return s.repo.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, string, error) {
	email = strings.ToLower(email)
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}

	token, err := geneareJwt(user.Email, user.Username)
	if err != nil {
		return "", "", err
	}
	return token, user.Username, nil

}

func geneareJwt(email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
