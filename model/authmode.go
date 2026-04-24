package model

import "time"

type Auth struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      *OTP   `json:"otp"`
}

type OTP struct {
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

var OTPStored = map[string]OTP{}
var Email string
