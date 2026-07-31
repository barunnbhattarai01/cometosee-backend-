package model

import "time"

type Auth struct {
	Auth_id  int    `json:"auth_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Otp      *OTP   `json:"otp"`
}

type OTP struct {
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

var OTPStored = map[string]OTP{}
var Email string

type PendingSignup struct {
	User      Auth
	ExpiredAt time.Time
}

var PendingUsers = map[string]PendingSignup{}
