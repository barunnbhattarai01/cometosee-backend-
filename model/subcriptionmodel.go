package model

type Subscription struct {
	ID        int     `json:"id"`
	UserEmail string  `json:"user_email"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
}
