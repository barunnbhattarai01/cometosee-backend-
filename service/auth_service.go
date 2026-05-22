package service

import (
	"cometosee/model"
	"cometosee/repository"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(user model.Auth) error
	Login(email, password string) (string, string, error)
	ForgetPassword(email string) error
	ResetPassword(email, otp, newPassword string) error
	GetProfile(email string) (model.Auth, error)
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
	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("invalid email format")
	}

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

	token, err := geneareJwt(user.Auth_id, user.Email, user.Username)
	if err != nil {
		return "", "", err
	}
	return token, user.Username, nil

}

func geneareJwt(authid int, email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"authId":   float64(authid),
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func genrateOtp() string {
	otp := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic(err)
		}
		otp[i] = byte(n.Int64()) + '0'
	}
	return string(otp)
}

func SaveOtp(idetifier string, otp string) {
	model.OTPStored[idetifier] = model.OTP{
		Code:      otp,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
}

func verifyotp(identifier string, otp string) error {
	data, exists := model.OTPStored[identifier]
	if !exists {
		return errors.New("otp not found")
	}
	//check expiry first
	if data.ExpiredAt.Before(time.Now()) {
		delete(model.OTPStored, identifier)
		return errors.New("otp expired")
	}
	//check code
	if data.Code != otp {
		return errors.New("otp didnot match")
	}
	//valid otp, delete it and return true
	delete(model.OTPStored, identifier)
	return nil
}

func (s *authService) ForgetPassword(email string) error {
	email = strings.ToLower(email)

	_, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	otp := genrateOtp()
	SaveOtp(email, otp)

	if err := sendEmail(email, otp); err != nil {
		return err
	}

	return nil
}

func sendEmail(to string, otp string) error {
	from := "barunnbhattarai@gmail.com"
	password := os.Getenv("APP_PASSWORD")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf("Subject: Password Reset OTP\n\nYour OTP is: %s", otp))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err

}

func (s *authService) ResetPassword(email, otp, newPassword string) error {
	email = strings.ToLower(email)

	// Verify OTP
	err := verifyotp(email, otp)
	if err != nil {
		return errors.New("error in verifying otp")
	}

	// hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}

	// Update password in DB
	err = s.repo.ForgetPassword(email, string(hash))
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) GetProfile(email string) (model.Auth, error) {
	email = strings.ToLower(email)

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return model.Auth{}, errors.New("user not found")
	}

	user.Password = ""
	return *user, nil
}
