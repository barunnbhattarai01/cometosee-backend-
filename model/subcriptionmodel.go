package model

import "time"

type Subscription struct {
	ID        int       `json:"id"`
	AuthID    int       `json:"auth_id"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Payment struct {
	ID              int       `json:"id"`
	AuthID          int       `json:"auth_id"`
	TransactionUUID string    `json:"transaction_uuid"`
	Plan            string    `json:"plan"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	PaymentDate     time.Time `json:"payment_date"`
}
