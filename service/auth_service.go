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
	VerifyEmail(email, otp string) (bool, error)
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

	otp := genrateOtp()
	SaveOtp(user.Email, otp)

	model.PendingUsers[user.Email] = model.PendingSignup{
		User:      user,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	if err := sendEmail(user.Email, otp); err != nil {
		delete(model.PendingUsers, user.Email)
		return err
	}

	return nil
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

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Email Verification</title>
</head>

<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,Helvetica,sans-serif;">

<table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px;background:#f4f6f9;">
<tr>
<td align="center">

<table width="600" cellpadding="0" cellspacing="0"
style="background:#ffffff;border-radius:12px;overflow:hidden;
box-shadow:0 6px 20px rgba(0,0,0,.08);">

<tr>
<td style="background:#ff6b00;padding:25px;text-align:center;color:white;">
<h1 style="margin:0;"> Cometosee</h1>
</td>
</tr>

<tr>
<td style="padding:40px;">

<h2 style="margin-top:0;color:#333;">
Verify Your Email
</h2>

<p style="font-size:16px;color:#555;line-height:28px;">
Welcome to <strong>Cometosee</strong>!
</p>

<p style="font-size:16px;color:#555;line-height:28px;">
Use the following One-Time Password (OTP) to verify your email address:
</p>

<div style="
margin:35px 0;
text-align:center;
font-size:34px;
font-weight:bold;
letter-spacing:8px;
color:#ff6b00;
background:#fff4eb;
padding:20px;
border-radius:10px;
border:2px dashed #ff6b00;
">
%s
</div>

<p style="font-size:15px;color:#555;line-height:26px;">
This OTP is valid for <strong>10 minutes</strong>.
Please do not share this code with anyone.
</p>

<p style="font-size:15px;color:#555;">
If you did not request this verification, you can safely ignore this email.
</p>

</td>
</tr>

<tr>
<td style="background:#f8f8f8;padding:20px;text-align:center;font-size:13px;color:#888;">
© 2026 Cometosee. All Rights Reserved.
</td>
</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`, otp)

	message := []byte(fmt.Sprintf(
		"From: Cometosee <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Verify Your Email - OTP\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n%s",
		from,
		to,
		html,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)
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

func (s *authService) VerifyEmail(email, otp string) (bool, error) {

	email = strings.ToLower(email)

	if err := verifyotp(email, otp); err != nil {
		return false, err
	}

	pending, exists := model.PendingUsers[email]
	if !exists {
		return false, errors.New("signup request not found")
	}

	if err := s.repo.CreateUser(pending.User); err != nil {
		return false, err
	}

	delete(model.PendingUsers, email)

	return true, nil
}
