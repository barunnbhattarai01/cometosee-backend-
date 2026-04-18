package model

import "time"

type Subscription struct {
	ID        int       `json:"id"`
	UserEmail string    `json:"user_email"`
	StartDate string    `json:"start_date"`
	EndDate   time.Time `json:"end_date,omitempty"`
}
