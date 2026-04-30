package model

import "time"

type ConnectionStatus string

const (
	Pending  ConnectionStatus = "pending"
	Accepted ConnectionStatus = "accepted"
	Blocked  ConnectionStatus = "blocked"
)

type Connection struct {
	ID          string
	UserLow     string
	UserHigh    string
	RequestedBy string
	Status      ConnectionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
