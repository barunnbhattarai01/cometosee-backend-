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

type AdminService interface {
	Login(email, password string) (string, *model.Admin, error)
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) Login(email, password string) (string, *model.Admin, error) {
	admin, err := s.repo.GetAdminByEmail(strings.TrimSpace(email))
	if err != nil || admin == nil || !admin.IsActive {
		return "", nil, errors.New("invalid admin credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid admin credentials")
	}

	claims := jwt.MapClaims{
		"adminId":  admin.AdminID,
		"username": admin.Username,
		"email":    admin.Email,
		"role":     "admin",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("ADMIN_SECRET")
	if secret == "" {
		secret = os.Getenv("SECRET")
	}
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}
	admin.Password = ""
	return signed, admin, nil
}
